package utils

import (
	"encoding/json"
	stdErrors "errors"
	"log"
	"net/http"
	"strconv"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/logger"
)

type ApiError struct {
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("json.NewEncoder error:", err)
		return err
	}
	return nil
}

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

func GetUserIDFromContext(r *http.Request) (int64, error) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0, errors.ErrUnauthorized
	}
	
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "Invalid user ID format")
	}
	
	return userID, nil
}

func GetParam(r *http.Request, key string) string {
	return r.PathValue(key)
}