package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/luckymahanty/library-management/models"
)

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := models.DB.Query("SELECT id, title, author, status FROM books")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		rows.Scan(&b.ID, &b.Title, &b.Author, &b.Status)
		books = append(books, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Add a new book
func addBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	book.Status = "available"
	result, err := models.DB.Exec("INSERT INTO books (title, author, status) VALUES (?, ?, ?)", book.Title, book.Author, book.Status)
	if err != nil {
		http.Error(w, "Database insert error", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	book.ID = int(id)
	json.NewEncoder(w).Encode(book)
}

// Delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := models.DB.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Database delete error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Book with ID %d deleted", id)
}

// Main function
func main() {
	models.InitDB() // ✅ Initialize SQLite database

	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/add", addBook)
	http.HandleFunc("/delete", deleteBook)

	// Serve static frontend files
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

