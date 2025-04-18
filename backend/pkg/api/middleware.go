package api

import (
	"context"
	"net/http"

	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/logger"
	"socialNetwork/pkg/storage"
)

// AuthMiddleware is a middleware that checks if the user is authenticated
func AuthMiddleware(next http.Handler, store storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get a logger
		log := logger.FromRequest(r)

		// Get the session token from the cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				utils.WriteError(w, errors.ErrUnauthorized.WithDetails("No session cookie found"))
				return
			}
			utils.WriteError(w, errors.Wrap(err, "Failed to read session cookie"))
			return
		}

		// Verify the session token
		sessionToken := cookie.Value
		userID, err := store.GetUserIdBySession(sessionToken)
		if err != nil {
			if errors.IsUnauthorized(err) || errors.IsNotFound(err) {
				utils.WriteError(w, errors.ErrSessionExpired)
			} else {
				log.Error("Session validation error", "error", err)
				utils.WriteError(w, errors.Wrap(err, "Failed to validate session"))
			}
			return
		}

		// Add user ID to the request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		
		// Session is valid, call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggerMiddleware adds a logger to the request context
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a logger for this request
		log := logger.New().With("request_path", r.URL.Path).With("method", r.Method)
		
		// Add the logger to the request context
		ctx := logger.WithLogger(r.Context(), log)
		
		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS)
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}