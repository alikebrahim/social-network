package sqlite

import (
	"fmt"
	"log"
	"strings"
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
	query := `SELECT p.id, p.user_id, p.content, p.image, p.privacy_level, p.allowed_followers,
              p.created_at, p.updated_at,
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
			&p.PrivacyLevel,
			&p.AllowedFollowers,
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

	// Validate privacy level
	if post.PrivacyLevel == "" {
		post.PrivacyLevel = posts.PrivacyPublic // Default to public
	} else if post.PrivacyLevel != posts.PrivacyPublic && 
		post.PrivacyLevel != posts.PrivacyAlmostPrivate && 
		post.PrivacyLevel != posts.PrivacyPrivate {
		return 0, errors.ErrInvalidInput
	}

	// For private posts, ensure there's a list of allowed followers
	if post.PrivacyLevel == posts.PrivacyPrivate && post.AllowedFollowers == "" {
		// If no allowed followers specified for a private post, default to almost_private
		post.PrivacyLevel = posts.PrivacyAlmostPrivate
	}

	// Insert post
	query := `INSERT INTO posts (user_id, content, image, privacy_level, allowed_followers, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := s.db.Exec(
		query,
		post.UserID,
		post.Content,
		post.Image,
		post.PrivacyLevel,
		post.AllowedFollowers,
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

	// Validate privacy level
	if post.PrivacyLevel != "" && 
	   post.PrivacyLevel != posts.PrivacyPublic && 
	   post.PrivacyLevel != posts.PrivacyAlmostPrivate && 
	   post.PrivacyLevel != posts.PrivacyPrivate {
		return errors.ErrInvalidInput
	}

	// For private posts, ensure there's a list of allowed followers
	if post.PrivacyLevel == posts.PrivacyPrivate && post.AllowedFollowers == "" {
		// If no allowed followers specified for a private post, default to almost_private
		post.PrivacyLevel = posts.PrivacyAlmostPrivate
	}

	// Update post
	query := `UPDATE posts SET 
			  content = ?, 
			  image = ?,
			  privacy_level = COALESCE(?, privacy_level),
			  allowed_followers = CASE 
			    WHEN ? IS NOT NULL THEN ?
				ELSE allowed_followers
			  END,
			  updated_at = ? 
			  WHERE id = ?`

	result, err := s.db.Exec(
		query,
		post.Content,
		post.Image,
		post.PrivacyLevel,
		post.AllowedFollowers, // For the IS NOT NULL check
		post.AllowedFollowers,
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
		return errors.ErrConflict
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

// GetUserPosts retrieves posts for a user profile
func (s *SQLiteStore) GetUserPosts(userID, requestorID int64, limit, offset int) ([]posts.Post, error) {
	// First check if requestor can see the user's posts
	if userID != requestorID {
		// Check if the target user has a public profile
		var profileType string
		profileQuery := `SELECT profile_type FROM users WHERE id = ?`
		err := s.db.QueryRow(profileQuery, userID).Scan(&profileType)
		if err != nil {
			log.Print("Error checking profile type:", err)
			return nil, errors.ErrInternalServer
		}

		// If private profile, check if requestor is an accepted follower
		if profileType == "private" {
			followQuery := `SELECT COUNT(*) FROM followers 
						   WHERE follower_id = ? AND followed_id = ? AND status = 'accepted'`
			var count int
			err = s.db.QueryRow(followQuery, requestorID, userID).Scan(&count)
			if err != nil {
				log.Print("Error checking follow status:", err)
				return nil, errors.ErrInternalServer
			}

			if count == 0 {
				// Not a follower of private profile, return empty result
				return []posts.Post{}, nil
			}
		}
	}

	// Query posts
	query := `SELECT p.id, p.user_id, p.content, p.image, p.privacy_level, p.allowed_followers,
			  p.created_at, p.updated_at,
			  (SELECT COUNT(*) FROM likes WHERE post_id = p.id) as likes_count
			  FROM posts p
			  WHERE p.user_id = ?
			  ORDER BY p.created_at DESC
			  LIMIT ? OFFSET ?`

	rows, err := s.db.Query(query, userID, limit, offset)
	if err != nil {
		log.Print("Error querying user posts:", err)
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
			&p.PrivacyLevel,
			&p.AllowedFollowers,
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

// CanUserSeePost checks if a user can see a post
func (s *SQLiteStore) CanUserSeePost(userID, postID int64) (bool, error) {
	// First, check if the post exists and get post details
	var postOwnerID int64
	var privacyLevel string
	var allowedFollowers string
	
	postQuery := `SELECT user_id, privacy_level, allowed_followers FROM posts WHERE id = ?`
	err := s.db.QueryRow(postQuery, postID).Scan(&postOwnerID, &privacyLevel, &allowedFollowers)
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

	// Check visibility based on privacy level
	switch privacyLevel {
	case posts.PrivacyPublic:
		// Public posts are visible to everyone with profile visibility considerations
		var profileType string
		profileQuery := `SELECT profile_type FROM users WHERE id = ?`
		err = s.db.QueryRow(profileQuery, postOwnerID).Scan(&profileType)
		if err != nil {
			log.Print("Error checking profile type:", err)
			return false, errors.ErrInternalServer
		}

		// If profile is public, public post is visible to everyone
		if profileType == "public" {
			return true, nil
		}

		// For private profiles, public posts still require follower status
		followQuery := `SELECT COUNT(*) FROM followers 
					   WHERE follower_id = ? AND followed_id = ? AND status = 'accepted'`
		var count int
		err = s.db.QueryRow(followQuery, userID, postOwnerID).Scan(&count)
		if err != nil {
			log.Print("Error checking follow status:", err)
			return false, errors.ErrInternalServer
		}
		return count > 0, nil

	case posts.PrivacyAlmostPrivate:
		// Almost private posts are visible only to followers
		followQuery := `SELECT COUNT(*) FROM followers 
					   WHERE follower_id = ? AND followed_id = ? AND status = 'accepted'`
		var count int
		err = s.db.QueryRow(followQuery, userID, postOwnerID).Scan(&count)
		if err != nil {
			log.Print("Error checking follow status:", err)
			return false, errors.ErrInternalServer
		}
		return count > 0, nil

	case posts.PrivacyPrivate:
		// Private posts are visible only to selected followers
		if allowedFollowers == "" {
			return false, nil
		}

		// Check if user is in the allowed followers list
		// This is a simple implementation using string search
		// In a production environment, you might want a many-to-many table for better performance
		for _, idStr := range strings.Split(allowedFollowers, ",") {
			var id int64
			_, err := fmt.Sscanf(idStr, "%d", &id)
			if err == nil && id == userID {
				// User is in the allowed list, verify they are still a follower
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
		}
		return false, nil

	default:
		// Unknown privacy level, default to not visible
		log.Print("Unknown privacy level:", privacyLevel)
		return false, nil
	}
}