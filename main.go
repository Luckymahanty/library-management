package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/luckymahanty/library-management/models"
	"golang.org/x/crypto/bcrypt"
)

var books []models.Book
var nextID = 1

// ==========================
// 📘 BOOK HANDLERS
// ==========================
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	book.ID = nextID
	nextID++
	book.Status = "available"
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			fmt.Fprintf(w, "Book with ID %d deleted", id)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

// ==========================
// 👤 USER HANDLERS
// ==========================
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // "admin" or "user"
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Insert user into DB
	_, err = models.DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)",
		creds.Username, string(hashedPassword), creds.Role)
	if err != nil {
		http.Error(w, "User already exists or DB error", http.StatusConflict)
		return
	}

	fmt.Fprintf(w, "✅ User %s registered successfully!", creds.Username)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var storedHash string
	var role string
	err := models.DB.QueryRow("SELECT password, role FROM users WHERE username = ?", creds.Username).Scan(&storedHash, &role)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password)); err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Login successful!",
		"role":    role,
	})
}

// ==========================
// 🚀 MAIN FUNCTION
// ==========================
func main() {
	models.InitDB() // Initialize database

	// Book APIs
	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/add", addBook)
	http.HandleFunc("/delete", deleteBook)

	// User APIs
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)

	// Serve static files (optional frontend)
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	fmt.Println("📚 Database initialized successfully with books & users tables!")
	fmt.Println("🚀 Server running on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
