package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
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

	// Check if the user can chat with the other user
	canChat, err := h.canUsersChat(userID, otherUserID)
	if err != nil {
		return err
	}
	if !canChat {
		return errors.ErrUnauthorized
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

	// Check if the user can chat with the other user
	canChat, err := h.canUsersChat(userID, otherUserID)
	if err != nil {
		http.Error(w, "Error verifying chat permissions", http.StatusInternalServerError)
		return
	}
	if !canChat {
		http.Error(w, "You can only chat with users who follow you or have public profiles", http.StatusForbidden)
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

			// Check for emoji support
			hasEmoji := containsEmoji(chatMessage.Content)
			if hasEmoji {
				// Ensure the content is UTF-8 encoded
				encodedContent := []byte(chatMessage.Content)
				if !isValidUTF8(encodedContent) {
					log.Printf("Invalid UTF-8 sequence in message")
					continue
				}
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
		http.Error(w, "You must be a member of the group to chat", http.StatusUnauthorized)
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

			// Check for emoji support
			hasEmoji := containsEmoji(chatMessage.Content)
			if hasEmoji {
				// Ensure the content is UTF-8 encoded
				encodedContent := []byte(chatMessage.Content)
				if !isValidUTF8(encodedContent) {
					log.Printf("Invalid UTF-8 sequence in message")
					continue
				}
			}

			// Verify the user is still a member of the group
			isMember, err := h.store.IsGroupMember(groupID, userID)
			if err != nil {
				log.Printf("error checking group membership: %v", err)
				continue
			}
			if !isMember {
				log.Printf("user %d is not a member of group %d", userID, groupID)
				// Send an error message to the client
				errorMsg := map[string]interface{}{
					"type":    "error",
					"message": "You are no longer a member of this group",
				}
				errorJSON, _ := json.Marshal(errorMsg)
				client.Send <- errorJSON
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

// Helper function to check if users can chat with each other
func (h *ChatHandler) canUsersChat(userID1, userID2 int64) (bool, error) {
	// Users can chat if:
	// 1. User1 follows User2 OR User2 follows User1
	// 2. OR User1 has a public profile OR User2 has a public profile

	// Check if User1 follows User2
	isFollowing1, err := h.store.IsFollowing(userID1, userID2)
	if err != nil {
		return false, err
	}

	// Check if User2 follows User1
	isFollowing2, err := h.store.IsFollowing(userID2, userID1)
	if err != nil {
		return false, err
	}

	// If either user follows the other, they can chat
	if isFollowing1 || isFollowing2 {
		return true, nil
	}

	// Check if User1 has a public profile
	var profileType1 string
	profileQuery := `SELECT profile_type FROM users WHERE id = ?`
	err = h.store.DB().QueryRow(profileQuery, userID1).Scan(&profileType1)
	if err != nil {
		log.Print("Error checking profile type:", err)
		return false, errors.ErrInternalServer
	}

	// Check if User2 has a public profile
	var profileType2 string
	err = h.store.DB().QueryRow(profileQuery, userID2).Scan(&profileType2)
	if err != nil {
		log.Print("Error checking profile type:", err)
		return false, errors.ErrInternalServer
	}

	// If either user has a public profile, they can chat
	return profileType1 == "public" || profileType2 == "public", nil
}

// Helper function to check if a string contains emoji characters
func containsEmoji(s string) bool {
	// Simple check for common emoji ranges
	// This is a simplified check, a more robust implementation would use a proper emoji library
	for _, r := range s {
		if r >= 0x1F600 && r <= 0x1F64F || // Emoticons
			r >= 0x1F300 && r <= 0x1F5FF || // Misc Symbols and Pictographs
			r >= 0x1F680 && r <= 0x1F6FF || // Transport and Map
			r >= 0x1F700 && r <= 0x1F77F || // Alchemical Symbols
			r >= 0x1F780 && r <= 0x1F7FF || // Geometric Shapes
			r >= 0x1F800 && r <= 0x1F8FF || // Supplemental Arrows-C
			r >= 0x1F900 && r <= 0x1F9FF || // Supplemental Symbols and Pictographs
			r >= 0x1FA00 && r <= 0x1FA6F || // Chess Symbols
			r >= 0x2600 && r <= 0x26FF || // Miscellaneous Symbols
			r >= 0x2700 && r <= 0x27BF { // Dingbats
			return true
		}
	}
	return false
}

// Helper function to validate UTF-8 encoding
func isValidUTF8(b []byte) bool {
	return strings.ToValidUTF8(string(b), "") == string(b)
}