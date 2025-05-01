mod arg;
mod git;
mod indexer;
mod parser;
mod solr;
mod utils;

use arg::Arg;
use indexer::Indexer;

#[tokio::main]
pub async fn main() {
    let mut arg = Arg::new();
    match arg.parse() {
        Ok(new_arg) => arg = new_arg,
        Err(msg) => {
            eprintln!("{}", msg);
            std::process::exit(1);
        }
    }

    let value: Vec<String> = utils::parse_json(&arg.index_file);
    for git_url in value {
        match git::get_repo(&git_url).await {
            Ok(_repo_id) => {
                let indexer_service = Indexer::new(
                    arg.folder.clone(),
                    &git_url,
                    &arg.solr_url,
                    arg.with_delete_folder,
                );
                indexer_service.process().await;
            }
            Err(e) => {
                print!("{}: Error {}\n", git_url, e);
                continue;
            }
        };
    }
}
