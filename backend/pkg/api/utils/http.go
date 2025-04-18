package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// ApiError represents an API error response
type ApiError struct {
	Error string
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

// MakeHTTPHandleFunc is a helper that allows handlers to return errors
func MakeHTTPHandleFunc(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(r *http.Request) (int64, error) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		return 0, errors.New("user ID not found in context")
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