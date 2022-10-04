mod index;
mod query;

use index::{Document, Index};
use query::Query;

fn main() {
    let docs = vec![
        Document{
            id: 1,
            text: "fn tokenize(text: &str) -> Vec<String> {".to_string()
        },
        Document{
            id: 2,
            text: "    if text.is_empty() {".to_string()
        },
    ];
    let mut indices = Index::new("test");
    indices.insert(docs.to_vec());

    let res = indices.search(Query { text: "token".to_string() });
    for id in res {
        let doc = docs.get(id).to_owned();
        if doc.is_none() {
            continue;
        }
        println!("{:?}", doc);
    }
}
