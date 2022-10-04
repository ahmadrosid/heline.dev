use crate::query::Query;

use std::collections::HashMap;
use std::collections::HashSet;
use std::iter::FromIterator;

#[derive(Debug, Clone)]
pub struct Document {
  pub id: usize,
  pub text: String,
}

#[derive(Debug)]
pub struct Index {
  pub name: String,
  pub data: HashMap<String, Vec<usize>>,
}

impl Index {
  pub fn new(name: &str) -> Self {
    Self {
      name: name.to_string(),
      data: HashMap::new()
    }
  }

  fn tokenize(text: &str) -> Vec<String> {
    if text.is_empty() {
        return vec!["".to_string()];
    }

    let chars = text.chars();
    let mut tokens = vec![];
    for c in chars.into_iter() {
      tokens.push(format!("{}", c))
    }

    tokens
  }

  pub fn insert(&mut self, docs: Vec<Document>) {
    for doc in docs {
      for tok in Self::tokenize(&doc.text) {
        self.data.insert(tok, vec![doc.id]);
      }
    }
  }

  fn intersection(a: Vec<usize>, b: Vec<usize>) -> Vec<usize> {
    println!("a.len: {}, b.len: {}", a.len(), b.len());
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
  
  fn hashset(data: &[usize]) -> HashSet<usize> {
    HashSet::from_iter(data.iter().cloned())
  }

  pub fn search(&self, query: Query) -> Vec<usize> {
    let mut ids = vec![];
    for tok in Self::tokenize(&query.text) {
      if let Some(doc_ids) = self.data.get(&tok as &str) {
        if ids.len() == 0 {
          ids = doc_ids.to_vec();
        } else {
          let a = Self::hashset(&ids.to_vec());
          let b = Self::hashset(&doc_ids.to_vec());

          let mut intersection = a.intersection(&b);
          ids = intersection.to_owned().into_iter().collect::<Vec<_>>();

        }
      }
    }
    ids
  }
}
