use std::env;
use std::path::PathBuf;

pub struct Arg {
    pub index_file: PathBuf,
    pub folder: PathBuf,
    pub solr_url: String,
    pub is_index_folder: bool,
    pub with_delete_folder: bool,
}

impl Arg {
    pub fn new() -> Self {
        Self {
            index_file: PathBuf::new(),
            folder: PathBuf::new(),
            solr_url: String::new(),
            is_index_folder: false,
            with_delete_folder: false,
        }
    }

    // Parse to: hli index_file.json --folder some/path
    pub fn parse(mut self) -> Result<Self, String> {
        // First try to get SOLR_BASE_URL, then BASE_URL, then use default
        self.solr_url = match env::var("SOLR_BASE_URL") {
            Ok(val) => val,
            Err(_) => match env::var("BASE_URL") {
                Ok(val) => val,
                Err(_) => "http://localhost:8984".to_string(),
            },
        };

        let arg_input: Vec<String> = env::args().collect();
        for input in &arg_input {
            if input == "-h" || input == "--help" {
                self.print_help();
                std::process::exit(0);
            }

            if input == "--delete-dir" {
                self.with_delete_folder = true;
            }
        }

        self.index_file = match arg_input.get(1) {
            Some(file) => PathBuf::from(file),
            None => {
                return Err(format!("Index file path is required!"));
            }
        };

        if !self.index_file.exists() {
            return Err(format!(
                "Please input the correct index file! Current file is not exists: {}",
                self.index_file.display()
            ));
        }

        let option = match arg_input.get(2) {
            Some(option) => option,
            None => "",
        };

        self.is_index_folder = option == "--folder";
        let mut repo_folder = PathBuf::new();
        match arg_input.get(3) {
            Some(folder) => repo_folder.push(folder),
            None => repo_folder.push("repos"),
        };

        self.folder = repo_folder;

        Ok(self)
    }

    fn print_help(&self) {
        let help_text = vec![
            "hli 0.1.0",
            "Ahmad Rosid <alahmadrosid@gmail.com>",
            "",
            "Heline.dev indexer, turn source code to github like html syntax highlighted!",
            "",
            "Usage :",
            "    hli data.json --folder some/folder",
            "",
            "Options :",
            "    <INDEX_FILE>    Index file with json format",
            "    --folder        Custom folder to source code",
            "    -h --help       Print help text",
            "    --delete-dir    Delete directory after indexing.",
            "",
        ];
        println!("{}", help_text.join("\n"));
    }
}
