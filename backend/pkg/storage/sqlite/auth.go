package sqlite

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"socialNetwork/pkg/domain/auth"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/logger"
)

// CreateUserAccount creates a new user account
func (s *SQLiteStore) CreateUserAccount(account *auth.UserAccount) (int64, error) {
	log := logger.New().With("method", "CreateUserAccount").With("email", account.Email)
	
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password", "error", err)
		return 0, errors.Wrap(err, "Failed to secure password")
	}

	// Check if email already exists
	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", account.Email).Scan(&count)
	if err != nil {
		log.Error("Error checking if email exists", "error", err)
		return 0, errors.Wrap(err, "Failed to check if email exists")
	}
	if count > 0 {
		return 0, errors.ErrUserAlreadyExists
	}

	// Insert the new user
	query := `INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth, 
			  avatar, nickname, about_me, profile_type, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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
		time.Now().Format(time.RFC3339),
	)
	if err != nil {
		log.Error("Error creating user account", "error", err)
		return 0, errors.Wrap(err, "Failed to create user account")
	}

	// Get the ID of the newly created user
	userID, err := result.LastInsertId()
	if err != nil {
		log.Error("Error getting last insert ID", "error", err)
		return 0, errors.Wrap(err, "Failed to get new user ID")
	}

	log.Info("User account created successfully", "user_id", userID)
	return userID, nil
}

// AuthenticateUser authenticates a user and returns a session token
func (s *SQLiteStore) AuthenticateUser(email, password string) (string, error) {
	log := logger.New().With("method", "AuthenticateUser").With("email", email)
	
	var userID int64
	var hashedPassword string

	// Get the user by email
	err := s.db.QueryRow("SELECT id, password_hash FROM users WHERE email = ?", email).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info("User not found")
			return "", errors.ErrInvalidCredentials
		}
		log.Error("Error retrieving user", "error", err)
		return "", errors.Wrap(err, "Failed to retrieve user")
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Info("Invalid password")
		return "", errors.ErrInvalidCredentials
	}

	// Generate a session token
	sessionToken := generateSessionToken()
	expiryTime := time.Now().Add(24 * time.Hour)
	
	// Store the session
	_, err = s.db.Exec(
		"INSERT INTO sessions (token, user_id, expiry_time) VALUES (?, ?, ?)",
		sessionToken, 
		userID, 
		expiryTime.Format(time.RFC3339),
	)
	
	if err != nil {
		log.Error("Error creating session", "error", err)
		return "", errors.Wrap(err, "Failed to create session")
	}

	log.Info("User authenticated successfully", "user_id", userID)
	return sessionToken, nil
}

// DeleteSession removes a session (logout)
func (s *SQLiteStore) DeleteSession(sessionToken string) error {
	log := logger.New().With("method", "DeleteSession")
	
	result, err := s.db.Exec("DELETE FROM sessions WHERE token = ?", sessionToken)
	if err != nil {
		log.Error("Error deleting session", "error", err)
		return errors.Wrap(err, "Failed to delete session")
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Error checking rows affected", "error", err)
		return errors.Wrap(err, "Failed to check if session was deleted")
	}
	
	if rowsAffected == 0 {
		log.Info("Session not found", "token", sessionToken)
		// Don't return an error for non-existent sessions
	} else {
		log.Info("Session deleted successfully")
	}
	
	return nil
}

// GetUserIdBySession retrieves a user ID from a session token
func (s *SQLiteStore) GetUserIdBySession(sessionToken string) (int64, error) {
	log := logger.New().With("method", "GetUserIdBySession")
	
	var userID int64
	var expiresAt string

	// Get the user associated with the session
	err := s.db.QueryRow("SELECT user_id, expiry_time FROM sessions WHERE token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info("Session not found", "token", sessionToken)
			return 0, errors.ErrSessionExpired
		}
		log.Error("Error retrieving user from session", "error", err)
		return 0, errors.Wrap(err, "Failed to retrieve session")
	}

	// Check if session is expired
	expiryTime, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		log.Error("Error parsing expiry time", "error", err)
		return 0, errors.Wrap(err, "Invalid session data")
	}
	
	if expiryTime.Before(time.Now()) {
		log.Info("Session expired", "token", sessionToken, "expiry", expiryTime)
		// Delete the expired session
		s.DeleteSession(sessionToken)
		return 0, errors.ErrSessionExpired
	}

	log.Info("Session validated successfully", "user_id", userID)
	return userID, nil
}

// Helper function to generate a random session token using UUID
func generateSessionToken() string {
	// Generate a secure random UUID for session tokens
	return uuid.New().String()
}