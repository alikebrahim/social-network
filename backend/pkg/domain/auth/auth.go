package auth

import "errors"

// Common auth domain errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrSessionExpired     = errors.New("session expired")
	ErrEmailTaken         = errors.New("email is already taken")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPassword    = errors.New("invalid password format")
)

// UserAccount represents a user account in the system
type UserAccount struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	First_name    string `json:"first_name"`
	Last_name     string `json:"last_name"`
	Date_of_birth string `json:"date_of_birth"`
	Avatar        string `json:"avatar"`
	Nickname      string `json:"nickname"`
	About_me      string `json:"about_me"`
	Profile_type  string `json:"profile_type"`
}

// LoginRequest is used when a user logs in
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegistrationRequest is used when a user registers
type RegistrationRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	First_name    string `json:"first_name"`
	Last_name     string `json:"last_name"`
	Date_of_birth string `json:"date_of_birth,omitempty"`
	Profile_type  string `json:"profile_type"`
}

// Session represents a user session
type Session struct {
	ID             string        `json:"id"`
	UserID         int64         `json:"user_id"`
	CreatedAt      string        `json:"created_at"`
	ExpiresAt      string        `json:"expires_at"`
	FollowCount    int64         `json:"followCount"`
	FollowingCount int64         `json:"followingCount"`
	Followers      []UserAccount `json:"followers"`
	Following      []UserAccount `json:"following"`
}