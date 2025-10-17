package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Status string `json:"status"`
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./library.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()

	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/addBook", addBookHandler)
	http.HandleFunc("/deleteBook", deleteBookHandler)
	http.HandleFunc("/borrowBook", borrowBookHandler)
	http.HandleFunc("/returnBook", returnBookHandler)

	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTables() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		role TEXT
	);`

	bookTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		author TEXT,
		status TEXT
	);`

	if _, err := db.Exec(userTable); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(bookTable); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Database initialized successfully")
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	_, err := db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)",
		user.Username, user.Password, user.Role)
	if err != nil {
		http.Error(w, "User already exists or DB error", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Signup successful"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	row := db.QueryRow("SELECT id, role FROM users WHERE username=? AND password=?", user.Username, user.Password)
	var id int
	var role string
	err := row.Scan(&id, &role)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "role": role})
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	rows, err := db.Query("SELECT id, title, author, status FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		rows.Scan(&b.ID, &b.Title, &b.Author, &b.Status)
		books = append(books, b)
	}
	json.NewEncoder(w).Encode(books)
}

func addBookHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	_, err := db.Exec("INSERT INTO books (title, author, status) VALUES (?, ?, ?)", book.Title, book.Author, "available")
	if err != nil {
		http.Error(w, "Error adding book", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Book added"})
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method != "DELETE" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	_, err := db.Exec("DELETE FROM books WHERE id=?", id)
	if err != nil {
		http.Error(w, "Error deleting book", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Book deleted"})
}

func borrowBookHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	_, err := db.Exec("UPDATE books SET status='borrowed' WHERE id=?", id)
	if err != nil {
		http.Error(w, "Error borrowing book", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Book borrowed"})
}

func returnBookHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	_, err := db.Exec("UPDATE books SET status='available' WHERE id=?", id)
	if err != nil {
		http.Error(w, "Error returning book", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Book returned"})
}
