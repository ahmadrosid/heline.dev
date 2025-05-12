package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmadrosid/heline/core/entity"
	"github.com/ahmadrosid/heline/core/module/solr"
)

// MockResetIndex is a mock function for solr.ResetIndex
type MockResetIndex func(recreateSchema bool) error

// TestHandleResetIndex tests the handleResetIndex function
func TestHandleResetIndex(t *testing.T) {
	// Save the original ResetIndex function and restore it after the test
	originalResetIndex := solr.ResetIndex
	defer func() {
		solr.ResetIndex = originalResetIndex
	}()

	// Test cases
	testCases := []struct {
		name           string
		method         string
		requestBody    map[string]interface{}
		mockResetIndex MockResetIndex
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Simple reset with POST",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"recreate_schema": false,
			},
			mockResetIndex: func(recreateSchema bool) error {
				if recreateSchema != false {
					t.Errorf("Expected recreateSchema to be false, got %v", recreateSchema)
				}
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"message": "Index reset successful",
				"details": map[string]interface{}{
					"recreate_schema": false,
				},
			},
		},
		{
			name:   "Full reset with schema recreation",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"recreate_schema": true,
			},
			mockResetIndex: func(recreateSchema bool) error {
				if recreateSchema != true {
					t.Errorf("Expected recreateSchema to be true, got %v", recreateSchema)
				}
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"message": "Index reset successful",
				"details": map[string]interface{}{
					"recreate_schema": true,
				},
			},
		},
		{
			name:        "Empty POST request (default values)",
			method:      http.MethodPost,
			requestBody: nil,
			mockResetIndex: func(recreateSchema bool) error {
				if recreateSchema != false {
					t.Errorf("Expected recreateSchema to be false, got %v", recreateSchema)
				}
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"message": "Index reset successful",
				"details": map[string]interface{}{
					"recreate_schema": false,
				},
			},
		},
		{
			name:           "Method not allowed",
			method:         http.MethodGet,
			requestBody:    nil,
			mockResetIndex: nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody: map[string]interface{}{
				"error": "Method not allowed. Use POST.",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the mock ResetIndex function if provided
			if tc.mockResetIndex != nil {
				solr.ResetIndex = tc.mockResetIndex
			}

			// Create a request
			var reqBody []byte
			var err error
			if tc.requestBody != nil {
				reqBody, err = json.Marshal(tc.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req, err := http.NewRequest(tc.method, "/api/index/reset", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handleResetIndex(rr, req)

			// Check the status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Check the response body
			var responseBody map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			// Compare the response body with the expected body
			for key, expectedValue := range tc.expectedBody {
				actualValue, ok := responseBody[key]
				if !ok {
					t.Errorf("Expected key %s in response body, but it was not found", key)
					continue
				}

				// Special handling for nested maps
				if expectedMap, ok := expectedValue.(map[string]interface{}); ok {
					if actualMap, ok := actualValue.(map[string]interface{}); ok {
						for nestedKey, nestedExpectedValue := range expectedMap {
							nestedActualValue, ok := actualMap[nestedKey]
							if !ok {
								t.Errorf("Expected nested key %s in response body, but it was not found", nestedKey)
								continue
							}
							if nestedActualValue != nestedExpectedValue {
								t.Errorf("Expected nested value %v for key %s, got %v", nestedExpectedValue, nestedKey, nestedActualValue)
							}
						}
					} else {
						t.Errorf("Expected map for key %s, got %T", key, actualValue)
					}
				} else if actualValue != expectedValue {
					t.Errorf("Expected value %v for key %s, got %v", expectedValue, key, actualValue)
				}
			}
		})
	}
}

// TestIntegrationWithHandler tests the integration of the reset handler with the main HTTP handler
func TestIntegrationWithHandler(t *testing.T) {
	// Save the original ResetIndex function and restore it after the test
	originalResetIndex := solr.ResetIndex
	defer func() {
		solr.ResetIndex = originalResetIndex
	}()

	// Mock the ResetIndex function
	resetCalled := false
	solr.ResetIndex = func(recreateSchema bool) error {
		resetCalled = true
		return nil
	}

	// Create a test server with our handler
	handler := Handler(nil)
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create a client
	client := &http.Client{}

	// Create a request
	reqBody, _ := json.Marshal(map[string]interface{}{
		"recreate_schema": true,
	})
	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/index/reset", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check if ResetIndex was called
	if !resetCalled {
		t.Errorf("Expected ResetIndex to be called, but it was not")
	}

	// Check the response body
	var responseBody entity.Map
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Verify the response contains the expected fields
	if success, ok := responseBody["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true, got %v", responseBody["success"])
	}

	if message, ok := responseBody["message"].(string); !ok || message != "Index reset successful" {
		t.Errorf("Expected message to be 'Index reset successful', got %v", responseBody["message"])
	}
}
