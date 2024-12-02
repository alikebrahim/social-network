package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader to upgrade HTTP connections to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow requests from any origin (adjust this for production security)
		return true
	},
}

// PrivateChatHandler handles WebSocket connections for private chats
func PrivateChatHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		log.Println("Private chat WebSocket connection established")

		// Handle messages
		for {
			var msg struct {
				SenderID   string `json:"sender_id"`
				ReceiverID string `json:"receiver_id"`
				Message    string `json:"message"`
			}

			// Read message from client
			err := conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("WebSocket connection closed")
					return
				}
				log.Printf("Error reading message: %v", err)
				continue
			}

			// Save message to the database
			query := `INSERT INTO private_messages (sender_id, receiver_id, message, created_at) 
					  VALUES (?, ?, ?, datetime('now'))`
			_, dbErr := db.Exec(query, msg.SenderID, msg.ReceiverID, msg.Message)
			if dbErr != nil {
				log.Printf("Error saving message to database: %v", dbErr)
				continue
			}

			// Echo message back to the sender
			err = conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Error sending message: %v", err)
				continue
			}
		}
	}
}

// GroupChatHandler handles WebSocket connections for group chats
func GroupChatHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		log.Println("Group chat WebSocket connection established")

		// Handle messages
		for {
			var msg struct {
				GroupID  string `json:"group_id"`
				SenderID string `json:"sender_id"`
				Message  string `json:"message"`
			}

			// Read message from client
			err := conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("WebSocket connection closed")
					return
				}
				log.Printf("Error reading message: %v", err)
				continue
			}

			// Save message to the database
			query := `INSERT INTO group_messages (group_id, sender_id, message, created_at) 
					  VALUES (?, ?, ?, datetime('now'))`
			_, dbErr := db.Exec(query, msg.GroupID, msg.SenderID, msg.Message)
			if dbErr != nil {
				log.Printf("Error saving message to database: %v", dbErr)
				continue
			}

			// Broadcast message to the group (implementation can vary)
			// For simplicity, we echo it back to the sender
			err = conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Error sending message: %v", err)
				continue
			}
		}
	}
}
