package models

import (
	"fmt"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Status string `json:"status"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func AddBook(title, author string) error {
	q := `INSERT INTO books (title, author, status) VALUES (?, ?, 'available')`
	_, err := DB.Exec(q, title, author)
	if err != nil {
		return fmt.Errorf("AddBook: %v", err)
	}
	return nil
}

func GetAllBooks() ([]Book, error) {
	rows, err := DB.Query(`SELECT id, title, author, status FROM books ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Status); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

func DeleteBook(id int) error {
	_, err := DB.Exec(`DELETE FROM books WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("DeleteBook: %v", err)
	}
	return nil
}
