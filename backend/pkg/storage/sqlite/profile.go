package sqlite

import (
	"log"
	"time"
	
	"socialNetwork/pkg/domain/profile"
	"socialNetwork/pkg/errors"
)

// GetProfileData retrieves a user's profile data
func (s *SQLiteStore) GetProfileData(p *profile.Profile) (*profile.Profile, error) {
	// Query to get basic user data
	query := `SELECT 
		first_name, last_name, nickname, about_me, avatar, date_of_birth, 
		profile_type, created_at 
		FROM users 
		WHERE id = ?`
	
	var createdAt string
	err := s.db.QueryRow(query, p.UserID).Scan(
		&p.FirstName,
		&p.LastName,
		&p.Nickname,
		&p.AboutMe,
		&p.Avatar,
		&p.DateOfBirth,
		&p.ProfileType,
		&createdAt,
	)
	
	if err != nil {
		log.Print("Error retrieving profile data:", err)
		return nil, errors.ErrInternalServer
	}
	
	// Parse created_at into a time.Time
	createdTime, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		log.Print("Error parsing creation time:", err)
		p.CreatedAt = time.Now() // Fallback to current time
	} else {
		p.CreatedAt = createdTime
	}
	
	p.LastActivity = time.Now() // Just use current time for last activity
	
	// Get follower count
	err = s.db.QueryRow("SELECT COUNT(*) FROM followers WHERE followed_id = ?", p.UserID).Scan(&p.Followers)
	if err != nil {
		log.Print("Error getting follower count:", err)
		p.Followers = 0 // Default to 0 on error
	}
	
	// Get following count
	err = s.db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ?", p.UserID).Scan(&p.Following)
	if err != nil {
		log.Print("Error getting following count:", err)
		p.Following = 0 // Default to 0 on error
	}
	
	// Get post count (assuming there is a posts table)
	err = s.db.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", p.UserID).Scan(&p.PostCount)
	if err != nil {
		log.Print("Error getting post count:", err)
		p.PostCount = 0 // Default to 0 on error
	}
	
	return p, nil
}

// SetProfilePrivacy updates a user's profile privacy setting
func (s *SQLiteStore) SetProfilePrivacy(p profile.Profile, privacyType string) error {
	_, err := s.db.Exec("UPDATE users SET profile_type = ? WHERE id = ?", privacyType, p.UserID)
	if err != nil {
		log.Print("Error updating profile privacy:", err)
		return errors.ErrInternalServer
	}
	return nil
}