package authservice

import (
	// Assuming your auth service package is named 'authservice'
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Declare variables and functions needed for testing that are likely defined in auth_service.go
// This is a placeholder and should be replaced with the actual declarations or by
// ensuring these are exported from the main package and imported correctly.
var rulesFile string

func setupRouter() http.Handler { return NewRouter() }

func createTemporaryRulesFile(t *testing.T) {
	// Save the original rules file path
	// originalRulesFile := rulesFile // This logic will be moved to the calling test

	// Create a temporary rules file for testing
	tempRulesFile, err := os.CreateTemp("", "rules-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary rules file: %v", err)
	}
	rulesFile = tempRulesFile.Name()

	// Write some sample rules to the temporary file
	sampleRules := `rules:
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
	originalRulesFile := rulesFile // Save global state
	createTemporaryRulesFile(t)    // This will set the global rulesFile
	defer func() {
		os.Remove(rulesFile)       // Remove the temp file
		rulesFile = originalRulesFile // Restore global state
	}()

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
	originalRulesFile := rulesFile // Save global state
	createTemporaryRulesFile(t)    // This will set the global rulesFile
	defer func() {
		os.Remove(rulesFile)       // Remove the temp file
		rulesFile = originalRulesFile // Restore global state
	}()
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
