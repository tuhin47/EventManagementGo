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
	http.HandleFunc("/events", handleEvents)
	http.HandleFunc("/event/", handleEvent)

	// Starting the server
	log.Printf("Starting the server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getEventByIDHandler(w, r)
	case http.MethodPut:
		updateEventHandler(w, r)
	case http.MethodDelete:
		deleteEventHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func handleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllEventsHandler(w, r)
	case http.MethodPost:
		createEventHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Welcome to the Event Management System!"}`)
}
