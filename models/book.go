package models

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

type Book struct {
	ID     int
	Title  string
	Author string
	Status string
}

type User struct {
	ID       int
	Username string
	Password string
	Role     string
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

// GetAllBooks returns all books
func GetAllBooks() ([]Book, error) {
	rows, err := DB.Query(`SELECT id, title, author, status FROM books;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Status)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

// DeleteBook removes a book by ID
func DeleteBook(id int) error {
	query := `DELETE FROM books WHERE id = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %v", err)
	}
	return nil
}
