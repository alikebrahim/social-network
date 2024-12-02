package services

import (
	"backend/pkg/utils"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUser registers a new user by hashing the password and saving the user details.
func RegisterUser(db *sql.DB, email, password, firstName, lastName, dateOfBirth, avatar, nickname, aboutMe string, isPublic bool) (string, error) {
	// Check if the email already exists in the database
	var existingUserID string
	query := `SELECT id FROM users WHERE email = ?`
	err := db.QueryRow(query, email).Scan(&existingUserID)
	if err == nil {
		return "", errors.New("email already in use")
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to check email availability: %w", err)
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate a new user ID and insert the user data into the database
	userID := utils.GenerateUUID() // Generate a new UUID for the user
	createdAt := time.Now()

	query = `INSERT INTO users (id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_public, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, userID, email, string(hashedPassword), firstName, lastName, dateOfBirth, avatar, nickname, aboutMe, isPublic, createdAt)
	if err != nil {
		return "", fmt.Errorf("failed to insert user into database: %w", err)
	}

	return userID, nil
}

// LoginUser logs a user in by verifying the provided password and returning a JWT token.
func LoginUser(db *sql.DB, email, password string) (string, error) {
	// Retrieve the user from the database
	var userID, hashedPassword string
	query := `SELECT id, password FROM users WHERE email = ?`
	err := db.QueryRow(query, email).Scan(&userID, &hashedPassword)
	if err == sql.ErrNoRows {
		return "", errors.New("invalid email or password")
	}
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user: %w", err)
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate a JWT token for the user
	token, err := utils.GenerateJWT(userID)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return token, nil
}

// VerifyToken verifies the provided JWT token and returns the user ID if valid.
func VerifyToken(db *sql.DB, token string) (string, error) {
	// Validate the token and extract the user ID from the claims
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return "", fmt.Errorf("invalid or expired token: %w", err)
	}

	// Retrieve the user from the database to ensure the user exists
	userID := claims["user_id"].(string)
	var existingUserID string
	query := `SELECT id FROM users WHERE id = ?`
	err = db.QueryRow(query, userID).Scan(&existingUserID)
	if err == sql.ErrNoRows {
		return "", errors.New("user not found")
	}
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user: %w", err)
	}

	return userID, nil
}
