mod api;
mod arg;
mod git;
mod indexer;
mod parser;
mod solr;
mod utils;

use arg::Arg;
use indexer::Indexer;
use std::env;

#[tokio::main]
pub async fn main() {
    let args: Vec<String> = env::args().collect();
    
    // Check if we're running in API mode
    if args.len() > 1 && args[1] == "api" {
        let port = env::var("API_PORT")
            .unwrap_or_else(|_| "8080".to_string())
            .parse::<u16>()
            .unwrap_or(8080);
            
        println!("Starting Heline Indexer in API mode on port {}", port);
        api::start_api_server(port).await;
        return;
    }
    
    // Otherwise, run in CLI mode
    println!("Starting Heline Indexer in CLI mode");
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
