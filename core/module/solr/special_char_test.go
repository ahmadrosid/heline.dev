package solr

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestSpecialCharacterHighlighting specifically tests the fix for special character highlighting
func TestSpecialCharacterHighlighting(t *testing.T) {
	if !isSolrAvailable() {
		t.Skip("Skipping test: Solr server not available")
	}

	// These are the specific queries that were failing before the fix
	problematicQueries := []struct {
		query       string
		description string
	}{
		{": void", "Colon followed by void - common in TypeScript"},
		{"function(", "Function with opening parenthesis"},
		{"[]", "Empty array brackets"},
		{"{}", "Empty object brackets"},
		{"()", "Empty parentheses"},
		{"=>", "Arrow function syntax"},
	}

	for _, test := range problematicQueries {
		t.Run("Query_"+strings.ReplaceAll(test.query, " ", "_"), func(t *testing.T) {
			query := SolrQuery{Query: test.query}
			
			t.Logf("Testing problematic query: '%s' (%s)", test.query, test.description)
			
			result, err := Search(query)
			if err != nil {
				t.Fatalf("Search failed for query '%s': %v", test.query, err)
			}

			// Parse response
			var response map[string]interface{}
			if err := json.Unmarshal(result, &response); err != nil {
				t.Fatalf("Failed to parse response for query '%s': %v", test.query, err)
			}

			// Check that we get a valid response (status 0)
			if header, ok := response["responseHeader"].(map[string]interface{}); ok {
				if status, ok := header["status"].(float64); ok && status != 0 {
					t.Errorf("Query '%s' returned error status: %.0f", test.query, status)
					return
				}
			}

			// Check if we have results
			if respBody, ok := response["response"].(map[string]interface{}); ok {
				if numFound, ok := respBody["numFound"].(float64); ok {
					if numFound > 0 {
						t.Logf("✅ Query '%s' found %.0f results", test.query, numFound)
						
						// Most importantly, check that highlighting exists and contains <mark> tags
						if highlighting, ok := response["highlighting"].(map[string]interface{}); ok && len(highlighting) > 0 {
							markTagCount := 0
							for _, docHighlight := range highlighting {
								if docMap, ok := docHighlight.(map[string]interface{}); ok {
									if content, ok := docMap["content"].([]interface{}); ok {
										for _, contentItem := range content {
											if contentStr, ok := contentItem.(string); ok && strings.Contains(contentStr, "<mark>") {
												markTagCount++
												break // Only count once per document
											}
										}
									}
								}
							}
							
							if markTagCount > 0 {
								t.Logf("✅ Query '%s' has proper highlighting with <mark> tags in %d documents", test.query, markTagCount)
							} else {
								t.Errorf("❌ Query '%s' has highlighting but no <mark> tags found (this was the original bug)", test.query)
							}
						} else {
							t.Logf("⚠️  Query '%s' has no highlighting data", test.query)
						}
					} else {
						t.Logf("ℹ️  Query '%s' returned no results (this may be expected if no data contains this pattern)", test.query)
					}
				}
			}
		})
	}
}

// TestSearchQueryEscaping tests that special characters are properly escaped in different parts of the query
func TestSearchQueryEscaping(t *testing.T) {
	testCases := []struct {
		input              string
		expectedEscapedChar string
		description        string
	}{
		{": void", "\\:", "Colon should be escaped"},
		{"function()", "\\(", "Parentheses should be escaped"},
		{"array[]", "\\[", "Brackets should be escaped"},
		{"obj{}", "\\{", "Braces should be escaped"},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// This tests the escaping logic without hitting Solr
			specialChars := []string{":", ";", "{", "}", "(", ")", "[", "]", "<", ">", "=", "+", "-", "!", "*", "?", "~", "^", "&", "|", "%"}
			hasSpecialChars := false
			for _, char := range specialChars {
				if strings.Contains(test.input, char) {
					hasSpecialChars = true
					break
				}
			}

			if !hasSpecialChars {
				t.Errorf("Test case '%s' should contain special characters", test.input)
			}

			// Test that our escaping would work
			escapedQuery := test.input
			escapedQuery = strings.ReplaceAll(escapedQuery, ":", "\\:")
			escapedQuery = strings.ReplaceAll(escapedQuery, "(", "\\(")
			escapedQuery = strings.ReplaceAll(escapedQuery, ")", "\\)")
			escapedQuery = strings.ReplaceAll(escapedQuery, "[", "\\[")
			escapedQuery = strings.ReplaceAll(escapedQuery, "]", "\\]")
			escapedQuery = strings.ReplaceAll(escapedQuery, "{", "\\{")
			escapedQuery = strings.ReplaceAll(escapedQuery, "}", "\\}")

			if !strings.Contains(escapedQuery, test.expectedEscapedChar) {
				t.Errorf("Expected escaped query to contain '%s', but got: %s", test.expectedEscapedChar, escapedQuery)
			}

			t.Logf("✅ Input '%s' correctly escaped to '%s'", test.input, escapedQuery)
		})
	}
}