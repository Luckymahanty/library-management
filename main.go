package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"

    "github.com/luckymahanty/library-management/models"
)

var books []models.Book
var nextID = 1

// ==========================
// 📘 Existing book handlers
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
// 🧍‍♂️ NEW USER HANDLERS
// ==========================
func signupHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")

        if username == "" || password == "" {
            http.Error(w, "Username and password are required", http.StatusBadRequest)
            return
        }

        db := models.GetDB()
        _, err := db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, password)
        if err != nil {
            http.Error(w, "User already exists or database error", http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, "✅ Signup successful! You can now log in.")
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")

        db := models.GetDB()
        row := db.QueryRow("SELECT password FROM users WHERE username = ?", username)

        var storedPassword string
        err := row.Scan(&storedPassword)
        if err != nil {
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        }

        if password != storedPassword {
            http.Error(w, "Invalid password", http.StatusUnauthorized)
            return
        }

        fmt.Fprintf(w, "✅ Login successful! Welcome %s", username)
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
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

    // 👇 Register new user routes here
    http.HandleFunc("/signup", signupHandler)
    http.HandleFunc("/login", loginHandler)

    // Serve static files
    fs := http.FileServer(http.Dir("./frontend"))
    http.Handle("/", fs)

    fmt.Println("🚀 Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

