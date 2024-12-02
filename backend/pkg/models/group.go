package models

import (
	"time"
)

// Group represents a group that can be created by a user.
type Group struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`   // The ID of the user who created the group
	CreatedAt   time.Time `json:"created_at"` // Timestamp when the group was created
	IsPublic    bool      `json:"is_public"`  // Flag to indicate if the group is public or private
}

// GroupMember represents a user who is part of a group.
type GroupMember struct {
	GroupID  string    `json:"group_id"`
	UserID   string    `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"` // Timestamp when the user joined the group
	IsAdmin  bool      `json:"is_admin"`  // Flag to indicate if the user is an admin
}

// GroupEvent represents an event created within a group.
type GroupEvent struct {
	ID          string    `json:"id"`
	GroupID     string    `json:"group_id"` // The group to which the event belongs
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DayTime     string    `json:"day_time"` // Date and time of the event (YYYY-MM-DD HH:MM:SS)
	CreatedAt   time.Time `json:"created_at"`
}

// GroupInvitation represents an invitation to join a group.
type GroupInvitation struct {
	GroupID   string    `json:"group_id"`
	UserID    string    `json:"user_id"`    // The user invited to the group
	Status    string    `json:"status"`     // Status of the invitation (pending, accepted, declined)
	InvitedAt time.Time `json:"invited_at"` // Timestamp when the invitation was sent
}
