package services

import (
	"backend/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// FollowUser sends a follow request from one user to another.
// If the target user has a public profile, the follow request is automatically accepted.
func FollowUser(db *sql.DB, followerID, targetUserID string) (string, error) {
	// Check if the target user exists
	var isPublic bool
	query := `SELECT is_public FROM users WHERE id = ?`
	err := db.QueryRow(query, targetUserID).Scan(&isPublic)
	if err == sql.ErrNoRows {
		return "", errors.New("target user not found")
	}
	if err != nil {
		return "", fmt.Errorf("failed to retrieve target user: %w", err)
	}

	// If the target user's profile is public, automatically follow the user
	if isPublic {
		// Automatically accept the follow request for public profiles
		query = `INSERT INTO followers (follower_id, following_id, status, created_at)
				  VALUES (?, ?, 'accepted', ?)`
		_, err := db.Exec(query, followerID, targetUserID, time.Now())
		if err != nil {
			return "", fmt.Errorf("failed to follow user: %w", err)
		}
		return "followed", nil
	}

	// Otherwise, send a follow request (pending status)
	query = `INSERT INTO followers (follower_id, following_id, status, created_at)
			  VALUES (?, ?, 'pending', ?)`
	_, err = db.Exec(query, followerID, targetUserID, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to send follow request: %w", err)
	}

	return "request_sent", nil
}

// AcceptFollowRequest allows the user to accept a follow request.
func AcceptFollowRequest(db *sql.DB, followerID, targetUserID string) error {
	// Check if there is a pending follow request
	query := `UPDATE followers SET status = 'accepted' 
			  WHERE follower_id = ? AND following_id = ? AND status = 'pending'`
	result, err := db.Exec(query, followerID, targetUserID)
	if err != nil {
		return fmt.Errorf("failed to accept follow request: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no pending follow request found")
	}

	return nil
}

// UnfollowUser allows a user to unfollow another user.
func UnfollowUser(db *sql.DB, followerID, targetUserID string) error {
	// Delete the follow relationship from the followers table
	query := `DELETE FROM followers WHERE follower_id = ? AND following_id = ?`
	result, err := db.Exec(query, followerID, targetUserID)
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("follow relationship not found")
	}

	return nil
}

// GetFollowers retrieves the list of users who follow the given user.
func GetFollowers(db *sql.DB, userID string) ([]models.User, error) {
	// Retrieve followers of the user
	query := `SELECT u.id, u.email, u.first_name, u.last_name, u.avatar
			  FROM users u 
			  JOIN followers f ON u.id = f.follower_id
			  WHERE f.following_id = ? AND f.status = 'accepted'`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve followers: %w", err)
	}
	defer rows.Close()

	var followers []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Avatar)
		if err != nil {
			return nil, fmt.Errorf("failed to parse follower: %w", err)
		}
		followers = append(followers, user)
	}

	return followers, nil
}

// GetFollowing retrieves the list of users that the given user is following.
func GetFollowing(db *sql.DB, userID string) ([]models.User, error) {
	// Retrieve users that the given user is following
	query := `SELECT u.id, u.email, u.first_name, u.last_name, u.avatar
			  FROM users u
			  JOIN followers f ON u.id = f.following_id
			  WHERE f.follower_id = ? AND f.status = 'accepted'`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve following users: %w", err)
	}
	defer rows.Close()

	var following []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Avatar)
		if err != nil {
			return nil, fmt.Errorf("failed to parse following user: %w", err)
		}
		following = append(following, user)
	}

	return following, nil
}
