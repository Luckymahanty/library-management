package models

import (
	
	"fmt"
)

// Book represents a library book
type Book struct {
	ID       int
	Title    string
	Author   string
	Status   string
}

// Transaction represents a borrow or return record
type Transaction struct {
	ID       int
	Username string
	BookID   int
	Action   string
	Date     string
}

// AddBook adds a new book to the database
func AddBook(title, author string) error {
	query := `INSERT INTO books (title, author, status) VALUES (?, ?, 'available');`
	_, err := DB.Exec(query, title, author)
	if err != nil {
		return fmt.Errorf("failed to add book: %v", err)
	}
	return nil
}

// GetAllBooks returns a list of all books
func GetAllBooks() ([]Book, error) {
	rows, err := DB.Query(`SELECT id, title, author, status FROM books;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		err = rows.Scan(&b.ID, &b.Title, &b.Author, &b.Status)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

// RecordTransaction logs a book borrow/return
func RecordTransaction(username string, bookID int, action, date string) error {
	query := `INSERT INTO transactions (username, book_id, action, date) VALUES (?, ?, ?, ?);`
	_, err := DB.Exec(query, username, bookID, action, date)
	if err != nil {
		return fmt.Errorf("failed to record transaction: %v", err)
	}
	return nil
}
