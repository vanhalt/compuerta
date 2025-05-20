package main

import (
	"log"
	"net/http"

	"vanhalt.com/authservice"
)

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
