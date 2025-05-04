package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func GetDB() *sql.DB {
	if DB == nil {
		InitDB()
	}
	return DB
}

func InitDB() {
	var err error
	log.Printf("Initializing database connection")
	dsn := "root:root@tcp(127.0.0.1:3306)/event_management?parseTime=true&loc=Local"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ping the database to ensure the connection is valid
	if err = DB.Ping(); err != nil {
		log.Fatalf("Database connection is not valid: %v", err)
	}

	log.Printf("Creating events table if it doesn't exist")
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
	if _, err := DB.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Printf("Database connection and initialization successful")
}
