package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// GET /chats
func (s *APIServer) HandleGetChats(w http.ResponseWriter, r *http.Request) error {
	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Get user's chat history
	chats, err := s.store.GetUserChats(userID)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, chats)
}

// GET /chats/{userId}
func (s *APIServer) HandleGetChatHistory(w http.ResponseWriter, r *http.Request) error {
	// Parse user ID from URL
	otherUserID, err := strconv.ParseInt(r.PathValue("userId"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Get chat history between the two users
	chatHistory, err := s.store.GetChatHistory(userID, otherUserID)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, chatHistory)
}

// GET /groups/{id}/chat
func (s *APIServer) HandleGetGroupChatHistory(w http.ResponseWriter, r *http.Request) error {
	// Parse group ID from URL
	groupID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Check if user is a member of the group
	isMember, err := s.store.IsGroupMember(groupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrUnauthorized
	}

	// Get group chat history
	groupChatHistory, err := s.store.GetGroupChatHistory(groupID)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, groupChatHistory)
}

// WebSocket /ws/chat/{userId}
func (s *APIServer) HandleChatWebSocket(w http.ResponseWriter, r *http.Request) {
	// Parse target user ID from URL
	receiverID, err := strconv.ParseInt(r.PathValue("userId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get user ID from session
	senderID, err := s.getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get client ID from query parameters (for testing purposes)
	clientID := r.URL.Query().Get("client_id")

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Create a WebSocket client
	client := &WebSocketClient{
		UserID:   senderID,
		Conn:     conn,
		Send:     make(chan []byte, 256), // Buffer up to 256 messages
		ClientID: clientID,               // Store the client ID for test client
	}

	// Register the client with the hub
	hub.register <- client

	// Record the user being connected for private chat
	log.Printf("Chat WebSocket established between users %d and %d (client: %s)", senderID, receiverID, clientID)

	// Set up a handler for reading private chat messages
	go func() {
		defer func() {
			hub.unregister <- client
			conn.Close()
		}()

		for {
			// Read message from the WebSocket
			var message ChatMessage
			err := conn.ReadJSON(&message)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				break
			}

			// Ensure message has correct type
			if message.Type != "private" {
				log.Printf("Invalid message received: incorrect type")
				continue
			}

			// Set message metadata
			message.SenderID = senderID
			message.ReceiverID = receiverID
			message.Timestamp = time.Now()

			// Save the message to the database
			chat := &Chat{
				SenderID:   senderID,
				ReceiverID: receiverID,
				Content:    message.Content,
				Image:      message.Image,
				CreatedAt:  message.Timestamp,
			}

			_, err = s.store.SaveChat(chat)
			if err != nil {
				log.Printf("Error saving private chat message: %v", err)
				continue
			}

			// Broadcast the message
			wsMessage := WebSocketMessage{
				Type:      "chat",
				Payload:   message,
				Client_ID: message.Client_ID, // Pass through the client ID for test client
			}

			messageJSON, err := json.Marshal(wsMessage)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			hub.broadcast <- messageJSON
		}
	}()

	// Start the client's write pump only since we're handling reads separately
	client.StartWritePumpOnly()
}

// WebSocket /ws/groups/{id}/chat
func (s *APIServer) HandleGroupChatWebSocket(w http.ResponseWriter, r *http.Request) {
	// Parse group ID from URL
	groupID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get client ID from query parameters (for testing purposes)
	clientID := r.URL.Query().Get("client_id")

	// Check if user is a member of the group
	isMember, err := s.store.IsGroupMember(groupID, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Create a WebSocket client
	client := &WebSocketClient{
		UserID:   userID,
		Conn:     conn,
		Send:     make(chan []byte, 256), // Buffer up to 256 messages
		ClientID: clientID,               // Store the client ID for test client
	}

	// Register the client with the hub
	hub.register <- client

	// Also register with the group
	hub.RegisterGroupClient(groupID, client)

	log.Printf("Group Chat WebSocket established for user %d in group %d (client: %s)", userID, groupID, clientID)

	// Set up a handler for incoming messages
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("WebSocket closed: %d %s", code, text)

		// Unregister the client
		hub.unregister <- client

		// Call the default close handler
		return nil
	})

	// Process incoming messages
	go func() {
		defer func() {
			hub.unregister <- client
			conn.Close()
		}()

		for {
			// Read message from the WebSocket
			var message ChatMessage
			err := conn.ReadJSON(&message)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				break
			}

			// Ensure message has correct type and group ID
			if message.Type != "group" || message.GroupID != groupID {
				log.Printf("Invalid message received: incorrect type or group ID")
				continue
			}

			// Set message metadata
			message.SenderID = userID
			message.Timestamp = time.Now()

			// Save the message to the database
			groupChat := &GroupChat{
				GroupID:   groupID,
				SenderID:  userID,
				Content:   message.Content,
				Image:     message.Image,
				CreatedAt: message.Timestamp,
			}

			_, err = s.store.SaveGroupChat(groupChat)
			if err != nil {
				log.Printf("Error saving group chat message: %v", err)
				continue
			}

			// Broadcast the message
			wsMessage := WebSocketMessage{
				Type:      "chat",
				Payload:   message,
				Client_ID: message.Client_ID, // Pass through the client ID for test client
			}

			messageJSON, err := json.Marshal(wsMessage)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			hub.broadcast <- messageJSON
		}
	}()

	// Start the client's write pump only since we're handling reads separately
	client.StartWritePumpOnly()
}

// Helper to handle an incoming private chat message
func (s *APIServer) handlePrivateChatMessage(senderID, receiverID int64, content, image string) error {
	// Create chat message
	chat := &Chat{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Image:      image,
		CreatedAt:  time.Now(),
	}

	// Save to database
	_, err := s.store.SaveChat(chat)
	if err != nil {
		return err
	}

	// Create WebSocket message
	message := ChatMessage{
		Type:       "private",
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Image:      image,
		Timestamp:  chat.CreatedAt,
	}

	// Broadcast
	wsMessage := WebSocketMessage{
		Type:    "chat",
		Payload: message,
	}

	messageJSON, err := json.Marshal(wsMessage)
	if err != nil {
		return err
	}

	hub.broadcast <- messageJSON
	return nil
}
