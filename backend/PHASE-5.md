# PHASE 5: Logger Design

## Overview
This phase defines the logger interface, configuration options, and overall design for the logging system. The goal is to create a flexible, easy-to-use logging system based on Go's built-in log package.

## Checklist

### 1. Logger Interface Definition
- [ ] Define the core Logger interface
- [ ] Define log levels (DEBUG, INFO, WARN, ERROR)
- [ ] Design methods for adding context to logs
- [ ] Design package-specific logger capabilities

### 2. Configuration Options
- [ ] Design configuration struct
- [ ] Define file output options
- [ ] Define console output options
- [ ] Create log level configuration

### 3. Request Logging Design
- [ ] Design HTTP middleware for request logging
- [ ] Define request ID generation approach
- [ ] Design request context integration

### 4. Error Logging Design
- [ ] Design error handling integration
- [ ] Define error log format
- [ ] Plan integration with existing error system

### 5. Format and Output
- [ ] Define log entry format
- [ ] Design timestamp format
- [ ] Plan file rotation approach
- [ ] Define colored console output approach

## Logger Interface Design

```go
// Level defines the log levels
type Level int

const (
    DEBUG Level = iota
    INFO
    WARN
    ERROR
)

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

// HTTP middleware interface
func RequestLogger(next http.Handler) http.Handler

// Helper to get logger from request
func FromRequest(r *http.Request) Logger

// Initialize the logger
func Init(config Config)

// Get a logger for a specific package
func GetLogger(pkg string) Logger
```

## Log Format Design

Each log entry will include:
1. Timestamp (ISO 8601 format)
2. Log level (DEBUG, INFO, WARN, ERROR)
3. Package name
4. File and line number
5. Message
6. Contextual fields (key-value pairs)

Example format:
```
2025-04-15 19:45:32.123 [INFO] [api/server.go:45] Server started on port 3000 {host=localhost, port=3000}
```

## Request Logging Design

HTTP requests will be logged with:
1. Request ID (UUID v4)
2. HTTP method
3. URL path
4. Status code
5. Response time
6. User agent (optional)
7. IP address (optional)

Example:
```
2025-04-15 19:45:35.234 [INFO] [api/middleware.go:30] Request processed {request_id=550e8400-e29b-41d4-a716-446655440000, method=GET, path=/profiles/1, status=200, duration=15ms}
```

## Error Logging Design

Errors will be logged with:
1. Error message
2. Error context (file, line)
3. Stack trace (optional for ERROR level)
4. Additional context from the Logger fields

Example:
```
2025-04-15 19:45:40.567 [ERROR] [storage/sqlite/auth.go:25] Failed to authenticate user {user_id=123, error=record not found}
```

## Implementation Considerations

1. **Thread Safety**: The logger will be thread-safe using mutex locks
2. **Performance**: Minimize allocations and formatting overhead
3. **Flexibility**: Allow context to be added at multiple levels
4. **Simplicity**: Keep the API simple and intuitive
5. **Standard Library**: Use only the built-in `log` package without external dependencies

## Next Steps

Once the design is finalized, the next phase will implement the logger according to this design in the `pkg/logger` package.