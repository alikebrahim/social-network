package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/domain/auth"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/logger"
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

// validateEmail performs basic email validation
func validateEmail(email string) error {
	if email == "" {
		return errors.ErrInvalidEmail
	}
	
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.ErrInvalidEmail
	}
	
	return nil
}

// validatePassword performs basic password validation
func validatePassword(password string) error {
	if password == "" {
		return errors.ErrPasswordTooWeak
	}
	
	if len(password) < 8 {
		return errors.ErrPasswordTooWeak.WithDetails("Password must be at least 8 characters long")
	}
	
	return nil
}

// HandleRegister handles user registration
func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	log := logger.FromRequest(r)
	
	// Parse the request body
	var req auth.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Failed to decode registration request", "error", err)
		return errors.ErrBadRequest.WithDetails("Invalid request format")
	}

	// Validate the user input
	validationErrors := make(errors.FieldErrors)
	
	// Validate email
	if err := validateEmail(req.Email); err != nil {
		validationErrors["email"] = "Invalid email format"
	}
	
	// Validate password
	if err := validatePassword(req.Password); err != nil {
		validationErrors["password"] = "Password must be at least 8 characters long"
	}
	
	// Validate name fields
	if req.First_name == "" {
		validationErrors["first_name"] = "First name is required"
	}
	
	if req.Last_name == "" {
		validationErrors["last_name"] = "Last name is required"
	}
	
	// Validate profile type
	if req.Profile_type != "public" && req.Profile_type != "private" {
		validationErrors["profile_type"] = "Profile type must be 'public' or 'private'"
	}
	
	// If there are validation errors, return them
	if len(validationErrors) > 0 {
		return errors.NewValidationError("Registration validation failed", validationErrors)
	}

	// Create user account object from request
	userAccount := &auth.UserAccount{
		Email:        req.Email,
		Password:     req.Password,
		First_name:   req.First_name,
		Last_name:    req.Last_name,
		Profile_type: req.Profile_type,
	}
	
	if req.Date_of_birth != "" {
		userAccount.Date_of_birth = req.Date_of_birth
	}

	// Create the user account
	userID, err := h.store.CreateUserAccount(userAccount)
	if err != nil {
		log.Error("Failed to create user account", "error", err)
		
		if strings.Contains(err.Error(), "already exists") {
			return errors.ErrUserAlreadyExists
		}
		
		return errors.Wrap(err, "Failed to create user account")
	}

	// Set the ID and return the user account (without password)
	userAccount.ID = userID
	userAccount.Password = ""

	return utils.WriteJson(w, http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"user":    userAccount,
	})
}

// HandleLogin handles user login
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	log := logger.FromRequest(r)
	
	// Parse the request body
	var loginRequest auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Error("Failed to decode login request", "error", err)
		return errors.ErrBadRequest.WithDetails("Invalid request format")
	}

	// Validate input
	validationErrors := make(errors.FieldErrors)
	
	if loginRequest.Email == "" {
		validationErrors["email"] = "Email is required"
	}
	
	if loginRequest.Password == "" {
		validationErrors["password"] = "Password is required"
	}
	
	if len(validationErrors) > 0 {
		return errors.NewValidationError("Login validation failed", validationErrors)
	}

	// Authenticate the user
	sessionToken, err := h.store.AuthenticateUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		log.Info("Authentication failed", "email", loginRequest.Email, "error", err)
		
		// Convert domain errors to app errors
		if err == auth.ErrInvalidCredentials {
			return errors.ErrInvalidCredentials
		}
		
		return errors.Wrap(err, "Authentication failed")
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
		log.Error("Failed to get user ID from session", "error", err)
		return errors.Wrap(err, "Failed to get user session")
	}

	// Create the response
	response := map[string]interface{}{
		"message": "Login successful",
		"user_id": userID,
		"token":   sessionToken,
	}

	return utils.WriteJson(w, http.StatusOK, response)
}

// HandleLogout handles user logout
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) error {
	log := logger.FromRequest(r)
	
	// Get the session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// No cookie, no session to logout from
			return utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Already logged out"})
		}
		log.Error("Failed to get session cookie", "error", err)
		return errors.Wrap(err, "Failed to get session cookie")
	}

	// Delete the session
	sessionToken := cookie.Value
	err = h.store.DeleteSession(sessionToken)
	if err != nil {
		log.Error("Failed to delete session", "error", err)
		return errors.Wrap(err, "Failed to delete session")
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