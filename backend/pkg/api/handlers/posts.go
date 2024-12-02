package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CreatePostRequest struct {
	Content  string `json:"content"`
	ImageURL string `json:"image_url,omitempty"`
	Privacy  string `json:"privacy"` // "public", "private", "almost_private"
}

type CommentRequest struct {
	PostID  string `json:"post_id"`
	Content string `json:"content"`
}

type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url,omitempty"`
	Privacy   string    `json:"privacy"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatePostHandler handles creating a new post
func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePostRequest

		// Parse JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate privacy settings
		if req.Privacy != "public" && req.Privacy != "private" && req.Privacy != "almost_private" {
			http.Error(w, "Invalid privacy setting", http.StatusBadRequest)
			return
		}

		// Insert the post into the database
		postID := uuid.New().String()
		userID := r.Context().Value("user_id").(string) // Assuming user ID is stored in context
		query := `INSERT INTO posts (id, user_id, content, image_url, privacy, created_at) VALUES (?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(query, postID, userID, req.Content, req.ImageURL, req.Privacy, time.Now())
		if err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully", "post_id": postID})
	}
}

// GetPostHandler retrieves a specific post by ID
func GetPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postID := r.URL.Query().Get("id")
		if postID == "" {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		// Retrieve the post from the database
		var post Post
		query := `SELECT id, user_id, content, image_url, privacy, created_at FROM posts WHERE id = ?`
		err := db.QueryRow(query, postID).Scan(&post.ID, &post.UserID, &post.Content, &post.ImageURL, &post.Privacy, &post.CreatedAt)
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
}

// AddCommentHandler handles adding a comment to a post
func AddCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CommentRequest

		// Parse JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Insert the comment into the database
		commentID := uuid.New().String()
		userID := r.Context().Value("user_id").(string) // Assuming user ID is stored in context
		query := `INSERT INTO comments (id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)`
		_, err := db.Exec(query, commentID, req.PostID, userID, req.Content, time.Now())
		if err != nil {
			http.Error(w, "Failed to add comment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Comment added successfully", "comment_id": commentID})
	}
}
