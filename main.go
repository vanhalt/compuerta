// Package main is the entry point for the authorization microservice.
// It initializes and starts the HTTP server.
package main

import (
	"log"
	"net/http"

	"vanhalt.com/authservice"
)

// main is the primary function for the authorization service.
// It performs the following steps:
// 1. Loads authorization rules from "rules.yaml".
// 2. Creates a new HTTP router using the authservice package.
// 3. Starts an HTTP server on port 8080.
// If any step fails, it logs a fatal error and exits.
func main() {
	err := authservice.LoadRules("rules.yaml")
	if err != nil {
		log.Fatalf("Error loading rules: %v", err)
	}
	router := authservice.NewRouter() // Assuming authservice has a function NewRouter that returns a *mux.Router

	log.Println("Authorization service starting on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
