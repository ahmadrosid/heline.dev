use serde::Serialize;

#[derive(Serialize, Clone, Debug)]
pub struct GitFile {
    pub id: String,
    pub file_id: String,
    pub owner_id: String,
    pub path: String,
    pub repo: String,
    pub branch: String,
    pub lang: String,
    pub content: Vec<String>,
}

#[derive(Serialize, Clone, Debug)]
pub struct GithubFileUpdate {
    pub id: String,
    pub content: AddString,
}

#[derive(Serialize, Clone, Debug)]
pub struct AddString {
    pub add: Vec<String>,
}

pub async fn insert(data: &GitFile, base_url: &str) -> Result<String, reqwest::Error> {
    let mut body: Vec<GitFile> = Vec::new();
    body.push(data.clone());
    let url = format!(
        "{}/solr/heline/update?&commitWithin=1000&overwrite=false&wt=json",
        base_url
    );
    let client = reqwest::Client::new();
    let res = client.post(url).json(&body).send().await?;
    let json = res.text().await?;
    Ok(json)
}

pub async fn update(data: &GitFile, base_url: &str) -> Result<String, reqwest::Error> {
    let data = data.clone();
    let mut body: Vec<GithubFileUpdate> = Vec::new();
    let update = GithubFileUpdate {
        id: data.id,
        content: AddString { add: data.content },
    };

    body.push(update);
    let url = format!("{}/solr/heline/update", base_url);
    let client = reqwest::Client::new();
    let res = client.post(url).json(&body).send().await?;
    let json = res.text().await?;
    Ok(json)
}
