package handlers

import (
	"encoding/json"
	"net/http"

	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/domain/auth"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/storage"
)

// AuthHandler manages authentication-related handlers
type AuthHandler struct {
	store storage.Storage
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(store storage.Storage) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

// HandleRegister handles user registration
func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	// Parse the request body
	var userAccount auth.UserAccount
	if err := json.NewDecoder(r.Body).Decode(&userAccount); err != nil {
		return errors.ErrBadRequest
	}

	// Validate the user input
	// This could be more detailed in a production application
	if userAccount.Email == "" || userAccount.Password == "" {
		return errors.ErrInvalidInput
	}

	// Create the user account
	userID, err := h.store.CreateUserAccount(&userAccount)
	if err != nil {
		return err
	}

	// Set the ID and return the user account (without password)
	userAccount.ID = userID
	userAccount.Password = ""

	return utils.WriteJson(w, http.StatusCreated, userAccount)
}

// HandleLogin handles user login
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	// Parse the request body
	var loginRequest auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		return errors.ErrBadRequest
	}

	// Validate input
	if loginRequest.Email == "" || loginRequest.Password == "" {
		return errors.ErrInvalidInput
	}

	// Authenticate the user
	sessionToken, err := h.store.AuthenticateUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return err
	}

	// Set session cookie
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		// Secure: true, // Uncomment in production with HTTPS
	}
	http.SetCookie(w, &cookie)

	// Get the user ID from the session
	userID, err := h.store.GetUserIdBySession(sessionToken)
	if err != nil {
		return err
	}

	// Create the response
	response := map[string]interface{}{
		"message":   "Login successful",
		"user_id":   userID,
		"token":     sessionToken,
	}

	return utils.WriteJson(w, http.StatusOK, response)
}

// HandleLogout handles user logout
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) error {
	// Get the session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// No cookie, no session to logout from
			return utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Already logged out"})
		}
		return errors.ErrInternalServer
	}

	// Delete the session
	sessionToken := cookie.Value
	err = h.store.DeleteSession(sessionToken)
	if err != nil {
		return err
	}

	// Clear the cookie
	cookie = &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		// Secure: true, // Uncomment in production with HTTPS
	}
	http.SetCookie(w, cookie)

	return utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Logout successful"})
}