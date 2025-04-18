package utils

import (
	"encoding/json"
	stdErrors "errors"
	"log"
	"net/http"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/logger"
)

// ApiError represents an API error response
type ApiError struct {
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

// WriteJson sends a JSON response
func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("json.NewEncoder error:", err)
		return err
	}
	return nil
}

// WriteError sends an error response with the appropriate status code
func WriteError(w http.ResponseWriter, err error) error {
	// Get the status code from the error, or default to 500
	statusCode := errors.HTTPStatusFromError(err)

	// If it's an AppError, extract the error details
	var appErr *errors.AppError
	if stdErrors.As(err, &appErr) {
		return WriteJson(w, statusCode, ApiError{
			Message: appErr.Message,
			Code:    appErr.Code,
			Details: appErr.Details,
		})
	}

	// Default error response for standard errors
	return WriteJson(w, statusCode, ApiError{
		Message: err.Error(),
		Code:    "internal_error",
	})
}

// MakeHTTPHandleFunc is a helper that allows handlers to return errors
func MakeHTTPHandleFunc(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a logger from the request context or create a new one
		log := logger.FromRequest(r)

		if err := f(w, r); err != nil {
			// Log the error
			log.Error("Handler error", "error", err)
			
			// Write the error response
			if writeErr := WriteError(w, err); writeErr != nil {
				log.Error("Failed to write error response", "error", writeErr)
			}
		}
	}
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(r *http.Request) (int64, error) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		return 0, errors.ErrUnauthorized
	}
	return userID, nil
}

// GetParam retrieves a URL parameter from the request
func GetParam(r *http.Request, key string) string {
	// Try to get path parameters from the request context
	params, ok := r.Context().Value("params").(map[string]string)
	if ok {
		if value, exists := params[key]; exists {
			return value
		}
	}
	
	// Try to get from path value (for Go 1.22+ compatibility)
	return r.PathValue(key)
}