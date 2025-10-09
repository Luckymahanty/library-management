package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var db *sql.DB

// Initialize database
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./library.db")
	if err != nil {
		log.Fatal(err)
	}

	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		role TEXT
	);`

	createBookTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		author TEXT,
		status TEXT
	);`

	_, err = db.Exec(createUserTable)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}
	_, err = db.Exec(createBookTable)
	if err != nil {
		log.Fatal("Error creating books table:", err)
	}

	fmt.Println("📚 Database initialized successfully with books & users tables!")
}

// --- Signup handler ---
func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)",
		user.Username, user.Password, user.Role)
	if err != nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "✅ User %s registered successfully!", user.Username)
}

// --- Login handler ---
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var dbUser User
	err := db.QueryRow("SELECT id, username, password, role FROM users WHERE username = ?", input.Username).
		Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password, &dbUser.Role)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if input.Password != dbUser.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "✅ Welcome, %s! You are logged in as %s.", dbUser.Username, dbUser.Role)
}

// --- Main function ---
func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)

	// Serve frontend
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
