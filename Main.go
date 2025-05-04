package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	// Registering handlers
	log.Printf("Registering HTTP handlers")
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/create", createEventHandler)
	http.HandleFunc("/events", getAllEventsHandler)
	http.HandleFunc("/event/", getEventByIDHandler)

	// Starting the server
	log.Printf("Starting the server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Welcome to the Event Management System!"}`)
}
