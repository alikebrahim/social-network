package routes

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64  `json:"id"`
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	result, err := DB.Exec(`
	INSERT INTO users (email, password, first_name, last_name, date_of_birth, bio, avatar, nickname, profile_type) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'public')`,
		user.Email, user.Password, user.FirstName, user.LastName,
		user.DOB, user.Bio, user.Avatar, user.Nickname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusInternalServerError)
		return
	}

	seesionToken, err := generateSessionToken()
	if err != nil {
		http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
		return
	}
	println(user.ID)
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
	println(user.ID)
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

func getUserIDFromSession(sessionToken string) (string, error) {
	log.Println("Getting session token")

	// ✅ Remove "session_token=" prefix if present
	token := strings.TrimPrefix(sessionToken, "session_token=")
	token = strings.TrimSpace(token)

	log.Println("Clean Session Token:", token) // Debugging log

	var userID string // ✅ Change userID from int to string

	err := DB.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", token).Scan(&userID)
	if err != nil {
		log.Println("Session token not found in database:", err)
		return "", err
	}

	log.Println("User ID Retrieved:", userID) // Debugging log
	return userID, nil
}

func GetSessionToken(r *http.Request) (string, error) {
	log.Println("GetSessionToken Called")

	// ✅ Try retrieving the session token from cookies
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "session_token" {
			token := strings.TrimSpace(cookie.Value)
			log.Println("Session Token Found in Cookies:", token)
			return token, nil
		}
	}

	// ✅ Try retrieving from `Session_token` header
	authHeader := r.Header.Get("Session_token")
	if authHeader != "" {
		token := strings.TrimSpace(strings.Split(authHeader, ";")[0]) // ✅ Remove extra attributes
		log.Println("Session Token Found in Header:", token)
		return token, nil
	}

	// ✅ Try `Authorization: Bearer token`
	authHeader = r.Header.Get("Authorization")
	if authHeader != "" {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		log.Println("Session Token Found in Authorization Header:", token)
		return token, nil
	}

	log.Println("No session token found")
	return "", errors.New("No session token found")
}
