use crate::git;
use crate::parser;
use crate::solr;
use crate::solr::client::GitFile;
use crate::utils;

use ignore::Walk;
use select::document::Document;
use select::predicate::{Class, Name};
use std::path::{Path, PathBuf};

pub struct MetaIndexFile {
    path: PathBuf,
    root_path_len: usize,
    git_repo: String,
    user_id: String,
    branch: String,
    base_url: String,
    git_host: String,
}

#[derive(Clone)]
pub struct Indexer {
    pub repo_dir: PathBuf,
    pub git_url: String,
    pub base_url: String,
    pub with_delete_dir: bool,
    pub git_host: String,
    pub repo_name: String,
}

impl Indexer {
    pub fn new(repo_dir: PathBuf, git_url: &str, base_url: &str, with_delete_dir: bool) -> Self {
        let git_host = utils::get_url_host(git_url).unwrap_or("github.com".to_string());
        let repo_name = utils::get_repo_name(git_url);

        Self {
            repo_dir,
            git_url: git_url.to_string(),
            base_url: base_url.to_string(),
            with_delete_dir,
            git_host,
            repo_name,
        }
    }

    pub async fn process(&self) {
        let ssh_url = utils::get_git_ssh_url(&self.git_url);
        let success = git::clone_repo(&self.repo_dir, &ssh_url, &self.repo_name);

        if success {
            self.index_directory().await;
            if self.with_delete_dir {
                utils::delete_dir(&self.repo_dir.join(Path::new(&self.repo_name)));
            }
        } else {
            print!("{}\n", format!("Failed to clone: {}", ssh_url));
        }
    }

    pub async fn index_directory(&self) {
        println!("Start indexing on folder: {}", self.repo_dir.display());

        let mut total = 0;
        let git_repo = utils::get_git_repo_path(&self.git_url);
        let username = git_repo.split("/").next().unwrap();
        let branch = git::get_branch_name(&self.repo_dir);

        let user_id = match &self.git_host[..] {
            "gitlab.com" => String::from("0000"),
            _ => match git::github::get_user_id(username).await {
                Ok(user_id) => user_id,
                Err(e) => {
                    print!("{}", e);
                    String::from("00000")
                }
            },
        };

        let repo_name = utils::get_repo_name(&self.git_url);
        let walk_dir_path = self.repo_dir.join(repo_name);
        let dirs = Walk::new(&walk_dir_path).into_iter().filter_map(|v| v.ok());
        let root_path_len = self
            .repo_dir
            .to_str()
            .unwrap()
            .split("/")
            .collect::<Vec<_>>()
            .len();

        let excluded_files = [
            "composer.lock",
            "package-lock.json",
            "yarn.lock",
            "Cargo.lock",
            "go.sum",
            ".DS_Store",
        ];
        
        let excluded_extensions = [
            ".lock",
            ".min.js",
            ".min.css",
            ".map",
            ".woff",
            ".woff2",
            ".ttf",
            ".eot",
            ".svg",
            ".png",
            ".jpg",
            ".jpeg",
            ".gif",
            ".ico",
            ".pdf",
            ".zip",
            ".tar",
            ".gz",
            ".mp3",
            ".mp4",
            ".mov",
            ".avi",
        ];

        for entry in dirs {
            if !entry.path().is_file() {
                continue;
            }
            
            let file_name = entry.path().file_name().unwrap_or_default().to_string_lossy();
            
            if excluded_files.contains(&file_name.as_ref()) || 
               excluded_extensions.iter().any(|&ext| file_name.ends_with(ext)) {
                println!("Skipping excluded file: {}", entry.path().display());
                continue;
            }
            
            if let Ok(metadata) = entry.metadata() {
                if metadata.len() > 5 * 1024 * 1024 {
                    println!("Skipping large file: {} ({} bytes)", entry.path().display(), metadata.len());
                    continue;
                }
            }

            print!("{}\n", format!("Indexing {}", entry.path().display()));
            let meta = MetaIndexFile {
                path: PathBuf::from(entry.path()),
                git_repo: git_repo.to_string(),
                user_id: user_id.to_string(),
                branch: branch.to_string(),
                base_url: self.base_url.to_string(),
                git_host: self.git_host.to_string(),
                root_path_len,
            };
            self.process_file(meta).await;
            total += 1;
        }

        if total == 0 {
            print!(
                "{}\n",
                format!("Folder '{}' not found!", walk_dir_path.display())
            );
        } else {
            print!(
                "{}\n",
                format!("Done indexing '{}' total {} files!", git_repo, total)
            );
        }
    }

    async fn process_file(&self, meta: MetaIndexFile) {
        match parser::read_file(&meta.path) {
            Ok((input, lang)) => {
                let html = parser::render_html(input, lang);
                let paths = meta.path.to_str().unwrap().split("/").collect::<Vec<_>>();
                let file_path = paths[meta.root_path_len..paths.len()].to_vec().join("/");
                let id = [
                    meta.git_repo.to_string(),
                    paths[meta.root_path_len..paths.len()].to_vec().join("/"),
                ]
                .join("/");
                let data = GitFile {
                    id: id.to_owned(),
                    file_id: format!(
                        "{}/{}/{}",
                        &meta.git_host,
                        &meta.git_repo,
                        file_path.to_string()
                    ),
                    owner_id: meta.user_id.to_string(),
                    path: if meta.root_path_len >= 2 {
                        paths[meta.root_path_len - 2..paths.len().saturating_sub(1)]
                            .to_vec()
                            .join("/")
                    } else {
                        file_path.clone()
                    },
                    repo: meta.git_repo.to_string(),
                    branch: meta.branch.to_owned(),
                    lang: lang.to_string(),
                    content: Vec::new(),
                };
                self.store(data, &html, &meta.base_url).await;
            }
            Err(msg) => {
                print!("{}\n", msg);
            }
        }
    }

    async fn store(&self, mut data: GitFile, html: &str, base_url: &str) {
        // Process the document outside of async context to avoid Send issues
        // This way, Document (which is not Send) doesn't cross an await point
        let chunks = {
            let document = Document::from(html);
            self.process_document(&document)
        };
        
        let mut update = false;
        for chunk in chunks {
            data.content = vec![chunk];
            self.create_or_update(&mut update, &data, base_url).await;
        }
    }

    fn process_document(&self, document: &Document) -> Vec<String> {
        let mut chunks = Vec::new();
        if let Some(el) = document.find(Class("highlight-table")).last() {
            let mut current_chunk = String::new();
            let mut line_count = 0;
            const CHUNK_SIZE: usize = 3;
            const MAX_CHARS: usize = 2000;

            for tr in el.find(Name("tr")) {
                line_count += 1;
                let html = tr.html();
                current_chunk.push_str(&html);
                current_chunk.push('\n');

                if line_count >= CHUNK_SIZE || current_chunk.len() >= MAX_CHARS {
                    if !current_chunk.is_empty() {
                        chunks.push(current_chunk);
                        current_chunk = String::new();
                        line_count = 0;
                    }
                }
            }

            if !current_chunk.is_empty() {
                chunks.push(current_chunk);
            }
        }
        chunks
    }

    async fn create_or_update(&self, update: &mut bool, data: &GitFile, base_url: &str) {
        if !*update {
            match solr::client::insert(&data, base_url).await {
                Ok(_) => {},
                Err(e) => print!("{}\n", e.to_string()),
            }
            *update = true;
        } else {
            match solr::client::update(&data, base_url).await {
                Ok(_) => {},
                Err(e) => print!("{}\n", e.to_string()),
            }
        }
    }
}
