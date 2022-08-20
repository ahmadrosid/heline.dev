use std::collections::HashMap;

#[derive(Debug)]
pub struct Document {
    pub id: usize,
    pub path_id: String,
    pub repo: String,
    pub path: String,
    pub host: String,
    pub lang: String,
    pub branch: String,
    pub contents: Vec<String>
}

trait CollectChars {
    fn collect_chars(&self) -> Vec<char>;
}

impl CollectChars for Vec<String> {
    fn collect_chars(&self) -> Vec<char> {
        let mut chars: Vec<char> = vec![];
        for item in self.into_iter() {
            chars.append(&mut item.chars().collect::<Vec<char>>());
        }
        chars
    }
}

#[derive(Debug)]
pub struct InvertedIndex {
    data: HashMap<String, Vec<usize>>,
}

impl InvertedIndex {
    pub fn new() -> Self {
        Self {
            data: HashMap::new(),
        }
    }

    pub fn insert(&mut self, docs: Vec<Document>) {
        for doc in docs {
            for token in analyze(&doc.contents) {
                if let Some(ids) = self.data.get(&token as &str) {
                    if ids[ids.len() - 1] == doc.id {
                        continue;
                    }
                    let mut new_ids = ids.to_owned();
                    new_ids.push(doc.id);
                    self.data.insert(token, new_ids);
                } else {
                    self.data.insert(token, vec![doc.id]);
                }
            }
        }
    }

    pub fn search(&self, query: &str) -> Vec<usize> {
        let mut result = vec![];
        let vec_query: Vec<String> = vec![query.to_string()];
        for token in analyze(&vec_query) {
            if let Some(data) = self.data.get(&token as &str) {
                if result.len() == 0 {
                    result = data.to_owned();
                } else {
                    result = intersection(result, data.to_vec())
                }
            }
        }

        result
    }
}

fn analyze(text: &Vec<String>) -> Vec<String> {
    let mut stop_words = HashMap::new();
    stop_words.insert("a", "");
    stop_words.insert("dia", "");
    let tokens = tokenize(text);
    let tokens = stopword_filter(&tokens, stop_words);
    tokens
}

fn tokenize(text: &Vec<String>) -> Vec<String> {
    if text.is_empty() {
        return vec!["".to_string()];
    }

    text.collect_chars()
        .into_iter()
        .filter(|c| c.is_alphanumeric() || c.is_whitespace())
        .collect::<String>()
        .split_whitespace()
        .map(String::from)
        .collect()
}

fn stopword_filter(tokens: &Vec<String>, stop_words: HashMap<&str, &str>) -> Vec<String> {
    let mut new_tokens = Vec::new();
    for token in tokens {
        if stop_words.get(&token as &str).is_none() {
            new_tokens.push(token.to_owned())
        }
    }
    new_tokens
}

fn intersection(a: Vec<usize>, b: Vec<usize>) -> Vec<usize> {
    let mut max_len = a.len();
    if b.len() > max_len {
        max_len = b.len();
    }

    let mut result: Vec<usize> = vec![0; max_len];
    let mut i = 0;
    let mut j = 0;
    while i < a.len() && j < b.len() {
        if a[i] < b[i] {
            i += 1;
        } else if a[i] < b[j] {
            j += 1;
        } else {
            result.push(a[i]);
            i += 1;
            j += 1;
        }
    }

    result
}

