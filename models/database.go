package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the SQLite database and creates required tables
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./library.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create Books table
	createBookTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		author TEXT,
		status TEXT
	);`

	// Create Users table
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		role TEXT DEFAULT 'user'
	);`

	// Execute table creation queries
	_, err = DB.Exec(createBookTable)
	if err != nil {
		log.Fatal("Error creating books table:", err)
	}

	_, err = DB.Exec(createUserTable)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	// Add a default admin if not exists
	addDefaultAdmin()

	fmt.Println("📚 Database initialized successfully!")
}

// GetDB returns the global DB connection
func GetDB() *sql.DB {
	return DB
}

// Add a default admin user if not already present
func addDefaultAdmin() {
	query := `INSERT OR IGNORE INTO users (username, password, role) VALUES ('admin', 'admin123', 'admin');`
	_, err := DB.Exec(query)
	if err != nil {
		log.Println("⚠️ Failed to add default admin:", err)
	}
}

