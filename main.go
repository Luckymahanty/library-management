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

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
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
	book.ID = nextID
	nextID++
	book.Status = "available"
	books = append(books, book)
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
	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			fmt.Fprintf(w, "Book with ID %d deleted", id)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/add", addBook)
	http.HandleFunc("/delete", deleteBook)

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

