package http

import (
	"encoding/json"
	"net/http"

	"github.com/ahmadrosid/heline/core/entity"
	"github.com/ahmadrosid/heline/core/module/solr"
)

// handleResetIndex handles requests to reset the Solr index
func handleResetIndex(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.Map{
			"error": "Method not allowed. Use POST.",
		})
		return
	}

	// Parse request body
	var requestBody struct {
		RecreateSchema bool `json:"recreate_schema"`
	}

	// Default to false if not specified
	requestBody.RecreateSchema = false

	// Try to decode the request body into the struct
	if r.Body != nil {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			// If there is an error, just use the default value
			// This allows for empty POST requests to work
		}
	}

	// Reset the index
	err := solr.ResetIndex(requestBody.RecreateSchema)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Map{
			"error":   "Failed to reset index",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.Map{
		"success": true,
		"message": "Index reset successful",
		"details": map[string]interface{}{
			"recreate_schema": requestBody.RecreateSchema,
		},
	})
}
