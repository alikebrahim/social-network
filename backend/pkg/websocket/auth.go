package websocket

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func getUserIDFromSession(sessionToken string) (int, error) {
	var userID int

	err := DB.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hashPassword(password string) (string, error) {
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedpassword), nil
}
