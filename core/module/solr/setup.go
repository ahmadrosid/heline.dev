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

	// Create code_syntax field type for better handling of code patterns
	codeSyntaxFieldType := map[string]interface{}{
		"add-field-type": map[string]interface{}{
			"name":                    "code_syntax",
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
					"class": "solr.ClassicTokenizerFactory",
				},
				"filters": []map[string]interface{}{
					{
						"class": "solr.LowerCaseFilterFactory",
					},
					{
						"class": "solr.ShingleFilterFactory",
						"minShingleSize": "2",
						"maxShingleSize": "5",
						"outputUnigrams": "true",
					},
					{
						"class": "solr.RemoveDuplicatesTokenFilterFactory",
					},
				},
			},
			"query": map[string]interface{}{
				"charFilters": []map[string]interface{}{
					{
						"class": "solr.HTMLStripCharFilterFactory",
					},
				},
				"tokenizer": map[string]interface{}{
					"class": "solr.ClassicTokenizerFactory",
				},
				"filters": []map[string]interface{}{
					{
						"class": "solr.LowerCaseFilterFactory",
					},
				},
			},
		},
	}

	// Create text_ngram field type for partial matching
	textNgramFieldType := map[string]interface{}{
		"add-field-type": map[string]interface{}{
			"name":                    "text_ngram",
			"class":                   "solr.TextField",
			"positionIncrementGap":    "100",
			"analyzer": map[string]interface{}{
				"charFilters": []map[string]interface{}{
					{
						"class": "solr.PatternReplaceCharFilterFactory",
						"pattern": "([\\p{Punct}&&[^_]])",
						"replacement": " $1 ",
					},
				},
				"tokenizer": map[string]interface{}{
					"class": "solr.NGramTokenizerFactory",
					"minGramSize": "2",
					"maxGramSize": "15",
				},
				"filters": []map[string]interface{}{
					{
						"class": "solr.LowerCaseFilterFactory",
					},
				},
			},
			"query": map[string]interface{}{
				"charFilters": []map[string]interface{}{
					{
						"class": "solr.PatternReplaceCharFilterFactory",
						"pattern": "([\\p{Punct}&&[^_]])",
						"replacement": " $1 ",
					},
				},
				"tokenizer": map[string]interface{}{
					"class": "solr.StandardTokenizerFactory",
				},
				"filters": []map[string]interface{}{
					{
						"class": "solr.LowerCaseFilterFactory",
					},
				},
			},
		},
	}

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
					{
						"class": "solr.PatternReplaceCharFilterFactory",
						"pattern": "([\\p{Punct}&&[^_]])",
						"replacement": " $1 ",
					},
				},
				"tokenizer": map[string]interface{}{
					"class": "solr.WhitespaceTokenizerFactory",
					"rule":  "java",
				},
				"filters": []map[string]interface{}{
					{
						"class": "solr.WordDelimiterFilterFactory",
						"generateWordParts": "1",
						"generateNumberParts": "1",
						"catenateWords": "1",
						"catenateNumbers": "1",
						"catenateAll": "0",
						"splitOnCaseChange": "1",
						"preserveOriginal": "1",
					},
					{
						"class": "solr.LowerCaseFilterFactory",
					},
					{
						"class": "solr.ASCIIFoldingFilterFactory",
					},
					{
						"class": "solr.StopFilterFactory",
						"ignoreCase": "true",
						"words": "stopwords.txt",
					},
				},
			},
			"query": map[string]interface{}{
				"charFilters": []map[string]interface{}{
					{
						"class": "solr.HTMLStripCharFilterFactory",
					},
					{
						"class": "solr.PatternReplaceCharFilterFactory",
						"pattern": "([\\p{Punct}&&[^_]])",
						"replacement": " $1 ",
					},
				},
				"tokenizer": map[string]interface{}{
					"class": "solr.WhitespaceTokenizerFactory",
					"rule":  "java",
				},
				"filters": []map[string]interface{}{
					{
						"class": "solr.WordDelimiterFilterFactory",
						"generateWordParts": "1",
						"generateNumberParts": "1",
						"catenateWords": "1",
						"catenateNumbers": "1",
						"catenateAll": "0",
						"splitOnCaseChange": "1",
						"preserveOriginal": "1",
					},
					{
						"class": "solr.LowerCaseFilterFactory",
					},
					{
						"class": "solr.ASCIIFoldingFilterFactory",
					},
					{
						"class": "solr.StopFilterFactory",
						"ignoreCase": "true",
						"words": "stopwords.txt",
					},
				},
			},
		},
	}

	codeSyntaxFieldTypeJSON, err := json.Marshal(codeSyntaxFieldType)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/solr/heline/schema", solrBaseURL), bytes.NewBuffer(codeSyntaxFieldTypeJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp1, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp1.Body.Close()

	textNgramFieldTypeJSON, err := json.Marshal(textNgramFieldType)
	if err != nil {
		return err
	}

	req, err = http.NewRequest("POST", fmt.Sprintf("%s/solr/heline/schema", solrBaseURL), bytes.NewBuffer(textNgramFieldTypeJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp2, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	// Marshal and send the text_html field type
	fieldTypeJSON, err := json.Marshal(fieldTypeData)
	if err != nil {
		return err
	}

	req, err = http.NewRequest("POST", fmt.Sprintf("%s/solr/heline/schema", solrBaseURL), bytes.NewBuffer(fieldTypeJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp3, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp3.Body.Close()

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
			{
				"name":        "code_content",
				"type":        "code_syntax",
				"multiValued": true,
				"stored":      true,
				"indexed":     true,
			},
			{
				"name":        "identifier_ngram",
				"type":        "text_ngram",
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

	resp4, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp4.Body.Close()

	return nil
}
