package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./library.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		author TEXT,
		status TEXT
	);
	`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	log.Println("📚 Database initialized successfully!")
}
