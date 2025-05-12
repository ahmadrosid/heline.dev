package solr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// ResetIndex completely resets the Solr index by:
// 1. Deleting all documents
// 2. Optionally recreating the schema
func ResetIndex(recreateSchema bool) error {
	// Get Solr URL from environment variables or use default
	solrBaseURL := os.Getenv("SOLR_BASE_URL")
	if solrBaseURL == "" {
		solrBaseURL = "http://localhost:8984"
	}

	fmt.Println("🧹 Resetting Solr index...")

	// Step 1: Delete all documents
	if err := deleteAllDocuments(solrBaseURL); err != nil {
		return fmt.Errorf("failed to delete all documents: %w", err)
	}

	// Step 2: Unload the core (optional)
	if recreateSchema {
		if err := unloadCore(solrBaseURL); err != nil {
			return fmt.Errorf("failed to unload core: %w", err)
		}

		// Step 3: Recreate the core and schema
		if err := createCores(solrBaseURL); err != nil {
			return fmt.Errorf("failed to create Solr cores: %w", err)
		}

		if err := setupHelineSchema(solrBaseURL); err != nil {
			return fmt.Errorf("failed to set up heline schema: %w", err)
		}
	}

	fmt.Println("✅ Solr index reset complete!")
	return nil
}

// deleteAllDocuments removes all documents from the Solr index
func deleteAllDocuments(solrBaseURL string) error {
	fmt.Println("Deleting all documents from index...")
	
	// Construct the delete-all query
	deleteQuery := map[string]interface{}{
		"delete": map[string]string{
			"query": "*:*",
		},
		"commit": map[string]interface{}{},
	}
	
	deleteJSON, err := json.Marshal(deleteQuery)
	if err != nil {
		return err
	}
	
	// Send the delete request
	url := fmt.Sprintf("%s/solr/heline/update", solrBaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(deleteJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Delete response:", string(body))
	
	return nil
}

// unloadCore unloads the Solr core
func unloadCore(solrBaseURL string) error {
	fmt.Println("Unloading Solr core...")
	
	// Construct the unload URL with parameters to delete the data
	unloadURL := fmt.Sprintf("%s/solr/admin/cores?action=UNLOAD&core=heline&deleteIndex=true&deleteDataDir=true&deleteInstanceDir=true", solrBaseURL)
	
	resp, err := http.Get(unloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Unload response:", string(body))
	
	return nil
}
