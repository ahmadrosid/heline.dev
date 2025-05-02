pub mod github;
pub mod gitlab;

use crate::utils;
use std::path::Path;
use std::process::Command;

pub fn get_branch_name(dir: &Path) -> String {
    let file_path = dir.join(".git/HEAD");
    return match std::fs::read_to_string(file_path) {
        Ok(file) => file
            .split("/")
            .last()
            .unwrap()
            .to_string()
            .trim_end_matches('\n')
            .to_string(),
        Err(_) => "main".to_string(),
    };
}

pub fn clone_repo(cwd: &Path, ssh_url: &str, repo_name: &str) -> bool {
    if !cwd.exists() {
        std::fs::create_dir(cwd).expect(&format!("Failed to create directory: {}", cwd.display()));
    }

    let cloned_repo_dir = cwd.join(Path::new(repo_name));
    if cloned_repo_dir.exists() {
        print!("Repository already cloned: {}\n", cloned_repo_dir.display());
        return true;
    }

    print!(
        "{}\n",
        format!("Cloning '{}' to {}/{}", ssh_url, cwd.display(), repo_name)
    );
    utils::exec_command(
        Command::new("git")
            .current_dir(cwd)
            .arg("clone")
            .arg(ssh_url)
            .arg(repo_name),
    )
}

pub async fn get_repo(git_url: &str) -> Result<String, String> {
    match utils::get_url_host(git_url) {
        None => Err(format!("Invalid git url {}", git_url)),
        Some(host) => {
            let git_repo = utils::get_git_repo_path(git_url);
            return match &host[..] {
                "github.com" => github::get_repo(&git_repo).await,
                "gitlab.com" => gitlab::get_repo(&git_repo).await,
                _ => Err(format!("Unsupported git host: {}", host)),
            };
        }
    }
}
