package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type Event struct {
	ID          int       `json:"id"`          // Unique identifier for the event
	Name        string    `json:"name"`        // Name of the event
	Description string    `json:"description"` // Description of the event
	Location    string    `json:"location"`    // Location where the event will take place
	StartTime   time.Time `json:"start_time"`  // Start time of the event
	EndTime     time.Time `json:"end_time"`    // End time of the event
	Organizer   string    `json:"organizer"`   // Name of the organizer
	Capacity    int       `json:"capacity"`    // Maximum number of attendees
	Attendees   []string  `json:"attendees"`   // List of attendees
	CreatedAt   time.Time `json:"created_at"`  // Timestamp when the event was created
	UpdatedAt   time.Time `json:"updated_at"`  // Timestamp when the event was last updated
}

func createEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var event Event
	// Debugging: Log the raw request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	log.Printf("Raw request body: %s", string(body))

	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	query := `INSERT INTO events (name, description, location, start_time, end_time, organizer, capacity) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, event.Name, event.Description, event.Location, event.StartTime, event.EndTime, event.Organizer, event.Capacity)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error fetching last insert ID: %v", err)
		http.Error(w, "Failed to retrieve event ID", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Event created successfully!",
		"id":      id,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response to JSON: %v", err)
		return
	}

	log.Printf("Event created successfully with ID: %d", id)
}

func getAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`SELECT id, name, description, location, start_time, end_time, organizer, capacity, created_at, updated_at FROM events`)
	if err != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		log.Printf("Error querying database: %v", err)
		return
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.StartTime, &event.EndTime, &event.Organizer, &event.Capacity, &event.CreatedAt, &event.UpdatedAt); err != nil {
			http.Error(w, "Failed to parse events", http.StatusInternalServerError)
			log.Printf("Error scanning row: %v", err)
			return
		}

		// Debugging: Log the scanned values
		log.Printf("Scanned Event: %+v", event)

		events = append(events, event)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode events", http.StatusInternalServerError)
		log.Printf("Error encoding events to JSON: %v", err)
		return
	}

	log.Println("Fetched list of all events")
}

func getEventByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the event ID from the URL path
	id := r.URL.Path[len("/event/"):] // Assumes the path is /event/{id}
	if id == "" {
		http.Error(w, "Missing event ID", http.StatusBadRequest)
		return
	}

	var event Event
	query := `SELECT id, name, description, location, start_time, end_time, organizer, capacity, created_at, updated_at FROM events WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.StartTime, &event.EndTime, &event.Organizer, &event.Capacity, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Event not found", http.StatusNotFound)
		} else {
			log.Printf("Error querying database: %v", err)
			http.Error(w, "Failed to fetch event", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, "Failed to encode event", http.StatusInternalServerError)
		log.Printf("Error encoding event to JSON: %v", err)
		return
	}

	log.Printf("Fetched event with ID: %s", id)
}
