package websocket

import (
	"log"
	"time"
)

func handlePostCreation(msg Message) {
	response := Message{Type: "post_response"}

	
	userID, err := getUserIDFromSession(msg.SessionToken)
	if err != nil {
		log.Printf("Post creation failed: %v", err)
		response.Content = "Authentication failed"
		sendResponseToClients(response)
		return
	}


	if msg.Post.Content == "" && msg.Post.Image == "" {
		response.Content = "Post cannot be empty"
		sendResponseToClients(response)
		return
	}

	privacy := msg.Post.Privacy
	if privacy == "" {
		privacy = "public" // Default value
	} else if privacy != "public" && privacy != "private"  {
		response.Content = "Invalid privacy setting"
		sendResponseToClients(response)
		return
	}

	_, err = DB.Exec("INSERT INTO posts (user_id, content, image, privacy, created_at) VALUES (?, ?, ?, ?, ?)",
		userID, msg.Post.Content, msg.Post.Image, privacy, time.Now().UTC())
	if err != nil {
		log.Printf("Post creation failed: %v", err)
		response.Content = "Post creation failed"
		sendResponseToClients(response)
		return
	}

	response.Content = "Post created successfully"
	sendResponseToClients(response)
}
