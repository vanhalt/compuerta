package authservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

// AuthorizationRules holds a list of authorization rules.
// It is used to parse the rules from a YAML file.
type AuthorizationRules struct {
	// Rules is a slice of individual authorization rules.
	Rules []Rule `yaml:"rules"`
}

// Rule defines a single authorization rule.
// It specifies which roles are allowed to access a resource with certain methods.
type Rule struct {
	// Resource is the identifier for the resource (e.g., a URL path).
	Resource string `yaml:"resource"`
	// AllowedMethods is a list of HTTP methods permitted for this resource under this rule.
	AllowedMethods []string `yaml:"allowed_methods"`
	// Roles is a list of user roles that this rule applies to.
	Roles []string `yaml:"roles"`
}

// AuthorizationRequest represents the JSON request body for an authorization check.
type AuthorizationRequest struct {
	// Resource is the resource the user is trying to access.
	Resource string `json:"resource"`
	// Method is the HTTP method being used for access.
	Method string `json:"method"`
	// Role is the role of the user making the request.
	Role string `json:"role"`
}

// AuthorizationResponse represents the JSON response body for an authorization check.
type AuthorizationResponse struct {
	// Authorized indicates whether the request is permitted or not.
	Authorized bool `json:"authorized"`
}

var rules AuthorizationRules

// LoadRules reads the authorization rules from the specified YAML file.
// It populates the global 'rules' variable (of type AuthorizationRules) with the parsed rules.
// The file should be in YAML format and match the structure of AuthorizationRules.
func LoadRules(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading rules file: %v", err)
	}

	err = yaml.Unmarshal(data, &rules)
	if err != nil {
		return fmt.Errorf("error unmarshalling rules: %v", err)
	}
	return nil
}

// NewRouter creates and configures a new HTTP router.
// It sets up the `/authorize` endpoint to be handled by AuthorizeHandler for POST requests.
// It returns an http.Handler which can be used to start an HTTP server.
func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/authorize", AuthorizeHandler).Methods("POST")
	return r
}

// AuthorizeHandler handles incoming HTTP requests to the /authorize endpoint.
// It decodes the JSON request body into an AuthorizationRequest,
// checks authorization using IsAuthorized, and sends an AuthorizationResponse.
// If the request body is invalid, it returns an HTTP 400 Bad Request error.
func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	var authRequest AuthorizationRequest
	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	authorized := IsAuthorized(authRequest)

	response := AuthorizationResponse{Authorized: authorized}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// IsAuthorized checks if the given AuthorizationRequest is permitted based on the loaded rules.
// It iterates through the rules stored in the global 'rules' variable.
// A request is authorized if a rule matches the request's resource, method, and role.
func IsAuthorized(req AuthorizationRequest) bool {
	for _, rule := range rules.Rules {
		if rule.Resource == req.Resource {
			methodAllowed := false
			for _, allowedMethod := range rule.AllowedMethods {
				if allowedMethod == req.Method {
					methodAllowed = true
					break
				}
			}

			if methodAllowed {
				for _, allowedRole := range rule.Roles {
					if allowedRole == req.Role {
						return true
					}
				}
			}
		}
	}
	return false
}
