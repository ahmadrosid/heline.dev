package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// IndexerClient handles communication with the heline-indexer API service
type IndexerClient struct {
	BaseURL string
	Client  *http.Client
}

// IndexRequest represents a request to index a git repository
type IndexRequest struct {
	GitURL string `json:"git_url"`
}

// IndexResponse represents the response from the indexer API
type IndexResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	JobID   *string `json:"job_id,omitempty"`
}

// JobStatus represents the status of an indexing job
type JobStatus struct {
	GitURL      string     `json:"git_url"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Message     *string    `json:"message,omitempty"`
}

func NewIndexerClient() *IndexerClient {
	indexerURL := os.Getenv("INDEXER_URL")
	if indexerURL == "" {
		if _, err := os.Stat("/app"); os.IsNotExist(err) {
			indexerURL = "http://localhost:8080"
		} else {
			indexerURL = "http://heline-indexer:8080"
		}
	}

	fmt.Printf("Connecting to indexer at: %s\n", indexerURL)

	return &IndexerClient{
		BaseURL: indexerURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// IndexRepository sends a request to index a git repository
func (c *IndexerClient) IndexRepository(gitURL string) (*IndexResponse, error) {
	reqBody := IndexRequest{
		GitURL: gitURL,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/index", c.BaseURL), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var indexResp IndexResponse
	if err := json.NewDecoder(resp.Body).Decode(&indexResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &indexResp, nil
}

// GetJobStatus retrieves the status of an indexing job
func (c *IndexerClient) GetJobStatus(jobID string) (*JobStatus, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/jobs/%s", c.BaseURL, jobID), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var jobStatus JobStatus
	if err := json.NewDecoder(resp.Body).Decode(&jobStatus); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &jobStatus, nil
}

// ListJobs retrieves a list of all indexing jobs
func (c *IndexerClient) ListJobs() ([]JobStatus, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/jobs", c.BaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "heline-app/1.0")
	
	resp, err := c.Client.Do(req)
	if err != nil {
		// If we can't connect, try a fallback URL if in Docker
		fmt.Printf("Error connecting to %s: %v\n", c.BaseURL, err)
		
		// For now, return a mock empty response instead of an error
		// This allows the UI to work even when the indexer is down
		fmt.Println("Returning empty jobs list as fallback")
		return []JobStatus{}, nil
	}
	defer resp.Body.Close()
	
	fmt.Printf("Received response with status code: %d\n", resp.StatusCode)
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var jobs []JobStatus
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return jobs, nil
}
