package websocket

import (
	"fmt"
	"log"
	"strings"
)

func handleLogin(msg Message) {
	response := Message{Type: "register_response"}
	fmt.Println(msg)
	login(msg)
	err := login(msg)
	if err != nil {
		log.Printf("Failed to sign up user: %v", err)
		if strings.Contains(err.Error(), "no rows in result set") {
			response.Content = "Email or password is incorrect"
		} else {
			response.Content = "Signup failed"
		}
	} else {
		log.Printf("User logged in: %s", msg.Data.Email)
		response.Content = "Login successful"
	}

	sendResponseToClients(response)
}

func login(msg Message) error {
	var userID int
	err := DB.QueryRow("SELECT id FROM users WHERE email = ? AND password = ?",
		msg.Data.Email, msg.Data.Password).Scan(&userID)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	fmt.Println("User does not exist")
		// 	return nil
		// }
		return err
	}
	fmt.Println("User exists with ID:", userID)
	logedInUsers[sender] = true
	return nil
}
