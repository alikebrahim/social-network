package utils

import (
	"encoding/json"
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