package routes

import (
	db "backend/pkg/db/sqlite"
	"encoding/json"
	"net/http"
)

type User struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DOB       string `json:"dob"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	AboutME   string `json:"about_me"`
}

// AUTH HANDLERS
// POST /auth/register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	json.NewDecoder(r.Body).Decode(&user)

	// NOTE: db processes to be handled through functions
	// NOTE: below query indicates wrong schema (to be checked)
	// TODO: below query to be amended to conform to the regiteration requirement (handle all user properties)
	_, err := db.DB.Exec("INSERT INTO users (username, password, bio) VALUES (?, ?, ?)", user.Email, user.Password, user.AboutME)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// POST /auth/login
func LoginHandler(w http.ResponseWriter, r *http.Request) {

}

// POST /auth/logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {

}
