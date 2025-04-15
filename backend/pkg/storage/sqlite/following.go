package sqlite

import (
	"log"
	"time"

	"socialNetwork/pkg/domain/auth"
	"socialNetwork/pkg/domain/following"
	"socialNetwork/pkg/errors"
)

// CreateFollowRequest creates a new follow request
func (s *SQLiteStore) CreateFollowRequest(req following.FollowRequest) error {
	// Check if the follow relationship already exists
	query := `SELECT COUNT(*) FROM followers 
              WHERE follower_id = ? AND followed_id = ?`

	var count int
	err := s.db.QueryRow(query, req.FollowerID, req.FollowedID).Scan(&count)
	if err != nil {
		log.Print("Error checking existing follow relationship:", err)
		return errors.ErrInternalServer
	}

	if count > 0 {
		return errors.ErrAlreadyExists
	}

	// Check if the target profile is private
	var profileType string
	profileQuery := `SELECT profile_type FROM users WHERE id = ?`
	err = s.db.QueryRow(profileQuery, req.FollowedID).Scan(&profileType)
	if err != nil {
		log.Print("Error checking profile type:", err)
		return errors.ErrInternalServer
	}

	// Set status based on profile type
	var status string
	if profileType == "public" {
		status = "accepted"
	} else {
		status = "pending"
	}

	// Create the follow request
	insertQuery := `INSERT INTO followers (follower_id, followed_id, status, created_at) 
                   VALUES (?, ?, ?, ?)`
	_, err = s.db.Exec(insertQuery, req.FollowerID, req.FollowedID, status, time.Now())
	if err != nil {
		log.Print("Error creating follow request:", err)
		return errors.ErrInternalServer
	}

	return nil
}

// AcceptFollowRequest accepts a follow request
func (s *SQLiteStore) AcceptFollowRequest(followerID, followedID int64) error {
	// Check if the request exists
	query := `SELECT COUNT(*) FROM followers 
              WHERE follower_id = ? AND followed_id = ? AND status = 'pending'`

	var count int
	err := s.db.QueryRow(query, followerID, followedID).Scan(&count)
	if err != nil {
		log.Print("Error checking follow request:", err)
		return errors.ErrInternalServer
	}

	if count == 0 {
		return errors.ErrNotFound
	}

	// Update to accepted status
	updateQuery := `UPDATE followers SET status = 'accepted' 
                   WHERE follower_id = ? AND followed_id = ?`
	_, err = s.db.Exec(updateQuery, followerID, followedID)
	if err != nil {
		log.Print("Error accepting follow request:", err)
		return errors.ErrInternalServer
	}

	return nil
}

// DeleteFollowRequest removes a follow request or unfollow a user
func (s *SQLiteStore) DeleteFollowRequest(followerID, followedID int64) error {
	// Delete the follow relationship or request
	deleteQuery := `DELETE FROM followers 
                   WHERE follower_id = ? AND followed_id = ?`
	result, err := s.db.Exec(deleteQuery, followerID, followedID)
	if err != nil {
		log.Print("Error removing follow relationship:", err)
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

// GetFollowRequests gets all pending follow requests for a user
func (s *SQLiteStore) GetFollowRequests(userID int64) ([]auth.UserAccount, error) {
	query := `SELECT u.id, u.email, u.first_name, u.last_name, u.date_of_birth, 
              u.avatar, u.nickname, u.about_me, u.profile_type
              FROM users u
              JOIN followers f ON u.id = f.follower_id
              WHERE f.followed_id = ? AND f.status = 'pending'
              ORDER BY f.created_at DESC`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		log.Print("Error querying follow requests:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	var usersList []auth.UserAccount
	for rows.Next() {
		var user auth.UserAccount
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.First_name,
			&user.Last_name,
			&user.Date_of_birth,
			&user.Avatar,
			&user.Nickname,
			&user.About_me,
			&user.Profile_type,
		)
		if err != nil {
			log.Print("Error scanning follow request row:", err)
			return nil, errors.ErrInternalServer
		}
		// Don't expose the password
		user.Password = ""
		usersList = append(usersList, user)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating follow request rows:", err)
		return nil, errors.ErrInternalServer
	}

	return usersList, nil
}