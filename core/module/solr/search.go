package solr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ahmadrosid/heline/core/entity"
	"github.com/ahmadrosid/heline/core/utils"
)

type SolrQuery struct {
	Query  string
	Filter []string
}

func Search(query SolrQuery) ([]byte, error) {
	// Get Solr URL from environment variables or use default
	solrBaseURL := os.Getenv("SOLR_BASE_URL")
	if solrBaseURL == "" {
		solrBaseURL = "http://localhost:8984"
	}
	
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
	q.Set("hl.regex.slop", "0.2")
	q.Set("hl.regex.pattern", ".*")
	q.Set("hl.bs.type", "WORD")
	q.Set("hl.bs.language", "en")
	u.RawQuery = q.Encode()

	// Check if the query contains special characters that might need phrase searching
	hasSpecialChars := false
	specialChars := []string{":", ";", "{", "}", "(", ")", "[", "]", "<", ">", "=", "+", "-", "!", "*", "?", "~", "^", "&", "|", "%"}
	for _, char := range specialChars {
		if strings.Contains(query.Query, char) {
			hasSpecialChars = true
			break
		}
	}

	// Construct the query based on content
	var solrQuery string
	var hlQuery string
	
	if hasSpecialChars {
		// For code patterns with special characters, we need to use a more precise approach
		// First, try an exact phrase match with proper escaping
		escapedQuery := strings.ReplaceAll(query.Query, ":", "\\:")
		escapedQuery = strings.ReplaceAll(escapedQuery, "(", "\\(")
		escapedQuery = strings.ReplaceAll(escapedQuery, ")", "\\)")
		escapedQuery = strings.ReplaceAll(escapedQuery, "[", "\\[")
		escapedQuery = strings.ReplaceAll(escapedQuery, "]", "\\]")
		
		// Build a query that searches for the exact phrase and also for parts of the phrase
		// This improves recall while still prioritizing exact matches
		solrQuery = fmt.Sprintf("content:\"%s\"^10", query.Query)
		
		// Add individual term matches with lower boost
		terms := strings.Fields(query.Query)
		for _, term := range terms {
			if len(term) > 1 { // Only add meaningful terms
				solrQuery += fmt.Sprintf(" OR content:%s^2", term)
			}
		}
		
		// Also search for the pattern without spaces
		noSpaceQuery := strings.ReplaceAll(query.Query, " ", "")
		if noSpaceQuery != query.Query {
			solrQuery += fmt.Sprintf(" OR content:%s^5", noSpaceQuery)
		}
		
		// For highlighting, use the exact phrase
		hlQuery = fmt.Sprintf("content:\"%s\"", escapedQuery)
	} else {
		// Use standard query for simple terms
		solrQuery = "content:" + query.Query
		hlQuery = solrQuery
	}
	
	// Set content field for highlighting
	q.Set("hl.fl", "content")
	
	// Set the highlight query parameter explicitly
	q.Set("hl.q", hlQuery)
	
	// Configure highlighting for better code pattern matching
	q.Set("hl.usePhraseHighlighter", "true")
	q.Set("hl.highlightMultiTerm", "true")
	q.Set("hl.tag.pre", "<mark>")
	q.Set("hl.tag.post", "</mark>")
	q.Set("hl.method", "unified")
	q.Set("hl.bs.type", "WORD")
	q.Set("hl.fragsize", "2500")
	q.Set("hl.maxAnalyzedChars", "500000")
	q.Set("hl.phraseLimit", "2000")
	q.Set("hl.multiValuedSeparatorChar", " ")
	
	// Debug info
	fmt.Println("\n==== SEARCH QUERY INFO ====")
	fmt.Println("Original query:", query.Query)
	fmt.Println("Solr query:", solrQuery)
	fmt.Println("Highlight query:", hlQuery)
	fmt.Println("==== END QUERY INFO ====\n")

	data := entity.Map{
		"query":  solrQuery,
		"fields": "id,file_id,repo,lang,branch,owner_id",
		"facet": entity.Map{
			"lang": entity.Map{
				"type":  "terms",
				"field": "lang",
				"limit": 10,
			},
			"path": entity.Map{
				"type":  "terms",
				"field": "path",
				"limit": 8,
			},
			"repo": entity.Map{
				"type":  "terms",
				"field": "repo",
				"limit": 7,
			},
		},
	}

	if query.Filter != nil {
		data["filter"] = query.Filter
	}

	queryData, _ := json.Marshal(data)

	println("search.go - solr_query:", utils.ByteToString(queryData))
	println("search.go - raw_query:", u.RawQuery)

	payload := bytes.NewReader(queryData)

	req, _ := http.NewRequest("POST", u.String(), payload)

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		println("ERROR", err.Error())
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// Debug: Print the query information
	fmt.Println("\n==== SOLR QUERY INFO ====")
	fmt.Println("Query:", solrQuery)
	fmt.Println("Highlight Query:", q.Get("hl.q"))
	fmt.Println("==== END QUERY INFO ====\n")

	return body, nil
}
