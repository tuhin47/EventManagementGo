package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func GetDB() *sql.DB {
	var dbLock sync.Mutex

	dbLock.Lock()
	defer dbLock.Unlock()
	if DB == nil {
		InitDB()
	}
	return DB
}

func InitDB() {
	var err error
	log.Printf("Initializing database connection")

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	if host == "" || port == "" || username == "" || password == "" {
		log.Fatalf("Environment variables DB_HOST, DB_PORT, DB_USERNAME, or DB_PASSWORD are not set")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/event_management?parseTime=true&loc=Local", username, password, host, port)

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
		title VARCHAR(255) NOT NULL,
		description TEXT,
		start_time DATETIME,
		end_time DATETIME,
		created_by VARCHAR(255),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`
	if _, err := DB.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Printf("Database connection and initialization successful")
}
