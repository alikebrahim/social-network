package groups

import "time"

// Group represents a social group
type Group struct {
	ID          int64     `json:"id"`
	CreatorID   int64     `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GroupMember represents a user's membership in a group
type GroupMember struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	UserID    int64     `json:"user_id"`
	Status    string    `json:"status"` // pending, accepted, rejected
	CreatedAt time.Time `json:"created_at"`
}

// GroupEvent represents an event in a group
type GroupEvent struct {
	ID          int64     `json:"id"`
	GroupID     int64     `json:"group_id"`
	CreatorID   int64     `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	CreatedAt   time.Time `json:"created_at"`
}

// EventResponse represents a user's response to a group event
type EventResponse struct {
	ID        int64     `json:"id"`
	EventID   int64     `json:"event_id"`
	UserID    int64     `json:"user_id"`
	Response  string    `json:"response"` // going, not_going
	CreatedAt time.Time `json:"created_at"`
}

// GroupInvite represents an invitation to join a group
type GroupInvite struct {
	GroupID   int64 `json:"group_id"`
	InviterID int64 `json:"inviter_id"`
	InviteeID int64 `json:"invitee_id"`
}

// GroupRequest represents a request to join a group
type GroupRequest struct {
	GroupID     int64 `json:"group_id"`
	RequesterID int64 `json:"requester_id"`
}

// CreateGroupRequest is used when creating a new group
type CreateGroupRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// CreateEventRequest is used when creating a new event
type CreateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
}

// EventResponseRequest is used when responding to an event
type EventResponseRequest struct {
	Response string `json:"response"` // going, not_going
}