package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/luckymahanty/library-management/models"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// ✅ Signup handler
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	db := models.GetDB()
	_, err := db.Exec(`INSERT INTO users (username, password, role) VALUES (?, ?, ?)`,
		user.Username, user.Password, user.Role)
	if err != nil {
		http.Error(w, "User already exists or DB error", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("✅ User created successfully"))
}

// ✅ Signin handler
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	db := models.GetDB()
	var dbPassword, role string
	err := db.QueryRow(`SELECT password, role FROM users WHERE username = ?`, creds.Username).Scan(&dbPassword, &role)
	if err == sql.ErrNoRows || dbPassword != creds.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println("DB error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"message": "✅ Login successful", "role": role}
	json.NewEncoder(w).Encode(resp)
}

