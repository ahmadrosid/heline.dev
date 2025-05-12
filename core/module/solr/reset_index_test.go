package solr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestResetIndex tests the ResetIndex function
func TestResetIndex(t *testing.T) {
	// Create a mock HTTP server to simulate Solr
	mockSolr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the request path and method
		path := r.URL.Path
		method := r.Method
		query := r.URL.Query()

		// Simulate different Solr endpoints
		switch {
		// Delete all documents endpoint
		case path == "/solr/heline/update" && method == "POST":
			// Check if this is a delete request
			var deleteReq map[string]interface{}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&deleteReq); err != nil {
				t.Logf("Error decoding delete request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Verify it's a delete all query
			if deleteQuery, ok := deleteReq["delete"].(map[string]interface{}); ok {
				if query, ok := deleteQuery["query"].(string); ok && query == "*:*" {
					w.WriteHeader(http.StatusOK)
					fmt.Fprintln(w, `{"responseHeader":{"status":0,"QTime":5}}`)
					return
				}
			}
			w.WriteHeader(http.StatusBadRequest)

		// Unload core endpoint
		case path == "/solr/admin/cores" && query.Get("action") == "UNLOAD":
			if query.Get("core") == "heline" {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, `{"responseHeader":{"status":0,"QTime":10}}`)
				return
			}
			w.WriteHeader(http.StatusBadRequest)

		// Create core endpoint
		case path == "/solr/admin/cores" && query.Get("action") == "CREATE":
			if query.Get("name") == "heline" || query.Get("name") == "docset" {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, `{"responseHeader":{"status":0,"QTime":100}}`)
				return
			}
			w.WriteHeader(http.StatusBadRequest)

		// Check core status endpoint
		case path == "/solr/admin/cores" && query.Get("action") == "STATUS":
			coreName := query.Get("core")
			if coreName == "heline" || coreName == "docset" {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{"responseHeader":{"status":0,"QTime":1},"status":{"%s":{"name":"%s","instanceDir":"path/to/%s"}}}`, coreName, coreName, coreName)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"responseHeader":{"status":0,"QTime":1},"status":{}}`)

		// Schema field check endpoint
		case path == "/solr/heline/schema/fields/content":
			// Return 404 to simulate field not existing (to trigger schema creation)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, `{"responseHeader":{"status":404,"QTime":5}}`)

		// Schema update endpoint
		case path == "/solr/heline/schema":
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"responseHeader":{"status":0,"QTime":10}}`)

		default:
			t.Logf("Unexpected request to %s with method %s", path, method)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockSolr.Close()

	// Set the mock server URL as the Solr base URL
	os.Setenv("SOLR_BASE_URL", mockSolr.URL)
	defer os.Unsetenv("SOLR_BASE_URL")

	// Test cases
	testCases := []struct {
		name           string
		recreateSchema bool
	}{
		{
			name:           "Reset index without schema recreation",
			recreateSchema: false,
		},
		{
			name:           "Reset index with schema recreation",
			recreateSchema: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ResetIndex(tc.recreateSchema)
			if err != nil {
				t.Errorf("ResetIndex(%v) failed: %v", tc.recreateSchema, err)
			}
		})
	}
}

// TestDeleteAllDocuments tests the deleteAllDocuments function
func TestDeleteAllDocuments(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is to the correct endpoint
		if r.URL.Path != "/solr/heline/update" {
			t.Errorf("Expected request to /solr/heline/update, got %s", r.URL.Path)
		}

		// Check if the request method is POST
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check if the content type is correct
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Check if the request body contains the delete query
		var requestBody map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}

		// Check if the delete query is correct
		deleteQuery, ok := requestBody["delete"].(map[string]interface{})
		if !ok {
			t.Errorf("Expected delete query in request body, got %v", requestBody)
		}

		query, ok := deleteQuery["query"].(string)
		if !ok || query != "*:*" {
			t.Errorf("Expected delete query to be *:*, got %v", query)
		}

		// Return a successful response
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"responseHeader":{"status":0,"QTime":5}}`)
	}))
	defer mockServer.Close()

	// Call the function with the mock server URL
	err := deleteAllDocuments(mockServer.URL)
	if err != nil {
		t.Errorf("deleteAllDocuments failed: %v", err)
	}
}

// TestUnloadCore tests the unloadCore function
func TestUnloadCore(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is to the correct endpoint
		if r.URL.Path != "/solr/admin/cores" {
			t.Errorf("Expected request to /solr/admin/cores, got %s", r.URL.Path)
		}

		// Check if the request method is GET
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Check if the query parameters are correct
		query := r.URL.Query()
		if query.Get("action") != "UNLOAD" {
			t.Errorf("Expected action=UNLOAD, got %s", query.Get("action"))
		}
		if query.Get("core") != "heline" {
			t.Errorf("Expected core=heline, got %s", query.Get("core"))
		}
		if query.Get("deleteIndex") != "true" {
			t.Errorf("Expected deleteIndex=true, got %s", query.Get("deleteIndex"))
		}
		if query.Get("deleteDataDir") != "true" {
			t.Errorf("Expected deleteDataDir=true, got %s", query.Get("deleteDataDir"))
		}
		if query.Get("deleteInstanceDir") != "true" {
			t.Errorf("Expected deleteInstanceDir=true, got %s", query.Get("deleteInstanceDir"))
		}

		// Return a successful response
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"responseHeader":{"status":0,"QTime":10}}`)
	}))
	defer mockServer.Close()

	// Call the function with the mock server URL
	err := unloadCore(mockServer.URL)
	if err != nil {
		t.Errorf("unloadCore failed: %v", err)
	}
}
