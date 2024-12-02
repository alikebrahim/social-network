package middleware

import (
	"log"
	"net/http"
	"time"
)

// RequestLogger logs incoming requests and their processing times
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time
		startTime := time.Now()

		// Log the incoming request details
		log.Printf("Incoming request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Create a response writer wrapper to capture the status code
		wrapper := &responseWriterWrapper{ResponseWriter: w}
		next.ServeHTTP(wrapper, r)

		// Log the details after the response is written
		duration := time.Since(startTime)
		log.Printf("Request processed: %s %s with status %d in %v", r.Method, r.URL.Path, wrapper.statusCode, duration)
	})
}

// responseWriterWrapper is a custom wrapper around http.ResponseWriter
// to capture the status code of the response.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code when it is written
func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the status code and then calls the original Write method
func (rw *responseWriterWrapper) Write(p []byte) (int, error) {
	if rw.statusCode == 0 {
		// If no status code is set, default to 200
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(p)
}
