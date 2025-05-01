use reqwest::Response;
use std::collections::HashMap;

pub async fn get_user_id(username: &str) -> Result<String, String> {
    let user_agent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36".to_owned();

    let url = format!("https://api.github.com/users/{}", username);
    let client = reqwest::Client::new();
    return match client
        .get(url)
        .header("User-Agent", user_agent)
        .send()
        .await
    {
        Ok(res) => extract_id(res).await,
        Err(e) => Err(e.to_string()),
    };
}

pub async fn get_repo(repo: &str) -> Result<String, String> {
    let user_agent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36".to_owned();

    let url = format!("https://api.github.com/repos/{}", repo);
    let client = reqwest::Client::new();
    return match client
        .get(url)
        .header("User-Agent", user_agent)
        .send()
        .await
    {
        Ok(res) => extract_id(res).await,
        Err(e) => Err(e.to_string()),
    };
}

async fn extract_id(res: Response) -> Result<String, String> {
    return match res.json::<HashMap<String, serde_json::Value>>().await {
        Ok(json) => {
            return match json.get("id") {
                Some(id) => Ok(id.to_string()),
                _ => Err("Not found!".to_string()),
            }
        }
        Err(e) => Err(e.to_string()),
    };
}
