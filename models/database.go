package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB opens/creates the SQLite DB and creates tables if missing.
// Returns the *sql.DB so main can hold it too.
func InitDB() *sql.DB {
	var err error
	DB, err = sql.Open("sqlite3", "./library.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// Create tables
	createUsers := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		role TEXT
	);`
	createBooks := `CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		author TEXT,
		status TEXT DEFAULT 'available'
	);`
	createTx := `CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		book_id INTEGER,
		action TEXT,
		date TEXT
	);`

	if _, err := DB.Exec(createUsers); err != nil {
		log.Fatalf("create users table: %v", err)
	}
	if _, err := DB.Exec(createBooks); err != nil {
		log.Fatalf("create books table: %v", err)
	}
	if _, err := DB.Exec(createTx); err != nil {
		log.Fatalf("create transactions table: %v", err)
	}

	// ensure an admin exists (username: admin, password: admin123)
	_, _ = DB.Exec(`INSERT OR IGNORE INTO users (username,password,role) VALUES ('admin','admin123','admin')`)

	fmt.Println("✅ Database initialized successfully")
	return DB
}
