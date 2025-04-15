package api

import (
	"net/http"

	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/storage"
)

// AuthMiddleware is a middleware that checks if the user is authenticated
func AuthMiddleware(next http.Handler, store storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session token from the cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
				return
			}
			utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			return
		}

		// Verify the session token
		sessionToken := cookie.Value
		_, err = store.GetUserIdBySession(sessionToken)
		if err != nil {
			if err == errors.ErrNotFound || err == errors.ErrUnauthorized {
				utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			} else {
				utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			}
			return
		}

		// Session is valid, call the next handler
		next.ServeHTTP(w, r)
	})
}