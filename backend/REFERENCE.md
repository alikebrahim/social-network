# Social Network Backend: Comprehensive Reference

This document serves as a complete reference for the Social Network Backend project, consolidating information from all phase documents and implementation details. It provides a single source of truth for understanding the project's structure, implementation, and development history.

## Table of Contents

1. [Project Overview](#project-overview)
2. [Project History](#project-history)
3. [Directory Structure](#directory-structure)
4. [Architectural Decisions](#architectural-decisions)
5. [Package Migration](#package-migration)
6. [Logger Implementation](#logger-implementation)
7. [Integration Summary](#integration-summary)
8. [API Endpoints](#api-endpoints)
9. [Future Considerations](#future-considerations)

## Project Overview

The Social Network Backend is a Go-based server application that provides REST API endpoints and WebSocket functionality for a social networking platform. It supports user authentication, profile management, social connections, post creation, groups, events, and real-time chat.

### Core Features

- Authentication and user management
- Profile creation and privacy controls
- Following system with request acceptance workflow
- Post creation with comments and likes
- Groups with membership management
- Events with RSVP functionality
- Real-time chat via WebSockets
- Structured logging with color-coded levels

### Technology Stack

- Go 1.22+
- SQLite for database storage
- Native net/http for API handling
- Gorilla WebSocket for real-time communication
- Custom structured logger

## Project History

### Initial Structure

The project originally used an `/internal` directory structure:

```
/backend
├── cmd/server                # Entry point
├── internal                  # Internal packages
│   ├── api                   # API server and handlers
│   ├── domain                # Domain models
│   ├── storage               # Storage interface and implementation
│   └── websocket             # WebSocket implementation
├── pkg                       # Public packages
│   └── errors                # Error definitions
└── test_api.sh               # API test script
```

### Restructuring Plan

A comprehensive restructuring plan was developed to:
1. Migrate all code from `/internal` to `/pkg` for better organization
2. Implement a robust logging system
3. Resolve import cycle issues
4. Enhance the overall architecture

The plan was executed in seven phases, each with specific goals and deliverables.

## Directory Structure

### Current Structure

The backend now follows a domain-driven architecture:

```
/backend
├── cmd/server                # Entry point for the application
├── pkg                       # Application packages
│   ├── api                   # API server and routes
│   │   ├── handlers          # HTTP handlers for API endpoints
│   │   └── utils             # HTTP utilities
│   ├── db                    # Database migrations and files
│   │   ├── migrations        # SQL migration scripts
│   │   └── sqlite            # SQLite database
│   ├── domain                # Domain models
│   │   ├── auth              # Authentication domain
│   │   ├── chat              # Chat domain
│   │   ├── following         # Following/followers domain
│   │   ├── groups            # Groups domain
│   │   ├── posts             # Posts domain
│   │   └── profile           # User profiles domain
│   ├── errors                # Error definitions
│   ├── logger                # Structured logging system
│   ├── storage               # Data storage interfaces
│   │   └── sqlite            # SQLite implementation
│   └── websocket             # WebSocket implementation
└── test_api.sh               # API test script
```

### Key Changes

- Removed the `/internal` directory completely
- Created a more modular structure in `/pkg`
- Added a `/pkg/logger` package for structured logging
- Created a `/pkg/api/utils` package to resolve import cycles
- Ensured proper separation of concerns across all packages

## Architectural Decisions

### Domain-Driven Design

The application follows domain-driven design principles:

- Clear separation between domain models and storage implementations
- Domain packages that encapsulate business logic
- Storage interfaces that allow for different database implementations
- Well-defined API handlers for each domain area

### Import Cycle Resolution

A critical architectural issue was resolved during the restructuring:

- Identified an import cycle between `internal/api` and `internal/api/handlers`
- Created a new `pkg/api/utils` package to extract shared functionality
- Moved HTTP utilities to break the circular dependency
- Ensured proper directional flow of dependencies

### Authentication System

The application uses a session-based authentication system:

- Session tokens stored securely in cookies
- JWT middleware for protected routes
- Privacy controls based on user preferences
- Authorization checks for all protected operations

## Package Migration

The migration from `/internal` to `/pkg` was executed methodically:

### Assessment Phase (PHASE-1)

- Analyzed directory structure and dependencies
- Created dependency graph to guide migration order
- Identified potential conflicts and risks
- Evaluated database path references

### Pre-Migration Checks (PHASE-2)

- Created branch backup for safety
- Established verification criteria
- Tested current application functionality
- Discovered import cycle issue

### Migration Execution (PHASE-3)

- Created new directory structure
- Migrated files in dependency order:
  1. Domain models
  2. Storage interface
  3. SQLite implementation
  4. WebSocket
  5. API components
  6. Main application
- Updated import statements
- Resolved database path references
- Created `utils` package to break import cycles

### Post-Migration Testing (PHASE-4)

- Verified application functionality
- Confirmed proper handling of database operations
- Tested all API endpoints
- Validated WebSocket connections
- Fixed any migration-related issues

## Logger Implementation

A comprehensive logging system was implemented in three phases:

### Logger Design (PHASE-5)

- Defined logger interface with key methods:
  ```go
  type Logger interface {
      Debug(msg string, keyvals ...interface{})
      Info(msg string, keyvals ...interface{})
      Warn(msg string, keyvals ...interface{})
      Error(msg string, keyvals ...interface{})
      WithField(key string, value interface{}) Logger
      WithFields(fields map[string]interface{}) Logger
      WithPackage(pkg string) Logger
      Close() error
  }
  ```
- Established log levels: DEBUG, INFO, WARN, ERROR
- Designed configuration options for console/file output
- Planned color-coding for improved readability

### Logger Implementation (PHASE-6)

- Created standard logger implementation
- Added structured logging with key-value pairs
- Implemented file and console output options
- Added package-specific logger instances
- Implemented context utilities for request logging

### Logger Integration (PHASE-7)

- Added logger initialization to main.go
- Integrated with HTTP middleware for request logging
- Added logging to critical application operations
- Implemented request context logging
- Added performance tracking for database operations

### Color Enhancement

The logger was enhanced with color-coded output:

- DEBUG level: Cyan
- INFO level: Green
- WARN level: Yellow
- ERROR level: Red

This improves readability in console output while maintaining plain text in log files.

### Configuration Options

The logger supports configuration through environment variables:

| Variable | Description | Default | Options |
|----------|-------------|---------|---------|
| `LOG_LEVEL` | Sets the minimum log level to display | `INFO` | `DEBUG`, `INFO`, `WARN`, `ERROR` |
| `LOG_CONSOLE` | Enable/disable console output | `true` | `true`, `false` |
| `LOG_COLORS` | Enable/disable color-coded log levels | `true` | `true`, `false` |

## Integration Summary

After the package migration, additional functionality from a separate branch (`backend`) was integrated into the restructured code.

### Feature Integration Status

| Feature Area | Implementation Status | Route Registration | Status |
|--------------|----------------------|-------------------|--------|
| Authentication | Complete | Complete | ✅ Unchanged |
| Profile | Complete | Complete | ✅ Integrated |
| Following | Complete | Complete | ✅ Integrated |
| Posts | Complete | Complete | ✅ Integrated |
| Comments | Complete | Complete | ✅ Integrated |
| Groups | Complete | Complete | ✅ Integrated |
| Events | Complete | Complete | ✅ Integrated |
| Chat | Complete | Complete | ✅ Unchanged |
| WebSockets | Complete | Complete | ✅ Unchanged |

### Integration Details

- Router Integration:
  - Comprehensive router in `pkg/api/router.go` with all endpoints
  - Authentication middleware for protected routes
  - Consistent route patterns across handlers

- Handler Integration:
  - Created multiple handler implementations:
    - `following_handlers.go`
    - `posts_handlers.go`
    - `groups_handlers.go`
    - `profile_handlers.go`

- Storage Integration:
  - Enhanced storage implementations for all functionality
  - Created proper database queries for all operations
  - Maintained interface compatibility

## API Endpoints

The application provides a comprehensive set of API endpoints:

### Auth
- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login with credentials
- `POST /auth/logout` - Logout user

### Profile
- `GET /profiles/{id}` - Get user profile
- `PUT /profiles/privacy` - Update profile privacy
- `GET /profiles/{id}/activity` - Get user activity

### Following
- `POST /follow/{id}` - Follow a user
- `GET /follow/requests` - Get follow requests
- `POST /follow/{id}/accept` - Accept follow request
- `DELETE /follow/{id}` - Unfollow or reject request

### Posts
- `POST /posts` - Create a post
- `GET /posts/{id}` - Get post details
- `PUT /posts/{id}` - Edit a post
- `DELETE /posts/{id}` - Delete a post
- `POST /posts/{id}/comments` - Comment on a post
- `POST /posts/{id}/likes` - Like a post
- `DELETE /posts/{id}/likes` - Unlike a post

### Groups
- `POST /groups` - Create a group
- `GET /groups` - List all groups
- `GET /groups/{id}` - Get group details
- `POST /groups/{id}/invite` - Invite user to group
- `GET /groups/search?q={query}` - Search for groups
- `POST /groups/{id}/join` - Request to join group
- `GET /groups/{id}/requests` - List join requests
- `POST /groups/{id}/requests/{userId}/accept` - Accept join request
- `POST /groups/{id}/requests/{userId}/reject` - Reject join request
- `GET /groups/{id}/members` - List group members

### Events
- `POST /groups/{id}/events` - Create event
- `GET /groups/{id}/events` - List group events
- `POST /groups/{id}/events/{eventId}/response` - Respond to event

### Chat
- `GET /chats` - List user's chats
- `GET /chats/{userId}` - Get chat history with specific user
- `GET /groups/{id}/chat` - Get group chat history
- `WebSocket /ws/chat/{userId}` - Real-time chat with user
- `WebSocket /ws/groups/{id}/chat` - Real-time group chat

## Future Considerations

While the application is now fully functional with enhanced architecture and logging, several areas have been identified for future improvement:

### Logging Enhancements
- Implement log rotation for better file management
- Add log compression for archival purposes
- Enhance log filtering capabilities
- Add per-domain logging levels

### API Improvements
- Add OpenAPI/Swagger documentation for API endpoints
- Implement request validation middleware
- Add rate limiting for security
- Improve error responses

### Performance Optimizations
- Implement Redis or other caching mechanisms
- Optimize database queries with indexes
- Add pagination for large result sets
- Implement connection pooling

### Feature Extensions
- Enhance search capabilities across content types
- Improve media file handling
- Add notification system
- Implement analytics tracking

### Testing Enhancements
- Add comprehensive unit tests
- Implement integration test suite
- Add performance benchmarks
- Create load testing scripts

## Conclusion

The Social Network Backend has undergone significant improvements through the package migration and logger implementation. The application now features a clean architecture, robust logging, and comprehensive feature set, providing a solid foundation for future development.

The codebase follows best practices for Go development, with clear separation of concerns, proper error handling, and consistent patterns across all components. The logger provides valuable insights into application behavior, aiding in debugging and monitoring.

By consolidating the information from all phase documents, this reference serves as the definitive guide to the project's structure, implementation, and development history.