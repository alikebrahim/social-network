package websocket

import (
	"database/sql"
	"log"
)

func handleLogout(msg Message) {
	response := Message{Type: "logout_response"}

	sessionToken := msg.SessionToken
	if sessionToken == "" {
		log.Println("Logout failed: session token missing")
		response.Content = "Logout failed: No session token provided"
		sendResponseToClients(response)
		return
	}

	userEmail, err := logout(sessionToken)
	if err != nil {
		log.Printf("Logout failed for session token %s: %v", sessionToken, err)
		response.Content = "Logout failed"
	} else {
		log.Printf("User logged out: %s", userEmail)
		response.Content = "Logout successful"

	
		delete(loggedInUsers, userEmail)
	}

	sendResponseToClients(response)
}

func logout(sessionToken string) (string, error) {
	var userEmail string

	err := DB.QueryRow("SELECT users.email FROM users INNER JOIN sessions ON users.id = sessions.user_id WHERE sessions.session_token = ?", sessionToken).Scan(&userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil 
		}
		return "", err
	}


	_, err = DB.Exec("DELETE FROM sessions WHERE session_token = ?", sessionToken)
	if err != nil {
		return "", err
	}

	return userEmail, nil
}
