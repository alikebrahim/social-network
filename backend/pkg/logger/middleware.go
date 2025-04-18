package logger

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

// contextKey is defined in logger.go

// RequestLogger is a middleware that logs HTTP requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Generate a request ID
		requestID := uuid.New().String()
		r.Header.Set("X-Request-ID", requestID)

		// Create a logger with request info
		requestLogger := GetLogger("api").WithFields(map[string]interface{}{
			"request_id": requestID,
			"method":     r.Method,
			"path":       r.URL.Path,
			"remote_ip":  r.RemoteAddr,
			"user_agent": r.UserAgent(),
		})

		// Store the logger in our map
		StoreLogger(requestID, requestLogger)

		// Create a response writer that tracks status code
		rw := newResponseWriter(w)

		// Call the next handler
		requestLogger.Info("Request started")
		next.ServeHTTP(rw, r)

		// Log request completion
		duration := time.Since(startTime)
		requestLogger.WithFields(map[string]interface{}{
			"status":   rw.status,
			"duration": duration.Milliseconds(),
		}).Info("Request completed")
		
		// Optionally clean up the logger (in a production app, you'd use a TTL cache)
		// delete(loggerMap, requestID)
	})
}

// responseWriter is a wrapper for http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}