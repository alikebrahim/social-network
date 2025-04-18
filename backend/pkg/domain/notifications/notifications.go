package notifications

import "time"

// Notification types
const (
	TypeFollowRequest = "follow_request"
	TypeGroupInvite   = "group_invite"
	TypeJoinRequest   = "join_request"
	TypeEventCreated  = "event_created"
	TypeComment       = "comment"
	TypeLike          = "like"
)

// Notification represents a user notification
type Notification struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Type        string    `json:"type"`
	Content     string    `json:"content"`
	IsRead      bool      `json:"is_read"`
	RelatedID   int64     `json:"related_id"`
	CreatedAt   time.Time `json:"created_at"`
	SenderID    int64     `json:"sender_id,omitempty"`
	SenderName  string    `json:"sender_name,omitempty"`
	SenderImage string    `json:"sender_image,omitempty"`
}

// NotificationsList represents a paginated list of notifications
type NotificationsList struct {
	Notifications []Notification `json:"notifications"`
	Count         int            `json:"count"`
}