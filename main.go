package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/luckymahanty/library-management/models"
)

var db *sql.DB

func main() {
	db = models.InitDB()
	models.DB = db

	fmt.Println("📚 Database ready with users & books!")

	// static files (frontend)
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	// API
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/books", booksHandler)         // GET list
	http.HandleFunc("/addbook", addBookHandler)     // POST add
	http.HandleFunc("/deletebook", deleteBookHandler) // DELETE ?id=#

	// start
	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ---- handlers ----

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if u.Username == "" || u.Password == "" {
		http.Error(w, "username/password required", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("INSERT INTO users (username,password,role) VALUES (?, ?, ?)", u.Username, u.Password, u.Role)
	if err != nil {
		http.Error(w, "user exists or db error", http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "signup successful"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	var u models.User
	err := db.QueryRow("SELECT id,username,password,role FROM users WHERE username=?", creds.Username).
		Scan(&u.ID, &u.Username, &u.Password, &u.Role)
	if err == sql.ErrNoRows || u.Password != creds.Password {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	// return user object (frontend will redirect by role)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		books, err := models.GetAllBooks()
		if err != nil {
			http.Error(w, "failed to fetch books", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func addBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}
	var b models.Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if b.Title == "" || b.Author == "" {
		http.Error(w, "title/author required", http.StatusBadRequest)
		return
	}
	if err := models.AddBook(b.Title, b.Author); err != nil {
		http.Error(w, "failed to add book", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "book added"})
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	// allow POST fallback from forms, but prefer DELETE
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := models.DeleteBook(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "book deleted"})
}
