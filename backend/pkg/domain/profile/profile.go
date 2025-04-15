package profile

import "time"

// Profile represents a user profile
type Profile struct {
	UserID       int64     `json:"user_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Nickname     string    `json:"nickname"`
	AboutMe      string    `json:"about_me"`
	Avatar       string    `json:"avatar"`
	DateOfBirth  string    `json:"date_of_birth"`
	ProfileType  string    `json:"profile_type"`
	Followers    int       `json:"followers"`
	Following    int       `json:"following"`
	PostCount    int       `json:"post_count"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
}

// ProfilePrivacyRequest is used when changing profile privacy settings
type ProfilePrivacyRequest struct {
	ProfileType string `json:"profile_type"`
}