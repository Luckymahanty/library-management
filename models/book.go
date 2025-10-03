package models

type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    Status string `json:"status"` // available, borrowed
}
