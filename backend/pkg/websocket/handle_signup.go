package websocket

import (
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func handleSignup(msg Message) {
	response := Message{Type: "register_response"}

	hashedpassword, errHashed := hashPassword(msg.Data.Password)
	if errHashed != nil {
		log.Printf("Failed to hash password: %v", errHashed)
		response.Content = "Signup failed"
		sendResponseToClients(response)
		return
	}

	err := registerUser(msg, hashedpassword)
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

func registerUser(msg Message, hashedPassword string) error {
	_, err := DB.Exec(
		"INSERT INTO users (email, password, first_name, last_name, date_of_birth, bio, avatar, nickname) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		msg.Data.Email, hashedPassword, msg.Data.FirstName, msg.Data.LastName, msg.Data.BirthDate, msg.Data.Bio, msg.Data.Avatar, msg.Data.Nickname,
	)
	return err
}

func hashPassword(password string) (string, error) {
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedpassword), nil
}
