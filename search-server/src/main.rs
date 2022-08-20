mod index;

use index::{Document, InvertedIndex};

fn main() {
    let mut index = InvertedIndex::new();
    let doc = Document{
        id: 1,
        path_id: "path/1".to_string(),
        repo: "username/repo-name".to_string(),
        path: "path/to/file".to_string(),
        host: "github.com".to_string(),
        lang: "JSON".to_string(),
        branch: "master".to_string(),
        contents: vec![
            "Some data string".to_string()
        ]
    };
    index.insert(vec![doc]);

    let res = index.search("Some");
    println!("{:?}", res);
}
