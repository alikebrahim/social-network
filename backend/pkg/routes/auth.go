package routes

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DOB       string `json:"date_of_birth"`
	Bio       string `json:"bio"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	AboutME   string `json:"about_me"`
}

var DB *sql.DB

// AUTH HANDLERS
// POST /auth/register
func RegisterHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user.ID = uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)


	_, err = DB.Exec(`
	INSERT INTO users (email, password, first_name, last_name, date_of_birth, bio, avatar, nickname, profile_type) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'public')`,
	user.Email, user.Password, user.FirstName, user.LastName, 
	user.DOB, user.Bio, user.Avatar, user.Nickname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	seesionToken, err := generateSessionToken()
	if err != nil {
		http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
		return
	}
	_, err = DB.Exec("INSERT INTO sessions (user_id, session_token) VALUES (?, ?)", user.ID, seesionToken)
	if err != nil {
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    seesionToken,
		HttpOnly: true,
		Secure:   true,
	})

	json.NewEncoder(w).Encode(map[string]string{"message": "Signup successful", "session_token": seesionToken})

	// NOTE: db processes to be handled through functions
	// NOTE: below query indicates wrong schema (to be checked)
	// TODO: below query to be amended to conform to the regiteration requirement (handle all user properties)
	// _, err = db.DB.Exec("INSERT INTO users (username, password, bio) VALUES (?, ?, ?)", user.Email, user.Password, user.AboutME)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// w.WriteHeader(http.StatusCreated)
}

// POST /auth/login
func LoginHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var user User
	err = DB.QueryRow("SELECT id, password FROM users WHERE email = ?", creds.Email).Scan(&user.ID, &user.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
		return
	}
	_, err = DB.Exec("INSERT INTO sessions (user_id, session_token) VALUES (?, ?)", user.ID, sessionToken)
	if err != nil {
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HttpOnly: true,
		Secure:   true,
	})

	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "session_token": sessionToken})

}

// POST /auth/logout
func LogoutHandler(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "No session token", http.StatusBadRequest)
		return
	}

	seesionToken := cookie.Value

	_, err = DB.Exec("DELETE FROM sessions WHERE session_token = ?", seesionToken)
	if err != nil {
		http.Error(w, "Error deleting session", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})

	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})

}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
