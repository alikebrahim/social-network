package posts

import (
	"fmt"
	"strings"
	"time"
)

// Privacy levels for posts
const (
	PrivacyPublic      = "public"      // Visible to all users
	PrivacyAlmostPrivate = "almost_private" // Visible only to followers
	PrivacyPrivate     = "private"     // Visible only to selected followers
)

// Post represents a user post
type Post struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Content         string    `json:"content"`
	Image           string    `json:"image"`
	PrivacyLevel    string    `json:"privacy_level"`
	AllowedFollowers string    `json:"allowed_followers,omitempty"` // Comma-separated list of user IDs for private posts
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Likes           int       `json:"likes"`
	Comments        []Comment `json:"comments"`
}

// GetAllowedFollowerIDs returns the slice of user IDs that are allowed to view a private post
func (p *Post) GetAllowedFollowerIDs() []int64 {
	if p.AllowedFollowers == "" {
		return []int64{}
	}
	
	// Convert comma-separated string to slice of int64
	followerIDs := []int64{}
	for _, idStr := range strings.Split(p.AllowedFollowers, ",") {
		var id int64
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err == nil {
			followerIDs = append(followerIDs, id)
		}
	}
	return followerIDs
}

// Comment represents a comment on a post
type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Like represents a like on a post
type Like struct {
	PostID int64 `json:"post_id"`
	UserID int64 `json:"user_id"`
}

// GroupPost represents a post in a group
type GroupPost struct {
	Post
	GroupID int64 `json:"group_id"`
}