package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	
	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/logger"
	"socialNetwork/pkg/storage"
)

type userAuthKey struct{}

var UserAuthKey = &userAuthKey{}

func RequestWithUserID(r *http.Request, userID int64) *http.Request {
	// Create a shallow copy of the request
	newRequest := *r
	// Create a shallow copy of the URL
	if r.URL != nil {
		url := *r.URL
		newRequest.URL = &url
	}
	// Store user ID directly in the request
	newRequest.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
	return &newRequest
}

func AuthMiddleware(next http.Handler, store storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromRequest(r)

		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				utils.WriteError(w, errors.ErrUnauthorized.WithDetails("No session cookie found"))
				return
			}
			utils.WriteError(w, errors.Wrap(err, "Failed to read session cookie"))
			return
		}

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

		// Create a new request with the user ID in header instead of context
		newReq := RequestWithUserID(r, userID)
		next.ServeHTTP(w, newReq)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a logger with request info
		log := logger.New().With("request_path", r.URL.Path).With("method", r.Method)
		
		// Store logger ID in request header
		requestID := uuid.New().String()
		r.Header.Set("X-Request-ID", requestID)
		logger.StoreLogger(requestID, log)
		
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}