package models

import (
	"time"
)

// Post represents a user's post in the system.
type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`             // The user who created the post
	Content   string    `json:"content"`             // The content of the post (text)
	ImageURL  string    `json:"image_url,omitempty"` // URL of an image associated with the post (optional)
	Privacy   string    `json:"privacy"`             // Privacy setting (e.g., "public", "private", "almost_private")
	CreatedAt time.Time `json:"created_at"`          // The timestamp when the post was created
}

// Comment represents a comment on a post.
type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`    // The ID of the post being commented on
	UserID    string    `json:"user_id"`    // The user who wrote the comment
	Content   string    `json:"content"`    // The content of the comment
	CreatedAt time.Time `json:"created_at"` // Timestamp when the comment was created
}

// PostWithComments represents a post along with its associated comments.
type PostWithComments struct {
	Post     Post      `json:"post"`     // The post itself
	Comments []Comment `json:"comments"` // List of comments on the post
}
