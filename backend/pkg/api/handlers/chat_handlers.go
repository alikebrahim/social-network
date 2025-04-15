package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	gorillaws "github.com/gorilla/websocket"
	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/domain/chat"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/storage"
	"socialNetwork/pkg/websocket"
)

// ChatHandler manages chat-related handlers
type ChatHandler struct {
	store storage.Storage
	hub   *websocket.Hub
}

// NewChatHandler creates a new chat handler
func NewChatHandler(store storage.Storage, hub *websocket.Hub) *ChatHandler {
	return &ChatHandler{
		store: store,
		hub:   hub,
	}
}

// HandleGetChats gets all chats for a user
func (h *ChatHandler) HandleGetChats(w http.ResponseWriter, r *http.Request) error {
	// Get the user ID from the session
	userID, err := getUserIDFromRequest(r, h.store)
	if err != nil {
		return err
	}

	// Get the user's chats
	chats, err := h.store.GetUserChats(userID)
	if err != nil {
		return err
	}

	return utils.WriteJson(w, http.StatusOK, chats)
}

// HandleGetChatHistory gets the chat history between two users
func (h *ChatHandler) HandleGetChatHistory(w http.ResponseWriter, r *http.Request) error {
	// Get the user ID from the session
	userID, err := getUserIDFromRequest(r, h.store)
	if err != nil {
		return err
	}

	// Parse the other user's ID from the URL
	otherUserIDStr := r.PathValue("userId")
	otherUserID, err := strconv.ParseInt(otherUserIDStr, 10, 64)
	if err != nil {
		return errors.ErrBadRequest
	}

	// Get the chat history
	chatHistory, err := h.store.GetChatHistory(userID, otherUserID)
	if err != nil {
		return err
	}

	return utils.WriteJson(w, http.StatusOK, chatHistory)
}

// HandleGetGroupChatHistory gets the chat history for a group
func (h *ChatHandler) HandleGetGroupChatHistory(w http.ResponseWriter, r *http.Request) error {
	// Get the user ID from the session
	userID, err := getUserIDFromRequest(r, h.store)
	if err != nil {
		return err
	}

	// Parse the group ID from the URL
	groupIDStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		return errors.ErrBadRequest
	}

	// Check if the user is a member of the group
	isMember, err := h.store.IsGroupMember(groupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.ErrUnauthorized
	}

	// Get the group chat history
	chatHistory, err := h.store.GetGroupChatHistory(groupID)
	if err != nil {
		return err
	}

	return utils.WriteJson(w, http.StatusOK, chatHistory)
}

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now
		// In production, should be restricted
		return true
	},
}

// HandleChatWebSocket handles WebSocket connections for direct user chats
func (h *ChatHandler) HandleChatWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the session
	userID, err := getUserIDFromRequest(r, h.store)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the other user's ID from the URL
	otherUserIDStr := r.PathValue("userId")
	otherUserID, err := strconv.ParseInt(otherUserIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	// Create a new client
	client := &websocket.Client{
		Hub:      h.hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		TargetID: otherUserID,
		IsGroup:  false,
	}

	// Register the client with the hub
	h.hub.Register <- client

	// Start the client's read and write pumps
	go client.ReadPump()
	go client.WritePump()

	// Save messages received from the WebSocket to the database
	// This is handled in the ReadPump method
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if gorillaws.IsUnexpectedCloseError(err, gorillaws.CloseGoingAway, gorillaws.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}

			// Parse the message
			var chatMessage chat.ChatMessage
			if err := json.Unmarshal(message, &chatMessage); err != nil {
				log.Printf("error unmarshaling message: %v", err)
				continue
			}

			// Create a chat object to save
			chatObj := &chat.Chat{
				SenderID:   userID,
				ReceiverID: otherUserID,
				Content:    chatMessage.Content,
				CreatedAt:  time.Now(),
			}

			// Save the message to the database
			_, err = h.store.SaveChat(chatObj)
			if err != nil {
				log.Printf("error saving chat message: %v", err)
				continue
			}

			// The message will be sent to both users in the client.ReadPump method
		}
	}()
}

// HandleGroupChatWebSocket handles WebSocket connections for group chats
func (h *ChatHandler) HandleGroupChatWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the session
	userID, err := getUserIDFromRequest(r, h.store)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the group ID from the URL
	groupIDStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Check if the user is a member of the group
	isMember, err := h.store.IsGroupMember(groupID, userID)
	if err != nil {
		http.Error(w, "Error verifying membership", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	// Create a new client
	client := &websocket.Client{
		Hub:      h.hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		TargetID: groupID,
		IsGroup:  true,
	}

	// Register the client with the hub
	h.hub.Register <- client

	// Start the client's read and write pumps
	go client.ReadPump()
	go client.WritePump()

	// Save messages received from the WebSocket to the database
	// This is handled in the ReadPump method
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if gorillaws.IsUnexpectedCloseError(err, gorillaws.CloseGoingAway, gorillaws.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}

			// Parse the message
			var chatMessage chat.ChatMessage
			if err := json.Unmarshal(message, &chatMessage); err != nil {
				log.Printf("error unmarshaling message: %v", err)
				continue
			}

			// Create a group chat object to save
			chatObj := &chat.GroupChat{
				GroupID:   groupID,
				SenderID:  userID,
				Content:   chatMessage.Content,
				CreatedAt: time.Now(),
			}

			// Save the message to the database
			_, err = h.store.SaveGroupChat(chatObj)
			if err != nil {
				log.Printf("error saving group chat message: %v", err)
				continue
			}

			// The message will be sent to all group members in the client.ReadPump method
		}
	}()
}

// Helper function to get user ID from request
func getUserIDFromRequest(r *http.Request, store storage.Storage) (int64, error) {
	// Get the session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return 0, errors.ErrUnauthorized
		}
		return 0, errors.ErrInternalServer
	}

	// Get the user ID from the session
	sessionToken := cookie.Value
	userID, err := store.GetUserIdBySession(sessionToken)
	if err != nil {
		return 0, err
	}

	return userID, nil
}