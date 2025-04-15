package sqlite

import (
	"log"
	"time"

	"socialNetwork/pkg/domain/posts"
	"socialNetwork/pkg/errors"
)

// GetPostByID retrieves a post by ID, considering visibility rules
func (s *SQLiteStore) GetPostByID(userID, postID int64) ([]posts.Post, error) {
	// Check if the user can see this post
	canSee, err := s.CanUserSeePost(userID, postID)
	if err != nil {
		log.Print("Error checking post visibility:", err)
		return nil, errors.ErrInternalServer
	}

	if !canSee {
		return nil, errors.ErrNotFound
	}

	// Query posts
	query := `SELECT p.id, p.user_id, p.content, p.image, p.created_at, p.updated_at,
              (SELECT COUNT(*) FROM likes WHERE post_id = p.id) as likes_count
              FROM posts p
              WHERE p.id = ?`

	rows, err := s.db.Query(query, postID)
	if err != nil {
		log.Print("Error querying post:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	result := []posts.Post{}
	for rows.Next() {
		p := posts.Post{}
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Content,
			&p.Image,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Likes,
		)
		if err != nil {
			log.Print("Error scanning post row:", err)
			return nil, errors.ErrInternalServer
		}

		// Get comments for this post
		comments, err := s.getPostComments(p.ID)
		if err != nil {
			log.Print("Error getting comments:", err)
			return nil, errors.ErrInternalServer
		}
		p.Comments = comments

		result = append(result, p)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating post rows:", err)
		return nil, errors.ErrInternalServer
	}

	return result, nil
}

// getPostComments retrieves comments for a post
func (s *SQLiteStore) getPostComments(postID int64) ([]posts.Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at
              FROM comments c
              WHERE c.post_id = ?
              ORDER BY c.created_at DESC`

	rows, err := s.db.Query(query, postID)
	if err != nil {
		log.Print("Error querying comments:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	result := []posts.Comment{}
	for rows.Next() {
		c := posts.Comment{}
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
		)
		if err != nil {
			log.Print("Error scanning comment row:", err)
			return nil, errors.ErrInternalServer
		}
		result = append(result, c)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating comment rows:", err)
		return nil, errors.ErrInternalServer
	}

	return result, nil
}

// CreatePost creates a new post
func (s *SQLiteStore) CreatePost(post posts.Post) (int64, error) {
	// Validate post
	if post.Content == "" && post.Image == "" {
		return 0, errors.ErrInvalidInput
	}

	// Insert post
	query := `INSERT INTO posts (user_id, content, image, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := s.db.Exec(
		query,
		post.UserID,
		post.Content,
		post.Image,
		now,
		now,
	)
	if err != nil {
		log.Print("Error creating post:", err)
		return 0, errors.ErrInternalServer
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Print("Error getting last insert ID:", err)
		return 0, errors.ErrInternalServer
	}

	return id, nil
}

// IsPostOwner checks if a user owns a post
func (s *SQLiteStore) IsPostOwner(userID, postID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM posts WHERE id = ? AND user_id = ?`

	var count int
	err := s.db.QueryRow(query, postID, userID).Scan(&count)
	if err != nil {
		log.Print("Error checking post ownership:", err)
		return false, errors.ErrInternalServer
	}

	return count > 0, nil
}

// EditPost updates a post
func (s *SQLiteStore) EditPost(post posts.Post) error {
	// Validate post
	if post.Content == "" && post.Image == "" {
		return errors.ErrInvalidInput
	}

	// Update post
	query := `UPDATE posts SET content = ?, image = ?, updated_at = ? WHERE id = ?`

	result, err := s.db.Exec(
		query,
		post.Content,
		post.Image,
		time.Now(),
		post.ID,
	)
	if err != nil {
		log.Print("Error updating post:", err)
		return errors.ErrInternalServer
	}

	// Check if any row was affected
	rows, err := result.RowsAffected()
	if err != nil {
		log.Print("Error getting rows affected:", err)
		return errors.ErrInternalServer
	}

	if rows == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// DeletePost deletes a post
func (s *SQLiteStore) DeletePost(postID int64) error {
	// Delete post
	query := `DELETE FROM posts WHERE id = ?`

	result, err := s.db.Exec(query, postID)
	if err != nil {
		log.Print("Error deleting post:", err)
		return errors.ErrInternalServer
	}

	// Check if any row was affected
	rows, err := result.RowsAffected()
	if err != nil {
		log.Print("Error getting rows affected:", err)
		return errors.ErrInternalServer
	}

	if rows == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// CreateComment adds a comment to a post
func (s *SQLiteStore) CreateComment(comment posts.Comment) error {
	// Insert comment
	query := `INSERT INTO comments (post_id, user_id, content, created_at)
              VALUES (?, ?, ?, ?)`

	_, err := s.db.Exec(
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
		time.Now(),
	)
	if err != nil {
		log.Print("Error creating comment:", err)
		return errors.ErrInternalServer
	}

	return nil
}

// CreateLike adds a like to a post
func (s *SQLiteStore) CreateLike(like posts.Like) error {
	// Check if like already exists
	query := `SELECT COUNT(*) FROM likes WHERE post_id = ? AND user_id = ?`

	var count int
	err := s.db.QueryRow(query, like.PostID, like.UserID).Scan(&count)
	if err != nil {
		log.Print("Error checking existing like:", err)
		return errors.ErrInternalServer
	}

	if count > 0 {
		return errors.ErrAlreadyExists
	}

	// Insert like
	insertQuery := `INSERT INTO likes (post_id, user_id, created_at)
                    VALUES (?, ?, ?)`

	_, err = s.db.Exec(
		insertQuery,
		like.PostID,
		like.UserID,
		time.Now(),
	)
	if err != nil {
		log.Print("Error creating like:", err)
		return errors.ErrInternalServer
	}

	return nil
}

// RemoveLikes removes a like from a post
func (s *SQLiteStore) RemoveLikes(like posts.Like) error {
	// Delete like
	query := `DELETE FROM likes WHERE post_id = ? AND user_id = ?`

	result, err := s.db.Exec(query, like.PostID, like.UserID)
	if err != nil {
		log.Print("Error removing like:", err)
		return errors.ErrInternalServer
	}

	// Check if any row was affected
	rows, err := result.RowsAffected()
	if err != nil {
		log.Print("Error getting rows affected:", err)
		return errors.ErrInternalServer
	}

	if rows == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// CanUserSeePost checks if a user can see a post
func (s *SQLiteStore) CanUserSeePost(userID, postID int64) (bool, error) {
	// First, check if the post exists and get the post owner
	var postOwnerID int64
	postQuery := `SELECT user_id FROM posts WHERE id = ?`
	err := s.db.QueryRow(postQuery, postID).Scan(&postOwnerID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil // Post doesn't exist
		}
		log.Print("Error checking post existence:", err)
		return false, errors.ErrInternalServer
	}

	// If user is the post owner, they can see it
	if userID == postOwnerID {
		return true, nil
	}

	// Check if the post owner has a public profile
	var profileType string
	profileQuery := `SELECT profile_type FROM users WHERE id = ?`
	err = s.db.QueryRow(profileQuery, postOwnerID).Scan(&profileType)
	if err != nil {
		log.Print("Error checking profile type:", err)
		return false, errors.ErrInternalServer
	}

	// If profile is public, post is visible
	if profileType == "public" {
		return true, nil
	}

	// For private profiles, check if the user is an accepted follower
	followQuery := `SELECT COUNT(*) FROM followers 
                   WHERE follower_id = ? AND followed_id = ? AND status = 'accepted'`
	var count int
	err = s.db.QueryRow(followQuery, userID, postOwnerID).Scan(&count)
	if err != nil {
		log.Print("Error checking follow status:", err)
		return false, errors.ErrInternalServer
	}

	return count > 0, nil
}