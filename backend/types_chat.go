package main

import "time"

// Chat represents a private message between two users
type Chat struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Content    string    `json:"content"`
	Image      string    `json:"image,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// GroupChat represents a message in a group chat
type GroupChat struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	Image     string    `json:"image,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatMessage represents a message sent through websocket
type ChatMessage struct {
	Type       string    `json:"type"` // "private" or "group"
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id,omitempty"` // Only for private chat
	GroupID    int64     `json:"group_id,omitempty"`    // Only for group chat
	Content    string    `json:"content"`
	Image      string    `json:"image,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
	Client_ID  string    `json:"client_id,omitempty"` // Used for test client to track message origin
}

// WebSocketConnection represents a connected websocket client
type WebSocketConnection struct {
	UserID int64
}

// ChatSession represents active chat sessions for a user
type ChatSession struct {
	UserID    int64  `json:"user_id"`
	OtherUser int64  `json:"other_user,omitempty"` // For private chat
	GroupID   int64  `json:"group_id,omitempty"`   // For group chat
	ChatType  string `json:"chat_type"`            // "private" or "group"
}
