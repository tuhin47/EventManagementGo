package handler

import (
	"EventManagement/config"
	"EventManagement/domain"
	"EventManagement/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	db       *sql.DB             = config.GetDB()
	validate *validator.Validate = validator.New()
)

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var event domain.Event

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

	if err := validate.Struct(event); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO events (title, description, start_time, end_time, created_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, event.Title, event.Description, event.StartTime, event.EndTime, event.CreatedBy, time.Now(), time.Now())
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

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response to JSON: %v", err)
		return
	}

	log.Printf("Event created successfully with ID: %d", id)
}

func GetAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get pagination parameters from utility function
	pageSize, offset := utils.GetPaginationParams(r)

	// TODO : Response should contain the total number of records
	rows, err := db.Query(`SELECT id, title, description, start_time, end_time, created_by, created_at, updated_at FROM events LIMIT ? OFFSET ?`, pageSize, offset)
	if err != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		log.Printf("Error querying database: %v", err)
		return
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.CreatedBy, &event.CreatedAt, &event.UpdatedAt); err != nil {
			http.Error(w, "Failed to parse events", http.StatusInternalServerError)
			log.Printf("Error scanning row: %v", err)
			return
		}

		log.Printf("Scanned Event: %+v", event)

		events = append(events, event)
	}

	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode events", http.StatusInternalServerError)
		log.Printf("Error encoding events to JSON: %v", err)
		return
	}

	log.Println("Fetched list of all events with pagination")
}

func GetEventByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.ExtractIDFromURL(r, "/event/")
	if err != nil {
		http.Error(w, "Missing event ID", http.StatusBadRequest)
		return
	}

	var event domain.Event
	query := `SELECT id, title, description, start_time, end_time, created_by, created_at, updated_at FROM events WHERE id = ?`
	err = db.QueryRow(query, id).Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.CreatedBy, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Event not found", http.StatusNotFound)
		} else {
			log.Printf("Error querying database: %v", err)
			http.Error(w, "Failed to fetch event", http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, "Failed to encode event", http.StatusInternalServerError)
		log.Printf("Error encoding event to JSON: %v", err)
		return
	}

	log.Printf("Fetched event with ID: %s", id)
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.ExtractIDFromURL(r, "/event/")
	if err != nil {
		http.Error(w, "Missing event ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	log.Printf("Raw request body: %s", string(body))

	if err := json.Unmarshal(body, &updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	if len(updates) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	query := "UPDATE events SET "
	args := []interface{}{}
	for field, value := range updates {
		query += field + " = ?, "
		args = append(args, value)
	}
	query += "updated_at = ?, "
	args = append(args, time.Now())

	query = query[:len(query)-2] // Remove the trailing comma and space
	query += " WHERE id = ?"
	args = append(args, id)

	_, err = db.Exec(query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Event updated successfully!",
		"id":      id,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response to JSON: %v", err)
		return
	}

	log.Printf("Event updated successfully with ID: %s", id)
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.ExtractIDFromURL(r, "/event/")
	if err != nil {
		http.Error(w, "Missing event ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM events WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching affected rows: %v", err)
		http.Error(w, "Failed to retrieve affected rows", http.StatusInternalServerError)
		return
	}

	if affectedRows == 0 {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"message": "Event deleted successfully!",
		"id":      id,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response to JSON: %v", err)
		return
	}

	log.Printf("Event deleted successfully with ID: %s", id)
}
