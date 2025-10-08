package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the database
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./library.db")
	if err != nil {
		log.Fatal(err)
	}

	createBookTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		author TEXT,
		status TEXT
	);`

	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		role TEXT
	);`

	_, err = DB.Exec(createBookTable)
	if err != nil {
		log.Fatal("Error creating books table:", err)
	}

	_, err = DB.Exec(createUserTable)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	fmt.Println("📚 Database initialized successfully with books & users tables!")
}
