package models

import (
	"time"
)

// Notification represents a notification for a user.
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`    // The user who is receiving the notification
	Message   string    `json:"message"`    // The notification message content
	IsRead    bool      `json:"is_read"`    // Flag to mark whether the notification has been read
	CreatedAt time.Time `json:"created_at"` // The timestamp when the notification was created
}

// NotificationType represents the type of notification (e.g., "follow", "group_invite").
type NotificationType struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	NotificationType string    `json:"notification_type"` // Type of notification (e.g., "follow", "group_invite")
	Message          string    `json:"message"`           // Message content specific to the notification type
	IsRead           bool      `json:"is_read"`           // Flag to indicate if the notification has been read
	CreatedAt        time.Time `json:"created_at"`        // Timestamp when the notification was created
}
