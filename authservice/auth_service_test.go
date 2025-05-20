package authservice

import (
	// Assuming your auth service package is named 'authservice'
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Declare variables and functions needed for testing that are likely defined in auth_service.go
// This is a placeholder and should be replaced with the actual declarations or by
// ensuring these are exported from the main package and imported correctly.
var rulesFile string

func setupRouter() http.Handler { return nil } // Placeholder

func createTemporaryRulesFile(t *testing.T) {
	// Save the original rules file path
	originalRulesFile := rulesFile
	defer func() { rulesFile = originalRulesFile }()

	// Create a temporary rules file for testing
	tempRulesFile, err := os.CreateTemp("", "rules-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary rules file: %v", err)
	}
	defer os.Remove(tempRulesFile.Name())
	rulesFile = tempRulesFile.Name()

	// Write some sample rules to the temporary file
	sampleRules := `
	rules:
	  - resource: /users
		allowed_methods:
		  - GET
		  - POST
		roles:
		  - admin
		  - user
	  - resource: /admin
		allowed_methods:
		  - GET
		roles:
		  - admin
	`
	if _, err := tempRulesFile.WriteString(sampleRules); err != nil {
		t.Fatalf("Failed to write sample rules to file: %v", err)
	}
	tempRulesFile.Close()
}

func TestAuthorizeHandler(t *testing.T) {
	createTemporaryRulesFile(t)

	// Reload rules after creating the temporary file
	if err := LoadRules(rulesFile); err != nil {
		t.Fatalf("Failed to load rules from temporary file: %v", err)
	}

	router := setupRouter()

	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
		expectedAuth   bool
	}{
		{
			name: "Successful Authorization - GET /users as admin",
			payload: map[string]string{
				"resource": "/users",
				"method":   "GET",
				"role":     "admin",
			},
			expectedStatus: http.StatusOK,
			expectedAuth:   true,
		},
		{
			name: "Successful Authorization - POST /users as user",
			payload: map[string]string{
				"resource": "/users",
				"method":   "POST",
				"role":     "user",
			},
			expectedStatus: http.StatusOK,
			expectedAuth:   true,
		},
		{
			name: "Unauthorized - Incorrect method",
			payload: map[string]string{
				"resource": "/users",
				"method":   "PUT",
				"role":     "admin",
			},
			expectedStatus: http.StatusOK,
			expectedAuth:   false,
		},
		{
			name: "Unauthorized - Incorrect role",
			payload: map[string]string{
				"resource": "/users",
				"method":   "GET",
				"role":     "guest",
			},
			expectedStatus: http.StatusOK,
			expectedAuth:   false,
		},
		{
			name: "Unauthorized - Unknown resource",
			payload: map[string]string{
				"resource": "/products",
				"method":   "GET",
				"role":     "admin",
			},
			expectedStatus: http.StatusOK,
			expectedAuth:   false,
		},
		{
			name: "Successful Authorization - GET /admin as admin",
			payload: map[string]string{
				"resource": "/admin",
				"method":   "GET",
				"role":     "admin",
			},
			expectedStatus: http.StatusOK,
			expectedAuth:   true,
		},
		{
			name: "Unauthorized - POST /admin as admin",
			payload: map[string]string{
				"resource": "/admin",
				"method":   "POST",
				"role":     "admin",
			},
			expectedStatus: http.StatusOK,
			expectedAuth:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatalf("Failed to marshal payload: %v", err)
			}

			req, err := http.NewRequest("POST", "/authorize", bytes.NewBuffer(payloadBytes))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			var response struct {
				Authorized bool `json:"authorized"`
			}
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			if response.Authorized != tt.expectedAuth {
				t.Errorf("handler returned unexpected authorization status: got %v want %v",
					response.Authorized, tt.expectedAuth)
			}
		})
	}
}

func TestAuthorizeHandler_InvalidJson(t *testing.T) {
	createTemporaryRulesFile(t)
	router := setupRouter()

	req, err := http.NewRequest("POST", "/authorize", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid JSON: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestAuthorizeHandler_MissingRulesFile(t *testing.T) {
	createTemporaryRulesFile(t)
	// Save the original rules file path
	originalRulesFile := rulesFile
	defer func() { rulesFile = originalRulesFile }()

	// Set a non-existent rules file path
	rulesFile = "non_existent_rules.yaml"

	// Reload rules, which should fail
	if err := LoadRules(rulesFile); err == nil {
		t.Fatal("loadRules unexpectedly succeeded with a non-existent file")
	}

	router := setupRouter()

	payload := map[string]string{
		"resource": "/users",
		"method":   "GET",
		"role":     "admin",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", "/authorize", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// The server should return an internal server error because rules failed to load
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code when rules file is missing: got %v want %v",
			status, http.StatusInternalServerError)
	}

	// Restore the original rules file path for other tests
	rulesFile = originalRulesFile
	// Attempt to reload the original rules, this might fail if the original file is also missing,
	// but it's necessary to clean up for other tests.
	if err := LoadRules(rulesFile); err != nil {
		fmt.Printf("Warning: Failed to reload original rules file after missing file test: %v\n", err)
	}
}
