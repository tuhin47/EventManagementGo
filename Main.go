package main

import (
	"EventManagement/router"
	"log"
	"net/http"
)

func main() {
	// Registering handlers
	router.RegisterRoutes()

	// Starting the server
	log.Printf("Starting the server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
