package models

import (
    "database/sql"
    "log"


    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() *sql.DB {
    db, err := sql.Open("sqlite3", "./library.db")
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE,
        password TEXT,
        role TEXT
    );`)
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        author TEXT,
        available INTEGER
    );`)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("✅ Database initialized successfully")
    return db
}
