package middleware

import (
	"net/http"
	"strings"
)

// CORS middleware to handle Cross-Origin Resource Sharing (CORS)
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins, for more control modify this list
		allowedOrigins := []string{
			"http://localhost:3000",    // Add your frontend URL here
			"https://yourfrontend.com", // Example for production
		}

		// Allow specific HTTP methods
		allowedMethods := []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		}

		// Allow specific headers
		allowedHeaders := []string{
			"Content-Type", "Authorization", "Origin", "Accept", "X-Requested-With",
		}

		// Check if the request's origin is allowed
		origin := r.Header.Get("Origin")
		allowOrigin := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowOrigin = true
				break
			}
		}

		if allowOrigin {
			// Allow credentials (cookies, HTTP authentication, etc.)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Set allowed origins, methods, and headers
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))
		}

		// Handle preflight requests (OPTIONS)
		if r.Method == "OPTIONS" {
			// If the method is OPTIONS, respond with status 200 and early return
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
