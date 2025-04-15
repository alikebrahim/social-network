# PHASE 6: Logger Implementation

## Overview
This phase implements the logger package based on the design from Phase 5. The implementation will use Go's built-in log package to create a flexible, simple logging system.

## Checklist

### 1. Core Logger Implementation
- [ ] Create the logger package structure
- [ ] Implement log levels
- [ ] Create the base logger implementation
- [ ] Implement log formatting

### 2. Configuration and Initialization
- [ ] Implement configuration handling
- [ ] Create file output functionality
- [ ] Implement console output
- [ ] Add initialization function

### 3. Context and Fields
- [ ] Implement WithField method
- [ ] Implement WithFields method
- [ ] Implement package context

### 4. HTTP Request Logging
- [ ] Create request ID generation
- [ ] Implement HTTP middleware
- [ ] Add request context integration

### 5. Testing
- [ ] Create basic unit tests
- [ ] Test different log levels
- [ ] Test field context
- [ ] Test file output

## File Structure

```
/pkg/logger/
  logger.go         - Core logger interface and implementation
  middleware.go     - HTTP middleware for request logging
  request.go        - Request context integration
  config.go         - Configuration handling
```

## Implementation Steps

### 1. Create Basic Package Structure

Create the logger.go file with the core interface and implementation:

```go
package logger

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "runtime"
    "strings"
    "sync"
    "time"
)

// Level defines the log levels
type Level int

const (
    DEBUG Level = iota
    INFO
    WARN
    ERROR
)

// String returns the string representation of the log level
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

// Config defines the configuration for the logger
type Config struct {
    Level     Level
    LogFile   string
    ToConsole bool
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
}

// logger is the implementation of the Logger interface
type logger struct {
    level      Level
    fields     map[string]interface{}
    packageStr string
    mu         sync.Mutex
    output     io.Writer
    stdLogger  *log.Logger
}

var (
    defaultLogger Logger
    once          sync.Once
)
```

### 2. Implement Initialization and Configuration

Add initialization logic to logger.go:

```go
// Init initializes the logger with the given configuration
func Init(config Config) {
    once.Do(func() {
        var writers []io.Writer

        // Add console output if configured
        if config.ToConsole {
            writers = append(writers, os.Stdout)
        }

        // Add file output if configured
        if config.LogFile != "" {
            // Ensure directory exists
            dir := filepath.Dir(config.LogFile)
            if err := os.MkdirAll(dir, 0755); err != nil {
                log.Printf("Failed to create log directory: %v", err)
            } else {
                file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
                if err != nil {
                    log.Printf("Failed to open log file: %v", err)
                } else {
                    writers = append(writers, file)
                }
            }
        }

        // Default to stdout if no writers
        if len(writers) == 0 {
            writers = append(writers, os.Stdout)
        }

        // Create multi-writer if needed
        var output io.Writer
        if len(writers) == 1 {
            output = writers[0]
        } else {
            output = io.MultiWriter(writers...)
        }

        // Create the logger
        stdLogger := log.New(output, "", 0)
        defaultLogger = &logger{
            level:     config.Level,
            fields:    make(map[string]interface{}),
            output:    output,
            stdLogger: stdLogger,
        }
    })
}

// GetLogger returns a logger instance with the given package name
func GetLogger(pkg string) Logger {
    if defaultLogger == nil {
        // Initialize with default config if not already done
        Init(Config{
            Level:     INFO,
            ToConsole: true,
        })
    }
    return defaultLogger.WithPackage(pkg)
}

// LevelFromString converts a string to a Level
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
```

### 3. Implement Logging Methods

Add the core logging methods to logger.go:

```go
// log logs a message at the given level with the given fields
func (l *logger) log(level Level, msg string, keyvals ...interface{}) {
    if level < l.level {
        return
    }

    l.mu.Lock()
    defer l.mu.Unlock()

    // Get caller information
    _, file, line, ok := runtime.Caller(2)
    fileLine := ""
    if ok {
        file = filepath.Base(file)
        fileLine = fmt.Sprintf("%s:%d", file, line)
    }

    // Format timestamp
    timestamp := time.Now().Format("2006-01-02 15:04:05.000")

    // Format package
    packageStr := ""
    if l.packageStr != "" {
        packageStr = fmt.Sprintf("[%s] ", l.packageStr)
    }

    // Format level
    levelStr := fmt.Sprintf("[%s]", level.String())

    // Format fields
    fieldsStr := ""
    if len(l.fields) > 0 {
        parts := make([]string, 0, len(l.fields))
        for k, v := range l.fields {
            parts = append(parts, fmt.Sprintf("%s=%v", k, v))
        }
        fieldsStr = fmt.Sprintf(" {%s}", strings.Join(parts, ", "))
    }

    // Format additional key-values
    additionalStr := ""
    if len(keyvals) > 0 && len(keyvals)%2 == 0 {
        parts := make([]string, 0, len(keyvals)/2)
        for i := 0; i < len(keyvals); i += 2 {
            key := fmt.Sprintf("%v", keyvals[i])
            value := keyvals[i+1]
            parts = append(parts, fmt.Sprintf("%s=%v", key, value))
        }
        additionalStr = fmt.Sprintf(" (%s)", strings.Join(parts, ", "))
    }

    // Build log entry
    logEntry := fmt.Sprintf("%s %s %s%s %s%s%s",
        timestamp,
        levelStr,
        packageStr,
        fileLine,
        msg,
        fieldsStr,
        additionalStr,
    )

    l.stdLogger.Println(logEntry)
}

// Debug logs a debug message
func (l *logger) Debug(msg string, keyvals ...interface{}) {
    l.log(DEBUG, msg, keyvals...)
}

// Info logs an info message
func (l *logger) Info(msg string, keyvals ...interface{}) {
    l.log(INFO, msg, keyvals...)
}

// Warn logs a warning message
func (l *logger) Warn(msg string, keyvals ...interface{}) {
    l.log(WARN, msg, keyvals...)
}

// Error logs an error message
func (l *logger) Error(msg string, keyvals ...interface{}) {
    l.log(ERROR, msg, keyvals...)
}
```

### 4. Implement Context Methods

Add context methods to logger.go:

```go
// WithField returns a new logger with the given field
func (l *logger) WithField(key string, value interface{}) Logger {
    return l.WithFields(map[string]interface{}{key: value})
}

// WithFields returns a new logger with the given fields
func (l *logger) WithFields(fields map[string]interface{}) Logger {
    newLogger := &logger{
        level:      l.level,
        fields:     make(map[string]interface{}),
        packageStr: l.packageStr,
        output:     l.output,
        stdLogger:  l.stdLogger,
    }

    // Copy existing fields
    for k, v := range l.fields {
        newLogger.fields[k] = v
    }

    // Add new fields
    for k, v := range fields {
        newLogger.fields[k] = v
    }

    return newLogger
}

// WithPackage returns a new logger with the given package
func (l *logger) WithPackage(pkg string) Logger {
    newLogger := &logger{
        level:      l.level,
        fields:     make(map[string]interface{}),
        packageStr: pkg,
        output:     l.output,
        stdLogger:  l.stdLogger,
    }

    // Copy existing fields
    for k, v := range l.fields {
        newLogger.fields[k] = v
    }

    return newLogger
}
```

### 5. Implement HTTP Middleware

Create a middleware.go file for HTTP request logging:

```go
package logger

import (
    "context"
    "net/http"
    "time"

    "github.com/google/uuid"
)

type contextKey string

const (
    requestIDKey contextKey = "request_id"
    loggerKey    contextKey = "logger"
)

// RequestMiddleware creates middleware for logging HTTP requests
func RequestMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            // Generate request ID
            requestID := uuid.New().String()

            // Add request ID to response headers
            w.Header().Set("X-Request-ID", requestID)

            // Create request-specific logger
            reqLogger := GetLogger("http").WithField("request_id", requestID)

            // Add logger to request context
            ctx := context.WithValue(r.Context(), loggerKey, reqLogger)
            ctx = context.WithValue(ctx, requestIDKey, requestID)
            r = r.WithContext(ctx)

            // Create a response wrapper to capture status code
            wrapper := newResponseWriter(w)

            // Log the request
            reqLogger.Info("Request started",
                "method", r.Method,
                "path", r.URL.Path,
                "remote_addr", r.RemoteAddr,
                "user_agent", r.UserAgent(),
            )

            // Call the next handler
            next.ServeHTTP(wrapper, r)

            // Log the response
            duration := time.Since(start)
            reqLogger.Info("Request completed",
                "method", r.Method,
                "path", r.URL.Path,
                "status", wrapper.status,
                "duration", duration.Milliseconds(),
                "duration_unit", "ms",
            )
        })
    }
}

// FromRequest gets the logger from the request context
func FromRequest(r *http.Request) Logger {
    logger, ok := r.Context().Value(loggerKey).(Logger)
    if !ok {
        return GetLogger("http")
    }
    return logger
}

// GetRequestID gets the request ID from the request context
func GetRequestID(r *http.Request) string {
    requestID, ok := r.Context().Value(requestIDKey).(string)
    if !ok {
        return ""
    }
    return requestID
}

// responseWriter is a wrapper around http.ResponseWriter to capture the status code
type responseWriter struct {
    http.ResponseWriter
    status int
}

// newResponseWriter creates a new responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
    return &responseWriter{w, http.StatusOK}
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(status int) {
    rw.status = status
    rw.ResponseWriter.WriteHeader(status)
}
```

### 6. Testing

Create a simple test file to verify the logger functionality:

```go
package logger

import (
    "bytes"
    "strings"
    "testing"
)

func TestLogger(t *testing.T) {
    // Create a buffer to capture log output
    var buf bytes.Buffer

    // Initialize logger with the buffer
    logger := &logger{
        level:      DEBUG,
        fields:     make(map[string]interface{}),
        packageStr: "test",
        output:     &buf,
        stdLogger:  log.New(&buf, "", 0),
    }

    // Test debug log
    buf.Reset()
    logger.Debug("Debug message")
    if !strings.Contains(buf.String(), "[DEBUG]") {
        t.Errorf("Expected debug log to contain [DEBUG], got %s", buf.String())
    }

    // Test info log
    buf.Reset()
    logger.Info("Info message")
    if !strings.Contains(buf.String(), "[INFO]") {
        t.Errorf("Expected info log to contain [INFO], got %s", buf.String())
    }

    // Test with fields
    buf.Reset()
    logger.WithField("key", "value").Info("With field")
    if !strings.Contains(buf.String(), "key=value") {
        t.Errorf("Expected log with field to contain key=value, got %s", buf.String())
    }

    // Test log level filtering
    buf.Reset()
    infoLogger := &logger{
        level:      INFO,
        fields:     make(map[string]interface{}),
        packageStr: "test",
        output:     &buf,
        stdLogger:  log.New(&buf, "", 0),
    }
    infoLogger.Debug("Should not log")
    if buf.String() != "" {
        t.Errorf("Expected debug log to be filtered out at INFO level, got %s", buf.String())
    }
}
```

## Implementation Notes

1. The implementation keeps the design simple and uses only the Go standard library.
2. File rotation is not included but could be added in a future enhancement.
3. The middleware uses UUID for request IDs which would require adding a dependency on github.com/google/uuid.
4. Error integration is handled through the Error log level and the ability to add context.

## Next Steps

Once the logger is implemented, the next phase will integrate it into the application components.