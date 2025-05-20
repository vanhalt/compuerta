package authservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

type AuthorizationRules struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Resource       string   `yaml:"resource"`
	AllowedMethods []string `yaml:"allowed_methods"`
	Roles          []string `yaml:"roles"`
}

type AuthorizationRequest struct {
	Resource string `json:"resource"`
	Method   string `json:"method"`
	Role     string `json:"role"`
}

type AuthorizationResponse struct {
	Authorized bool `json:"authorized"`
}

var rules AuthorizationRules

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

func NewRouter() http.Handler { return mux.NewRouter() } // Placeholder

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
