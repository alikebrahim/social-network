package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"unicode"

	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/domain/auth"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/logger"
	"socialNetwork/pkg/storage"
)

type AuthHandler struct {
	store storage.Storage
}

func NewAuthHandler(store storage.Storage) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

func validateEmail(email string) error {
	if email == "" {
		return errors.ErrInvalidEmail
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return errors.ErrInvalidEmail.WithDetails("Email must contain a username and domain separated by @")
	}
	
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return errors.ErrPasswordTooWeak
	}
	
	if len(password) < 6 {
		return errors.ErrPasswordTooWeak.WithDetails("Password must be at least 6 characters long")
	}
	
	return nil
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	log := logger.FromRequest(r)
	
	var req auth.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Failed to decode registration request", "error", err)
		return errors.ErrBadRequest.WithDetails("Invalid request format")
	}

	validationErrors := make(errors.FieldErrors)
	
	if err := validateEmail(req.Email); err != nil {
		errDetails, ok := err.(*errors.AppError)
		if ok && errDetails.Details != nil {
			validationErrors["email"] = errDetails.Details.(string)
		} else {
			validationErrors["email"] = "Invalid email format"
		}
	}
	
	if err := validatePassword(req.Password); err != nil {
		errDetails, ok := err.(*errors.AppError)
		if ok && errDetails.Details != nil {
			validationErrors["password"] = errDetails.Details.(string)
		} else {
			validationErrors["password"] = "Password must meet security requirements"
		}
	}
	
	if req.First_name == "" {
		validationErrors["first_name"] = "First name is required"
	}
	
	if req.Last_name == "" {
		validationErrors["last_name"] = "Last name is required"
	}
	
	if req.Profile_type != "public" && req.Profile_type != "private" {
		validationErrors["profile_type"] = "Profile type must be 'public' or 'private'"
	}
	
	if len(validationErrors) > 0 {
		return errors.NewValidationError("Registration validation failed", validationErrors)
	}

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

	userID, err := h.store.CreateUserAccount(userAccount)
	if err != nil {
		log.Error("Failed to create user account", "error", err)
		
		if strings.Contains(err.Error(), "already exists") {
			return errors.ErrUserAlreadyExists
		}
		
		return errors.Wrap(err, "Failed to create user account")
	}

	userAccount.ID = userID
	userAccount.Password = ""

	return utils.WriteJson(w, http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"user":    userAccount,
	})
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	log := logger.FromRequest(r)
	
	var loginRequest auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Error("Failed to decode login request", "error", err)
		return errors.ErrBadRequest.WithDetails("Invalid request format")
	}

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

	sessionToken, err := h.store.AuthenticateUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		log.Info("Authentication failed", "email", loginRequest.Email, "error", err)
		
		if err == auth.ErrInvalidCredentials {
			return errors.ErrInvalidCredentials
		}
		
		return errors.Wrap(err, "Authentication failed")
	}

	expiry := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expiry,
		MaxAge:   int(24 * time.Hour.Seconds()),
	}
	http.SetCookie(w, &cookie)

	userID, err := h.store.GetUserIdBySession(sessionToken)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.Wrap(err, "Failed to get user session")
	}

	response := map[string]interface{}{
		"message": "Login successful",
		"user_id": userID,
		"token":   sessionToken,
	}

	return utils.WriteJson(w, http.StatusOK, response)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) error {
	log := logger.FromRequest(r)
	
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Already logged out"})
		}
		log.Error("Failed to get session cookie", "error", err)
		return errors.Wrap(err, "Failed to get session cookie")
	}

	sessionToken := cookie.Value
	err = h.store.DeleteSession(sessionToken)
	if err != nil {
		log.Error("Failed to delete session", "error", err)
		return errors.Wrap(err, "Failed to delete session")
	}

	cookie = &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	return utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Logout successful"})
}