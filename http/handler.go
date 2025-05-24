package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/ahmadrosid/heline/core/entity"
	"github.com/ahmadrosid/heline/core/module/solr"
	"github.com/ahmadrosid/heline/core/utils"
	queryparam "github.com/tomwright/queryparam/v4"
)

func Handler() http.Handler {

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Welcome to Heline API",
		})
	}))
	mux.HandleFunc("/api/search", handleSearch)
	
	// Add indexer API endpoints
	mux.HandleFunc("/api/index", handleIndexRepository)
	mux.HandleFunc("/api/index/status/", handleJobStatus)
	mux.HandleFunc("/api/index/jobs", handleListJobs)
	
	// Add index management endpoints
	mux.HandleFunc("/api/index/reset", handleResetIndex)

	return wrapCORSHandler(mux, &CorsConfig{
		allowedOrigin: "*",
	})
}

func getQueryFilter(param entity.QueryParam) []string {
	var filter []string

	if len(param.Lang) > 0 {
		filter = append(filter, fmt.Sprintf("lang:(%s)", utils.Join(param.Lang, " ", "*")))
	}

	if len(param.Path) > 0 {
		filter = append(filter, fmt.Sprintf("path:(%s)", utils.Join(param.Path, " ", "*")))
	}

	if len(param.Repo) > 0 {
		filter = append(filter, fmt.Sprintf("repo:(%s)", utils.Join(param.Repo, " ", "*")))
	}

	return filter
}

// handleIndexRepository processes requests to index a git repository
func handleIndexRepository(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	// Parse the request body
	var req struct {
		GitURL string `json:"git_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate the git URL
	if req.GitURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Git URL is required",
		})
		return
	}

	// Create a new indexer client
	client := NewIndexerClient()

	// Send the indexing request to the indexer API
	resp, err := client.IndexRepository(req.GitURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to index repository: " + err.Error(),
		})
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleJobStatus retrieves the status of an indexing job
func handleJobStatus(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	// Extract the job ID from the URL path
	jobID := strings.TrimPrefix(r.URL.Path, "/api/index/status/")
	if jobID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Job ID is required",
		})
		return
	}

	// Create a new indexer client
	client := NewIndexerClient()

	// Get the job status from the indexer API
	status, err := client.GetJobStatus(jobID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get job status: " + err.Error(),
		})
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// handleListJobs retrieves a list of all indexing jobs
func handleListJobs(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	// Create a new indexer client
	client := NewIndexerClient()

	// Get the list of jobs from the indexer API
	jobs, err := client.ListJobs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to list jobs: " + err.Error(),
		})
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// enhanceHighlighting improves the highlighting of code patterns with special characters
func enhanceHighlighting(result []byte, originalQuery string) []byte {
	// Parse the JSON result
	var data map[string]interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		// If we can't parse the JSON, return the original result
		return result
	}

	// Check if highlighting data exists
	highlighting, ok := data["highlighting"].(map[string]interface{})
	if !ok {
		return result
	}

	// Prepare the pattern to be highlighted
	// Escape special regex characters
	escapedQuery := regexp.QuoteMeta(originalQuery)

	// Process each document's highlighting
	for docID, highlightData := range highlighting {
		docHighlight, ok := highlightData.(map[string]interface{})
		if !ok {
			continue
		}

		// Process content field highlighting
		contentHighlights, ok := docHighlight["content"].([]interface{})
		if !ok {
			continue
		}

		// Process each highlighted snippet
		for i, snippet := range contentHighlights {
			snippetStr, ok := snippet.(string)
			if !ok {
				continue
			}

			// Apply custom highlighting to ensure the full pattern is highlighted
			// Simple case-insensitive match for the escaped query
			re := regexp.MustCompile(fmt.Sprintf("(?i)(%s)", escapedQuery))
			enhancedSnippet := re.ReplaceAllString(snippetStr, "<mark>$1</mark>")

			// Update the snippet in the data
			contentHighlights[i] = enhancedSnippet
		}

		// Update the content highlights
		docHighlight["content"] = contentHighlights
		highlighting[docID] = docHighlight
	}

	// Update the highlighting in the data
	data["highlighting"] = highlighting

	// Convert back to JSON
	enhancedResult, err := json.Marshal(data)
	if err != nil {
		// If we can't convert back to JSON, return the original result
		return result
	}

	return enhancedResult
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	param := entity.QueryParam{}
	err := queryparam.Parse(r.URL.Query(), &param)
	switch err {
	case nil:
		break
	case queryparam.ErrInvalidBoolValue:
		println("Failed parse query param")
		return
	default:
		println("return empty query param", err.Error())
		return
	}

	q := strings.Replace(param.Query, "*", "\\*", -1)
	q = strings.Replace(q, "\"", "\\\"", -1)
	q = strings.Replace(q, "'", "\\'", -1)
	// q = "*" + q + "*"
	println(q)

	// Original search query for post-processing highlighting
	originalQuery := q

	result, err := solr.Search(solr.SolrQuery{
		Query:  q,
		Filter: getQueryFilter(param),
	})
	
	// Post-process the result to improve highlighting if needed
	if err == nil && len(originalQuery) > 0 && strings.ContainsAny(originalQuery, ":(){}[]") {
		// For queries with special characters, we may need to enhance the highlighting
		result = enhanceHighlighting(result, originalQuery)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(entity.Map{
			"error": err.Error(),
		})
		return
	}

	dec := json.NewDecoder(bytes.NewReader(result))
	var data entity.SolrResult
	err = dec.Decode(&data)
	if err != nil {
		enc.Encode(entity.Map{
			"error": err.Error(),
		})
		return
	}

	println("hints:", data.Response.NumFound, q)

	w.Header().Set("Content-Type", "application/json")
	var content []entity.ContentData
	for _, item := range data.Response.Docs {
		contents := data.Highlight[item.ID].Content
		if len(contents) == 0 {
			continue
		}
		content = append(content, entity.ContentData{
			ID: entity.Map{
				"raw": item.ID,
			},
			Branch: entity.Map{
				"raw": item.Branch,
			},
			OwnerID: entity.Map{
				"raw": item.OwnerID,
			},
			FileID: entity.Map{
				"raw": item.FileID,
			},
			Content: entity.Map{
				"snippet": contents,
				// "snippet": contents[len(contents)-1:],
			},
			Repo: entity.Map{
				"raw": item.Repo,
			},
		})
	}
	enc.Encode(entity.CodeSearchResult{
		Response: entity.CodeHits{
			Hits:   content,
			Facets: data.Facet,
			Total:  data.Response.NumFound,
		},
	})
}
