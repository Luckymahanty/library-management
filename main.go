package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    _ "github.com/mattn/go-sqlite3"
    "github.com/luckymahanty/library-management/models"
)

var db *sql.DB

func main() {
    db = models.InitDB() // ✅ properly initialize DB
     models.DB = db
    fmt.Println("📚 Database ready with users & books!")

    fs := http.FileServer(http.Dir("./frontend"))
    http.Handle("/", fs)

    http.HandleFunc("/signup", signupHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/books", models.GetBooksHandler)

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
		author TEXT
	);`

	if _, err := db.Exec(userTable); err != nil {
		log.Fatal("❌ Failed to create users table:", err)
	}
	if _, err := db.Exec(bookTable); err != nil {
		log.Fatal("❌ Failed to create books table:", err)
	}
}

// ---------------- HANDLERS ---------------

func signupHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var user struct {
        Username string `json:"username"`
        Password string `json:"password"`
        Role     string `json:"role"`
    }

    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    _, err := db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", user.Username, user.Password, user.Role)
    if err != nil {
        http.Error(w, "User already exists or database error", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "Signup successful"}`))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var dbUser models.User
	err := db.QueryRow("SELECT id, username, password, role FROM users WHERE username = ?", creds.Username).
		Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password, &dbUser.Role)
	if err != nil {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}

	if dbUser.Password != creds.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// send user info for redirect
	resp, _ := json.Marshal(dbUser)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, author FROM books")
	if err != nil {
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		rows.Scan(&b.ID, &b.Title, &b.Author)
		books = append(books, b)
	}

	json.NewEncoder(w).Encode(books)
}

func addBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var b models.Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", b.Title, b.Author)
	if err != nil {
		http.Error(w, "Failed to add book", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Book added successfully"))
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Book deleted"))
}
