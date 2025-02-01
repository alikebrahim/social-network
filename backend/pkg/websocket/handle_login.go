package websocket

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var loggedInUsers = make(map[string]bool)

func handleLogin(msg Message) {
	response := Message{Type: "login_response"}

	sessionToken,err := login(msg)
	if err != nil {
		log.Printf("Login failed for user %s: %v", msg.Data.Email, err)
		if err == sql.ErrNoRows {
			response.Content = "Email or password is incorrect"
		} else {
			response.Content = "Login failed"
		}
	} else {
		log.Printf("User logged in: %s", msg.Data.Email)
		response.Content = "Login successful"
		response.SessionToken = sessionToken // ✅ Send the token
		log.Print("Session token sent to client: ", response.SessionToken)
	}

	sendResponseToClients(response)
}

func login(msg Message) (string, error) {
	var userID int
	var hashedPassword string
	
	err := DB.QueryRow("SELECT id, password FROM users WHERE email = ?", msg.Data.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(msg.Data.Password)); err != nil {
		return "", err
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		log.Println("Error generating session token:", err)
		return "", err
	}

	_, err = DB.Exec("INSERT INTO sessions (user_id, session_token) VALUES (?, ?)", userID, sessionToken)
	if err != nil {
		log.Println("Error saving session:", err)
		return "", err
	}

	loggedInUsers[msg.Data.Email] = true
	fmt.Println("User logged in:", msg.Data.Email)

	return sessionToken, nil
}

