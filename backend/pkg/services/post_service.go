package services

import (
	"backend/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreatePost creates a new post.
func CreatePost(db *sql.DB, userID, content, imageURL, privacy string) (string, error) {
	// Validate privacy setting
	if privacy != "public" && privacy != "private" && privacy != "almost_private" {
		return "", errors.New("invalid privacy setting")
	}

	// Generate a new post ID and insert the new post into the database
	postID := uuid.New().String()
	query := `INSERT INTO posts (id, user_id, content, image_url, privacy, created_at)
			  VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, postID, userID, content, imageURL, privacy, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to create post: %w", err)
	}

	return postID, nil
}

// GetPost retrieves a specific post by its ID.
func GetPost(db *sql.DB, postID string) (*models.Post, error) {
	// Retrieve the post from the database
	var post models.Post
	query := `SELECT id, user_id, content, image_url, privacy, created_at 
			  FROM posts WHERE id = ?`
	err := db.QueryRow(query, postID).Scan(&post.ID, &post.UserID, &post.Content, &post.ImageURL, &post.Privacy, &post.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("post not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve post: %w", err)
	}

	return &post, nil
}

// GetPostWithComments retrieves a post along with its comments.
func GetPostWithComments(db *sql.DB, postID string) (*models.PostWithComments, error) {
	// Retrieve the post
	post, err := GetPost(db, postID)
	if err != nil {
		return nil, err
	}

	// Retrieve the comments for the post
	query := `SELECT id, post_id, user_id, content, created_at 
			  FROM comments WHERE post_id = ? ORDER BY created_at ASC`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse comment: %w", err)
		}
		comments = append(comments, comment)
	}

	// Combine the post and its comments
	postWithComments := &models.PostWithComments{
		Post:     *post,
		Comments: comments,
	}

	return postWithComments, nil
}

// AddComment adds a comment to a post.
func AddComment(db *sql.DB, postID, userID, content string) (string, error) {
	// Generate a new comment ID and insert the comment into the database
	commentID := uuid.New().String()
	query := `INSERT INTO comments (id, post_id, user_id, content, created_at)
			  VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, commentID, postID, userID, content, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to add comment: %w", err)
	}

	return commentID, nil
}

// GetPostsByUser retrieves all posts created by a specific user.
func GetPostsByUser(db *sql.DB, userID string) ([]models.Post, error) {
	// Retrieve posts created by the user
	query := `SELECT id, user_id, content, image_url, privacy, created_at 
			  FROM posts WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.ImageURL, &post.Privacy, &post.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetPublicPosts retrieves all public posts.
func GetPublicPosts(db *sql.DB) ([]models.Post, error) {
	// Retrieve all public posts
	query := `SELECT id, user_id, content, image_url, privacy, created_at 
			  FROM posts WHERE privacy = 'public' ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve public posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.ImageURL, &post.Privacy, &post.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// UpdatePostPrivacy updates the privacy setting of an existing post.
func UpdatePostPrivacy(db *sql.DB, postID, userID, newPrivacy string) error {
	// Validate privacy setting
	if newPrivacy != "public" && newPrivacy != "private" && newPrivacy != "almost_private" {
		return errors.New("invalid privacy setting")
	}

	// Check if the user is the creator of the post
	var creatorID string
	query := `SELECT user_id FROM posts WHERE id = ?`
	err := db.QueryRow(query, postID).Scan(&creatorID)
	if err == sql.ErrNoRows {
		return errors.New("post not found")
	}
	if err != nil {
		return fmt.Errorf("failed to retrieve post creator: %w", err)
	}

	// Only the creator of the post can change its privacy
	if creatorID != userID {
		return errors.New("you can only update your own posts")
	}

	// Update the privacy setting of the post
	query = `UPDATE posts SET privacy = ? WHERE id = ?`
	_, err = db.Exec(query, newPrivacy, postID)
	if err != nil {
		return fmt.Errorf("failed to update post privacy: %w", err)
	}

	return nil
}
