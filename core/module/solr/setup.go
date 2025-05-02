package solr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// SetupSchema checks if the Solr schema is properly set up and creates it if not
func SetupSchema() error {
	// Get Solr URL from environment variables or use default
	solrBaseURL := os.Getenv("SOLR_BASE_URL")
	if solrBaseURL == "" {
		solrBaseURL = "http://localhost:8984"
	}

	fmt.Println("🔍 Checking Solr schema setup...")

	// First, check if cores exist and create them if they don't
	if err := createCores(solrBaseURL); err != nil {
		return fmt.Errorf("failed to create Solr cores: %w", err)
	}

	// Then set up the schema for the heline core
	if err := setupHelineSchema(solrBaseURL); err != nil {
		return fmt.Errorf("failed to set up heline schema: %w", err)
	}

	fmt.Println("✅ Solr schema setup complete!")
	return nil
}

// createCores creates the necessary Solr cores if they don't exist
func createCores(solrBaseURL string) error {
	// Check if heline core exists
	resp, err := http.Get(fmt.Sprintf("%s/solr/admin/cores?action=STATUS&core=heline", solrBaseURL))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var statusResp map[string]interface{}
	if err := json.Unmarshal(body, &statusResp); err != nil {
		return err
	}

	// Create heline core if it doesn't exist
	if status, ok := statusResp["status"].(map[string]interface{}); !ok || status["heline"] == nil {
		fmt.Println("Creating heline core...")
		createURL := fmt.Sprintf("%s/solr/admin/cores?action=CREATE&name=heline&instanceDir=heline&config=solrconfig.xml&dataDir=data", solrBaseURL)
		if _, err := http.Get(createURL); err != nil {
			return err
		}
	}

	// Check if docset core exists
	resp, err = http.Get(fmt.Sprintf("%s/solr/admin/cores?action=STATUS&core=docset", solrBaseURL))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &statusResp); err != nil {
		return err
	}

	// Create docset core if it doesn't exist
	if status, ok := statusResp["status"].(map[string]interface{}); !ok || status["docset"] == nil {
		fmt.Println("Creating docset core...")
		createURL := fmt.Sprintf("%s/solr/admin/cores?action=CREATE&name=docset&instanceDir=docset&config=solrconfig.xml&dataDir=data", solrBaseURL)
		if _, err := http.Get(createURL); err != nil {
			return err
		}
	}

	return nil
}

// setupHelineSchema sets up the schema for the heline core
func setupHelineSchema(solrBaseURL string) error {
	// Check if the schema is already set up by checking for a field
	resp, err := http.Get(fmt.Sprintf("%s/solr/heline/schema/fields/content", solrBaseURL))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// If the field exists, we assume the schema is set up
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Solr schema already set up.")
		return nil
	}

	fmt.Println("Setting up Solr schema...")

	// Create text_html field type
	fieldTypeData := map[string]interface{}{
		"add-field-type": map[string]interface{}{
			"name":                    "text_html",
			"class":                   "solr.TextField",
			"positionIncrementGap":    "100",
			"autoGeneratePhraseQueries": "true",
			"analyzer": map[string]interface{}{
				"charFilters": []map[string]interface{}{
					{
						"class": "solr.HTMLStripCharFilterFactory",
					},
				},
				"tokenizer": map[string]interface{}{
					"class": "solr.WhitespaceTokenizerFactory",
					"rule":  "java",
				},
				"filters": []map[string]interface{}{
					{
						"class": "solr.WordDelimiterFilterFactory",
					},
					{
						"class": "solr.LowerCaseFilterFactory",
					},
					{
						"class": "solr.ASCIIFoldingFilterFactory",
					},
				},
			},
			"query": map[string]interface{}{
				"tokenizer": map[string]interface{}{
					"class": "solr.WhitespaceTokenizerFactory",
					"rule":  "java",
				},
				"filters": []map[string]interface{}{
					{
						"class": "solr.WordDelimiterFilterFactory",
					},
					{
						"class": "solr.LowerCaseFilterFactory",
					},
					{
						"class": "solr.ASCIIFoldingFilterFactory",
					},
				},
			},
		},
	}

	fieldTypeJSON, err := json.Marshal(fieldTypeData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/solr/heline/schema", solrBaseURL), bytes.NewBuffer(fieldTypeJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create schema fields
	fieldsData := map[string]interface{}{
		"add-field": []map[string]interface{}{
			{
				"name":   "branch",
				"type":   "string",
				"stored": true,
			},
			{
				"name":   "path",
				"type":   "string",
				"stored": true,
			},
			{
				"name":   "file_id",
				"type":   "string",
				"stored": true,
			},
			{
				"name":   "owner_id",
				"type":   "string",
				"stored": true,
			},
			{
				"name":   "lang",
				"type":   "string",
				"stored": true,
			},
			{
				"name":   "repo",
				"type":   "string",
				"stored": true,
			},
			{
				"name":        "content",
				"type":        "text_html",
				"multiValued": true,
				"stored":      true,
				"indexed":     true,
			},
		},
	}

	fieldsJSON, err := json.Marshal(fieldsData)
	if err != nil {
		return err
	}

	req, err = http.NewRequest("POST", fmt.Sprintf("%s/solr/heline/schema", solrBaseURL), bytes.NewBuffer(fieldsJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
