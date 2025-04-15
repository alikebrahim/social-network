# PHASE 7: Logger Integration

## Overview
This phase focuses on integrating the logger into various application components to provide comprehensive logging throughout the codebase. The integration will follow a consistent pattern to ensure uniform logging behavior.

## Checklist

### 1. Application Initialization
- [ ] Add logger initialization to main.go
- [ ] Configure log level from environment or config
- [ ] Set up file logging if needed

### 2. API Layer Integration
- [ ] Add request logging middleware to the HTTP router
- [ ] Update API handlers to use request-specific loggers
- [ ] Add logging to critical API operations
- [ ] Enhance error responses with logging

### 3. Storage Layer Integration
- [ ] Add logging to database initialization
- [ ] Log database operations and performance
- [ ] Enhance error handling with contextual logging
- [ ] Track query execution times

### 4. WebSocket Integration
- [ ] Log WebSocket connections and disconnections
- [ ] Track real-time message flow
- [ ] Add error logging for WebSocket operations
- [ ] Monitor connection health

### 5. Domain Layer Integration
- [ ] Add contextual logging to domain operations
- [ ] Log validation and business rule outcomes
- [ ] Track complex operations across components

## Implementation Steps

### 1. Update Main Application

Modify the main.go file to initialize the logger:

```go
package main

import (
    "log"
    "os"
    "path/filepath"

    "socialNetwork/pkg/api"
    "socialNetwork/pkg/logger"
    "socialNetwork/pkg/storage/sqlite"
)

func main() {
    // Initialize the logger
    logLevel := os.Getenv("LOG_LEVEL")
    if logLevel == "" {
        logLevel = "INFO"
    }

    logFile := os.Getenv("LOG_FILE")
    toConsole := logFile == "" || os.Getenv("LOG_CONSOLE") == "true"

    logger.Init(logger.Config{
        Level:     logger.LevelFromString(logLevel),
        LogFile:   logFile,
        ToConsole: toConsole,
    })

    // Get logger for main package
    log := logger.GetLogger("main")
    log.Info("Application starting")

    // Initialize the SQLite storage
    store, err := sqlite.NewSQLiteStore()
    if err != nil {
        log.Error("Failed to initialize database", "error", err)
        os.Exit(1)
    }
    if err := store.Init(); err != nil {
        log.Error("Failed to initialize database", "error", err)
        os.Exit(1)
    }
    log.Info("Database initialized successfully")

    // Create and run the API server
    server := api.NewAPIServer(":3000", store)
    log.Info("Starting server", "address", ":3000")
    server.Run()
}
```

### 2. Update API Router

Modify router.go to add request logging middleware:

```go
package api

import (
    "net/http"

    "socialNetwork/pkg/api/handlers"
    "socialNetwork/pkg/logger"
    "socialNetwork/pkg/storage"
)

// Router sets up all API routes and handlers
func (s *APIServer) setupRouter() *http.ServeMux {
    mux := http.NewServeMux()

    // Initialize handlers
    authHandler := handlers.NewAuthHandler(s.Store)
    chatHandler := handlers.NewChatHandler(s.Store, s.Hub)
    // More handlers will be initialized here...

    // Add logger middleware
    logMiddleware := logger.RequestMiddleware()

    // AUTH ROUTES - with logging middleware
    mux.Handle("POST /auth/register", logMiddleware(makeHTTPHandleFunc(authHandler.HandleRegister)))
    mux.Handle("POST /auth/login", logMiddleware(makeHTTPHandleFunc(authHandler.HandleLogin)))
    mux.Handle("POST /auth/logout", logMiddleware(makeHTTPHandleFunc(authHandler.HandleLogout)))

    // CHAT ROUTES - with logging middleware
    mux.Handle("GET /chats", logMiddleware(makeHTTPHandleFunc(chatHandler.HandleGetChats)))
    mux.Handle("GET /chats/{userId}", logMiddleware(makeHTTPHandleFunc(chatHandler.HandleGetChatHistory)))
    mux.Handle("GET /groups/{id}/chat", logMiddleware(makeHTTPHandleFunc(chatHandler.HandleGetGroupChatHistory)))
    
    // WebSocket handlers (these don't use makeHTTPHandleFunc because they handle their own responses)
    mux.Handle("/ws/chat/{userId}", logMiddleware(http.HandlerFunc(chatHandler.HandleChatWebSocket)))
    mux.Handle("/ws/groups/{id}/chat", logMiddleware(http.HandlerFunc(chatHandler.HandleGroupChatWebSocket)))

    // Additional routes will be added here as we implement more handlers...

    return mux
}
```

### 3. Update Auth Handler Example

Modify auth_handlers.go to use request-specific logging:

```go
package handlers

import (
    "encoding/json"
    "net/http"

    "socialNetwork/pkg/api"
    "socialNetwork/pkg/domain/auth"
    "socialNetwork/pkg/errors"
    "socialNetwork/pkg/logger"
    "socialNetwork/pkg/storage"
)

// AuthHandler manages authentication-related handlers
type AuthHandler struct {
    store storage.Storage
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(store storage.Storage) *AuthHandler {
    return &AuthHandler{
        store: store,
    }
}

// HandleRegister handles user registration
func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) error {
    // Get logger from request context
    log := logger.FromRequest(r)
    
    // Parse the request body
    var userAccount auth.UserAccount
    if err := json.NewDecoder(r.Body).Decode(&userAccount); err != nil {
        log.Error("Failed to parse request body", "error", err)
        return errors.ErrBadRequest
    }

    // Validate the user input
    // This could be more detailed in a production application
    if userAccount.Email == "" || userAccount.Password == "" {
        log.Warn("Invalid input for registration", 
            "email_present", userAccount.Email != "",
            "password_present", userAccount.Password != "")
        return errors.ErrInvalidInput
    }

    log.Info("Creating new user account", "email", userAccount.Email)

    // Create the user account
    userID, err := h.store.CreateUserAccount(&userAccount)
    if err != nil {
        log.Error("Failed to create user account", "error", err)
        return err
    }

    log.Info("User account created successfully", "user_id", userID)

    // Set the ID and return the user account (without password)
    userAccount.ID = userID
    userAccount.Password = ""

    return api.WriteJson(w, http.StatusCreated, userAccount)
}
```

### 4. Update Storage Layer Example

Modify sqlite.go to add logging to database operations:

```go
package sqlite

import (
    "database/sql"
    "log"
    "time"

    _ "github.com/mattn/go-sqlite3"
    "socialNetwork/pkg/domain/auth"
    "socialNetwork/pkg/domain/chat"
    "socialNetwork/pkg/domain/following"
    "socialNetwork/pkg/domain/groups"
    "socialNetwork/pkg/domain/posts"
    "socialNetwork/pkg/domain/profile"
    "socialNetwork/pkg/errors"
    "socialNetwork/pkg/logger"
)

// SQLiteStore implements the Storage interface using SQLite
type SQLiteStore struct {
    db     *sql.DB
    logger logger.Logger
}

// Init verifies the database connection
func (s *SQLiteStore) Init() error {
    s.logger.Info("Initializing database connection")
    return s.db.Ping()
}

// NewSQLiteStore creates a new SQLite storage instance
func NewSQLiteStore() (*SQLiteStore, error) {
    // Get logger for sqlite package
    log := logger.GetLogger("storage/sqlite")
    
    // Using relative path for the database
    connStr := "../../../pkg/db/sqlite/main.db"
    log.Info("Opening database connection", "path", connStr)

    db, err := sql.Open("sqlite3", connStr)
    if err != nil {
        log.Error("Failed to open database", "error", err)
        return nil, err
    }

    // Configure connection pool
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(time.Hour)
    
    log.Debug("Database connection pool configured",
        "max_open", 10,
        "max_idle", 5,
        "max_lifetime", "1h")

    return &SQLiteStore{
        db:     db,
        logger: log,
    }, nil
}
```

### 5. Update WebSocket Example

Modify websocket.go to add logging to WebSocket operations:

```go
package websocket

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"

    "github.com/gorilla/websocket"
    "socialNetwork/pkg/logger"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mutex      sync.Mutex
    logger     logger.Logger
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        logger:     logger.GetLogger("websocket"),
    }
}

// Run starts the Hub's message handling loop
func (h *Hub) Run() {
    h.logger.Info("Starting WebSocket hub")
    for {
        select {
        case client := <-h.register:
            h.mutex.Lock()
            h.clients[client] = true
            h.mutex.Unlock()
            h.logger.Info("Client connected to WebSocket",
                "client_id", client.id,
                "remote_addr", client.conn.RemoteAddr().String())
        case client := <-h.unregister:
            h.mutex.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
                h.logger.Info("Client disconnected from WebSocket",
                    "client_id", client.id,
                    "remote_addr", client.conn.RemoteAddr().String())
            }
            h.mutex.Unlock()
        case message := <-h.broadcast:
            h.mutex.Lock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                    h.logger.Warn("Failed to send message to client, closing connection",
                        "client_id", client.id)
                }
            }
            h.mutex.Unlock()
        }
    }
}
```

## Integration Patterns

### 1. Package-Specific Loggers
For each package, create a logger instance at the package or struct level:
```go
// Package level
var log = logger.GetLogger("package/name")

// Struct level
type MyStruct struct {
    // ...
    logger logger.Logger
}

func NewMyStruct() *MyStruct {
    return &MyStruct{
        // ...
        logger: logger.GetLogger("package/name"),
    }
}
```

### 2. Request Context Logging
For HTTP handlers, get the logger from the request context:
```go
func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) error {
    log := logger.FromRequest(r)
    // Use log throughout the handler
}
```

### 3. Structured Logging Pattern
Add context with structured fields:
```go
// Add single field
log.WithField("user_id", userID).Info("User logged in")

// Add multiple fields
log.WithFields(map[string]interface{}{
    "user_id": userID,
    "ip": clientIP,
}).Info("User logged in")
```

### 4. Performance Tracking Pattern
Track operation durations:
```go
start := time.Now()
// Perform operation
duration := time.Since(start)
log.Info("Operation completed", 
    "operation", "query_users",
    "duration_ms", duration.Milliseconds())
```

## Success Criteria
Logger integration is successful when:
1. All application components use the logger consistently
2. Log messages provide meaningful context for debugging
3. Request tracking works across components
4. Errors are logged with appropriate context
5. Performance insights can be gathered from logs

## Next Steps
After integrating the logger, test the application thoroughly to ensure logging works as expected across all components.