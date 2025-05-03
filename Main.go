package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	// Updated DSN to include parseTime and loc parameters
	dsn := "root:root@tcp(127.0.0.1:3306)/event_management?parseTime=true&loc=Local"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create events table if it doesn't exist
	createTableQuery := `CREATE TABLE IF NOT EXISTS events (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		location VARCHAR(255),
		start_time DATETIME,
		end_time DATETIME,
		organizer VARCHAR(255),
		capacity INT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Welcome to the Event Management System!"}`)
	})

	// Registering handlers
	http.HandleFunc("/create", createEventHandler)
	http.HandleFunc("/events", getAllEventsHandler)

	// Starting the server
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
