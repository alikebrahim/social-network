package models

import (
	"time"
)

// User represents a user in the system.
type User struct {
	ID          string    `json:"id"`                 // Unique user identifier
	Email       string    `json:"email"`              // User's email address
	Password    string    `json:"-"`                  // User's hashed password (not returned in response)
	FirstName   string    `json:"first_name"`         // User's first name
	LastName    string    `json:"last_name"`          // User's last name
	DateOfBirth string    `json:"date_of_birth"`      // User's date of birth
	Avatar      string    `json:"avatar,omitempty"`   // URL to the user's avatar image
	Nickname    string    `json:"nickname,omitempty"` // User's chosen nickname
	AboutMe     string    `json:"about_me,omitempty"` // A short bio or description
	IsPublic    bool      `json:"is_public"`          // Flag indicating if the profile is public or private
	CreatedAt   time.Time `json:"created_at"`         // Timestamp when the user account was created
}

// UserProfile represents a user's profile data, combining public and private details.
type UserProfile struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth string    `json:"date_of_birth"`
	Avatar      string    `json:"avatar,omitempty"`
	Nickname    string    `json:"nickname,omitempty"`
	AboutMe     string    `json:"about_me,omitempty"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
}
