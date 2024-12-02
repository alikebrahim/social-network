package models

import (
	"time"
)

// PrivateChatMessage represents a message sent in a private chat.
type PrivateChatMessage struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

// GroupChatMessage represents a message sent in a group chat.
type GroupChatMessage struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	SenderID  string    `json:"sender_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// PrivateChat represents a conversation between two users.
type PrivateChat struct {
	ID        string               `json:"id"`
	User1ID   string               `json:"user1_id"`
	User2ID   string               `json:"user2_id"`
	Messages  []PrivateChatMessage `json:"messages"`
	CreatedAt time.Time            `json:"created_at"`
}

// GroupChat represents a group chat conversation.
type GroupChat struct {
	ID           string             `json:"id"`
	GroupID      string             `json:"group_id"`
	Participants []string           `json:"participants"` // List of user IDs participating in the group chat
	Messages     []GroupChatMessage `json:"messages"`
	CreatedAt    time.Time          `json:"created_at"`
}
