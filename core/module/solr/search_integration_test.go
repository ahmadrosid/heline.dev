package solr

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// Simple test to verify search functionality works with real Solr data
func TestSearchIntegration(t *testing.T) {
	// Skip test if Solr is not available
	if !isSolrAvailable() {
		t.Skip("Skipping test: Solr server not available")
	}

	testCases := []struct {
		name          string
		query         string
		expectResults bool
		description   string
	}{
		{
			name:          "Search for colon void pattern",
			query:         ": void",
			expectResults: true,
			description:   "Should find TypeScript/JavaScript void patterns",
		},
		{
			name:          "Search for function pattern",
			query:         "function(",
			expectResults: true,
			description:   "Should find function declarations",
		},
		{
			name:          "Search for simple term",
			query:         "function",
			expectResults: true,
			description:   "Should find the word function",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create search query
			query := SolrQuery{
				Query: tc.query,
			}

			t.Logf("Testing query: '%s' - %s", tc.query, tc.description)

			// Execute search
			result, err := Search(query)
			if err != nil {
				t.Fatalf("Search failed for query '%s': %v", tc.query, err)
			}

			// Parse JSON response to check basic structure
			var solrResp map[string]interface{}
			if err := json.Unmarshal(result, &solrResp); err != nil {
				t.Fatalf("Failed to parse Solr response for query '%s': %v", tc.query, err)
			}

			// Check response header
			responseHeader, ok := solrResp["responseHeader"].(map[string]interface{})
			if !ok {
				t.Fatalf("Invalid response header for query '%s'", tc.query)
			}

			status, ok := responseHeader["status"].(float64)
			if !ok {
				t.Fatalf("Invalid status in response header for query '%s'", tc.query)
			}

			if status != 0 {
				t.Errorf("Solr returned error status %.0f for query '%s'", status, tc.query)
				// Print response for debugging
				t.Logf("Response: %s", string(result)[:min(500, len(result))])
				return
			}

			// Check response body
			response, ok := solrResp["response"].(map[string]interface{})
			if !ok {
				t.Fatalf("Invalid response body for query '%s'", tc.query)
			}

			numFound, ok := response["numFound"].(float64)
			if !ok {
				t.Fatalf("Invalid numFound in response for query '%s'", tc.query)
			}

			hasResults := numFound > 0
			
			if tc.expectResults && !hasResults {
				t.Errorf("Expected results for query '%s' but got none. NumFound: %.0f", tc.query, numFound)
				return
			}

			if hasResults {
				t.Logf("✅ Query '%s' returned %.0f results", tc.query, numFound)

				// Check if highlighting exists
				highlighting, hasHighlighting := solrResp["highlighting"].(map[string]interface{})
				if hasHighlighting && len(highlighting) > 0 {
					highlightCount := 0
					markTagCount := 0
					
					for docID, highlight := range highlighting {
						if highlightMap, ok := highlight.(map[string]interface{}); ok {
							if content, hasContent := highlightMap["content"].([]interface{}); hasContent {
								highlightCount++
								// Check for <mark> tags in highlighting
								for _, contentItem := range content {
									if contentStr, ok := contentItem.(string); ok {
										if strings.Contains(contentStr, "<mark>") {
											markTagCount++
											// Log a sample
											if markTagCount <= 2 {
												sample := contentStr
												if len(sample) > 100 {
													sample = sample[:100] + "..."
												}
												t.Logf("Found <mark> tags in doc %s: %s", docID, sample)
											}
										}
									}
								}
							}
						}
					}

					if markTagCount > 0 {
						t.Logf("✅ Query '%s' has <mark> tags in %d out of %d highlighted documents", tc.query, markTagCount, highlightCount)
					} else {
						t.Logf("⚠️  Query '%s' has highlighting in %d documents but no <mark> tags found", tc.query, highlightCount)
					}
				} else {
					t.Logf("⚠️  Query '%s' has no highlighting data", tc.query)
				}
			}
		})
	}
}

// Helper function to check if Solr is available
func isSolrAvailable() bool {
	// Try to get the SOLR_BASE_URL or use default
	solrBaseURL := os.Getenv("SOLR_BASE_URL")
	if solrBaseURL == "" {
		solrBaseURL = "http://localhost:8984"
	}

	// Create a simple test query
	query := SolrQuery{
		Query: "test",
	}

	// Try to execute search
	_, err := Search(query)
	return err == nil
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}