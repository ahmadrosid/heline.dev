use crate::indexer::Indexer;
use std::path::PathBuf;
use warp::{Filter, Rejection, Reply};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::sync::Mutex;
use std::collections::HashMap;

#[derive(Debug, Deserialize)]
pub struct IndexRequest {
    git_url: String,
}

#[derive(Debug, Serialize)]
pub struct IndexResponse {
    status: String,
    message: String,
    job_id: Option<String>,
}

#[derive(Debug, Clone, Serialize)]
pub struct JobStatus {
    git_url: String,
    status: String,
    created_at: chrono::DateTime<chrono::Utc>,
    completed_at: Option<chrono::DateTime<chrono::Utc>>,
    message: Option<String>,
}

type Jobs = Arc<Mutex<HashMap<String, JobStatus>>>;

pub async fn start_api_server(port: u16) {
    let jobs = Jobs::default();
    let jobs_filter = warp::any().map(move || jobs.clone());

    // Health check endpoint
    let health = warp::path("health")
        .and(warp::get())
        .map(|| warp::reply::json(&IndexResponse {
            status: "ok".to_string(),
            message: "Heline Indexer API is running".to_string(),
            job_id: None,
        }));

    // Index endpoint
    let index = warp::path("index")
        .and(warp::post())
        .and(warp::body::json())
        .and(jobs_filter.clone())
        .and_then(handle_index);

    // Job status endpoint
    let job_status = warp::path!("jobs" / String)
        .and(warp::get())
        .and(jobs_filter.clone())
        .and_then(handle_job_status);

    // Jobs list endpoint
    let jobs_list = warp::path("jobs")
        .and(warp::get())
        .and(jobs_filter.clone())
        .and_then(handle_jobs_list);

    let routes = health
        .or(index)
        .or(job_status)
        .or(jobs_list)
        .with(warp::cors().allow_any_origin());

    println!("Starting Heline Indexer API server on port {}", port);
    warp::serve(routes).run(([0, 0, 0, 0], port)).await;
}

async fn handle_index(
    req: IndexRequest,
    jobs: Jobs,
) -> Result<impl Reply, Rejection> {
    let job_id = uuid::Uuid::new_v4().to_string();
    let git_url = req.git_url.clone();
    
    println!("Received indexing request for: {}", git_url);
    
    // Create a new job status entry
    let job_status = JobStatus {
        git_url: git_url.clone(),
        status: "queued".to_string(),
        created_at: chrono::Utc::now(),
        completed_at: None,
        message: None,
    };
    
    // Store the job status
    {
        let mut jobs_map = jobs.lock().await;
        jobs_map.insert(job_id.clone(), job_status);
    }
    
    // Clone the jobs for the async task
    let jobs_clone = jobs.clone();
    let job_id_clone = job_id.clone();
    
    // Spawn a background task to handle the indexing
    tokio::spawn(async move {
        let base_url = std::env::var("SOLR_BASE_URL").unwrap_or_else(|_| "http://localhost:8984".to_string());
        let repo_dir = PathBuf::from("repos");
        
        // Update job status to running
        {
            let mut jobs_map = jobs_clone.lock().await;
            if let Some(job) = jobs_map.get_mut(&job_id_clone) {
                job.status = "running".to_string();
            }
        }
        
        // Create and run the indexer
        let indexer = Indexer::new(
            repo_dir,
            &git_url,
            &base_url,
            false, // Don't delete folder after indexing
        );
        
        // Process the repository
        let result = indexer.process().await;
        
        // Update job status based on result
        {
            let mut jobs_map = jobs_clone.lock().await;
            if let Some(job) = jobs_map.get_mut(&job_id_clone) {
                job.status = "completed".to_string();
                job.completed_at = Some(chrono::Utc::now());
                job.message = Some("Indexing completed successfully".to_string());
            }
        }
    });
    
    // Return the job ID to the client
    Ok(warp::reply::json(&IndexResponse {
        status: "accepted".to_string(),
        message: format!("Indexing job for {} has been queued", git_url),
        job_id: Some(job_id),
    }))
}

async fn handle_job_status(
    job_id: String,
    jobs: Jobs,
) -> Result<impl Reply, Rejection> {
    let jobs_map = jobs.lock().await;
    
    match jobs_map.get(&job_id) {
        Some(job) => Ok(warp::reply::json(job)),
        None => {
            let response = IndexResponse {
                status: "error".to_string(),
                message: format!("Job with ID {} not found", job_id),
                job_id: None,
            };
            Ok(warp::reply::json(&response))
        }
    }
}

async fn handle_jobs_list(
    jobs: Jobs,
) -> Result<impl Reply, Rejection> {
    let jobs_map = jobs.lock().await;
    let jobs_vec: Vec<&JobStatus> = jobs_map.values().collect();
    
    Ok(warp::reply::json(&jobs_vec))
}
