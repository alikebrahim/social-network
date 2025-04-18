package sqlite

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"socialNetwork/pkg/domain/auth"
	"socialNetwork/pkg/errors"
)

// CreateUserAccount creates a new user account
func (s *SQLiteStore) CreateUserAccount(account *auth.UserAccount) (int64, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print("Failed to hash password", err)
		return 0, errors.ErrInternalServer
	}

	// Check if email already exists
	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", account.Email).Scan(&count)
	if err != nil {
		log.Print("Error checking if email exists:", err)
		return 0, errors.ErrInternalServer
	}
	if count > 0 {
		return 0, errors.ErrAlreadyExists
	}

	// Insert the new user
	query := `INSERT INTO users (email, password, first_name, last_name, date_of_birth, 
              avatar, nickname, about_me, profile_type) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := s.db.Exec(
		query,
		account.Email,
		string(hashedPassword),
		account.First_name,
		account.Last_name,
		account.Date_of_birth,
		account.Avatar,
		account.Nickname,
		account.About_me,
		account.Profile_type,
	)
	if err != nil {
		log.Print("Error creating user account:", err)
		return 0, errors.ErrInternalServer
	}

	// Get the ID of the newly created user
	userID, err := result.LastInsertId()
	if err != nil {
		log.Print("Error getting last insert ID:", err)
		return 0, errors.ErrInternalServer
	}

	return userID, nil
}

// AuthenticateUser authenticates a user and returns a session token
func (s *SQLiteStore) AuthenticateUser(email, password string) (string, error) {
	var userID int64
	var hashedPassword string

	// Get the user by email
	err := s.db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.ErrNotFound
		}
		log.Print("Error retrieving user:", err)
		return "", errors.ErrInternalServer
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", errors.ErrUnauthorized
	}

	// Generate a session token and store it
	// This is a simplified version; in a real-world app you would use a more secure method
	sessionToken := generateSessionToken()
	expireAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	// Store the session
	_, err = s.db.Exec("INSERT INTO sessions (id, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)",
		sessionToken, userID, time.Now().Format(time.RFC3339), expireAt)
	if err != nil {
		log.Print("Error creating session:", err)
		return "", errors.ErrInternalServer
	}

	return sessionToken, nil
}


// DeleteSession removes a session (logout)
func (s *SQLiteStore) DeleteSession(sessionToken string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE id = ?", sessionToken)
	if err != nil {
		log.Print("Error deleting session:", err)
		return errors.ErrInternalServer
	}
	return nil
}

// GetUserIdBySession retrieves a user ID from a session token
func (s *SQLiteStore) GetUserIdBySession(sessionToken string) (int64, error) {
	var userID int64
	var expiresAt string

	// Get the user associated with the session
	err := s.db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE id = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.ErrNotFound
		}
		log.Print("Error retrieving user from session:", err)
		return 0, errors.ErrInternalServer
	}

	// Check if session is expired
	expiryTime, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil || expiryTime.Before(time.Now()) {
		// Delete the expired session
		s.DeleteSession(sessionToken)
		return 0, errors.ErrUnauthorized
	}

	return userID, nil
}

// Helper function to generate a random session token using UUID
func generateSessionToken() string {
	// Generate a secure random UUID for session tokens
	return uuid.New().String()
}