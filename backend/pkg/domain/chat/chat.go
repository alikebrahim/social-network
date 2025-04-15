package chat

import "time"

// Chat represents a chat message between users
type Chat struct {
	ID        int64     `json:"id"`
	SenderID  int64     `json:"sender_id"`
	ReceiverID int64    `json:"receiver_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// GroupChat represents a chat message in a group
type GroupChat struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatMessage represents a message sent over WebSocket
type ChatMessage struct {
	Type      string    `json:"type"`
	SenderID  int64     `json:"sender_id"`
	ReceiverID int64    `json:"receiver_id,omitempty"`
	GroupID   int64     `json:"group_id,omitempty"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}