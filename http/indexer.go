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

// NewIndexerClient creates a new client for the heline-indexer API
func NewIndexerClient() *IndexerClient {
	// Get the indexer URL from environment or use default
	indexerURL := os.Getenv("INDEXER_URL")
	if indexerURL == "" {
		indexerURL = "http://heline-indexer:8080"
	}

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

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

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
