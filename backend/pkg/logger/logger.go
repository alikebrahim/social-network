package logger

import (
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
	With(key string, value interface{}) Logger  // Alias for WithField for easier chaining
	WithFields(fields map[string]interface{}) Logger
	WithPackage(pkg string) Logger
	Close() error // Close any open resources
}

var (
	// Global logger storage map - in a real app, you'd use a concurrency-safe map
	loggerMap = make(map[string]Logger)
	defaultLogger Logger
)

// StoreLogger stores a logger with the given request ID
func StoreLogger(requestID string, logger Logger) {
	loggerMap[requestID] = logger
}

// FromRequest retrieves the logger from the request header
func FromRequest(r *http.Request) Logger {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		return GetLogger("unknown")
	}
	
	logger, ok := loggerMap[requestID]
	if !ok {
		return GetLogger("unknown")
	}
	
	return logger
}

// GetLogger returns a logger instance for the given package
func GetLogger(pkg string) Logger {
	return defaultLogger.WithPackage(pkg)
}

// New returns a new logger instance
func New() Logger {
	return defaultLogger
}

// RequestID returns the request ID from the header
func RequestID(r *http.Request) string {
	id := r.Header.Get("X-Request-ID")
	if id == "" {
		return "no-request-id"
	}
	return id
}

// Init initializes the logger with the given configuration
func Init(config Config) {
	defaultLogger = newStdLogger(config)
}

// RequestMiddleware creates a middleware that logs HTTP requests
func RequestMiddleware() func(http.Handler) http.Handler {
	return RequestLogger
}