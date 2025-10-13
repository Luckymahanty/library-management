package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Transaction struct {
	Username string `json:"username"`
	BookID   int    `json:"book_id"`
	Action   string `json:"action"`
	Date     string `json:"date"`
}

func BorrowBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t Transaction
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Check if book is available
		var available int
		err := db.QueryRow("SELECT available FROM books WHERE id = ?", t.BookID).Scan(&available)
		if err != nil || available == 0 {
			http.Error(w, "Book not available", http.StatusBadRequest)
			return
		}

		// Mark book as borrowed
		_, err = db.Exec("UPDATE books SET available = 0 WHERE id = ?", t.BookID)
		if err != nil {
			http.Error(w, "Error updating book", http.StatusInternalServerError)
			return
		}

		// Record transaction
		_, err = db.Exec("INSERT INTO transactions (username, book_id, action, date) VALUES (?, ?, 'borrow', ?)",
			t.Username, t.BookID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			http.Error(w, "Error recording transaction", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Book borrowed successfully"))
	}
}

func ReturnBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t Transaction
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Mark book as returned
		_, err := db.Exec("UPDATE books SET available = 1 WHERE id = ?", t.BookID)
		if err != nil {
			http.Error(w, "Error updating book", http.StatusInternalServerError)
			return
		}

		// Record transaction
		_, err = db.Exec("INSERT INTO transactions (username, book_id, action, date) VALUES (?, ?, 'return', ?)",
			t.Username, t.BookID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			http.Error(w, "Error recording transaction", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Book returned successfully"))
	}
}

