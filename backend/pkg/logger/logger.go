package logger

import (
	"context"
	"net/http"
	"strings"
)

// Level defines the log levels
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LevelFromString converts a string level to Level type
func LevelFromString(level string) Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

// Config defines the configuration for the logger
type Config struct {
	Level     Level
	LogFile   string
	ToConsole bool
	UseColors bool // Enable colored output in console
}

// Logger defines the interface for the logger
type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithPackage(pkg string) Logger
	Close() error // Close any open resources
}

type contextKey struct {
	name string
}

var loggerContextKey = &contextKey{"logger"}

// FromRequest retrieves the logger from the request context
func FromRequest(r *http.Request) Logger {
	logger, ok := r.Context().Value(loggerContextKey).(Logger)
	if !ok {
		return GetLogger("unknown")
	}
	return logger
}

// WithContext adds a logger to a context
func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// GetLogger returns a logger instance for the given package
func GetLogger(pkg string) Logger {
	return defaultLogger.WithPackage(pkg)
}

// RequestID returns the request ID from the context
func RequestID(r *http.Request) string {
	id, ok := r.Context().Value(requestIDKey).(string)
	if !ok {
		return "no-request-id"
	}
	return id
}

// Init initializes the logger with the given configuration
func Init(config Config) {
	defaultLogger = newStdLogger(config)
}

// RequestMiddleware creates a middleware that logs HTTP requests
// Note: This is an alias for RequestLogger to match the PHASE-7 documentation
func RequestMiddleware() func(http.Handler) http.Handler {
	return RequestLogger
}