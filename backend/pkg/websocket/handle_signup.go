package websocket

import (
	"log"
	"strings"
)

func handleSignup(msg Message) {
	response := Message{Type: "register_response"}

	err := registerUser(msg)
	if err != nil {
		log.Printf("Failed to sign up user: %v", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			response.Content = "Email already exists"
		} else {
			response.Content = "Signup failed"
		}
	} else {
		log.Printf("User signed up: %s", msg.Data.Email)
		response.Content = "Signup successful"
	}

	sendResponseToClients(response)
}

func registerUser(msg Message) error {
	_, err := DB.Exec("INSERT INTO users (email, password, first_name, last_name, birth_date) VALUES (?, ?, ?, ?, ?)",
		msg.Data.Email, msg.Data.Password, msg.Data.FirstName, msg.Data.LastName, msg.Data.BirthDate)
	return err
}