package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// Simple test to examine raw Solr response for highlighting
func main() {
	// Get Solr URL from environment variables or use default
	solrBaseURL := os.Getenv("SOLR_BASE_URL")
	if solrBaseURL == "" {
		solrBaseURL = "http://localhost:8984"
	}

	// Test query with special characters
	testQuery := "boot():"

	// Set up the Solr URL with highlighting parameters
	solrURL := fmt.Sprintf("%s/solr/heline/select", solrBaseURL)
	u, _ := url.Parse(solrURL)
	q := u.Query()
	q.Set("hl", "on")
	q.Set("hl.fl", "content")
	q.Set("hl.simple.pre", "<mark>")
	q.Set("hl.simple.post", "</mark>")
	q.Set("hl.snippets", "3")
	q.Set("hl.usePhraseHighlighter", "true")
	q.Set("hl.requireFieldMatch", "true")
	q.Set("hl.highlightMultiTerm", "true")
	q.Set("hl.mergeContiguous", "true")
	q.Set("hl.fragsize", "2500")
	q.Set("hl.maxAnalyzedChars", "100000")
	q.Set("hl.method", "unified")
	
	// Escape special characters in the query
	escapedQuery := testQuery
	// Set the highlight query parameter explicitly
	q.Set("hl.q", fmt.Sprintf("content:\"%s\"", escapedQuery))
	q.Set("hl.qparser", "lucene")
	u.RawQuery = q.Encode()

	// Construct the JSON payload
	data := map[string]interface{}{
		"query":  fmt.Sprintf("content:\"%s\"", testQuery),
		"fields": "id,file_id,repo,lang,branch,owner_id",
	}

	queryData, _ := json.Marshal(data)
	fmt.Println("Query:", string(queryData))
	fmt.Println("URL:", u.String())

	// Send the request
	payload := bytes.NewReader(queryData)
	req, _ := http.NewRequest("POST", u.String(), payload)
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	}

	// Read and print the response
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// Pretty print the JSON for better readability
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "  ")
	if err == nil {
		fmt.Println("\n==== RAW SOLR RESPONSE ====")
		fmt.Println(prettyJSON.String())
		fmt.Println("==== END RAW SOLR RESPONSE ====")
	} else {
		fmt.Println("Error formatting JSON:", err)
		fmt.Println(string(body))
	}
}
