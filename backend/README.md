# Social Network Backend

This is a Go-based backend for a social network application. It provides REST API endpoints and WebSocket functionality for user interaction, group management, and real-time chat.

## Features

- User authentication and profile management
- Social connections (following/followers)
- Posts with privacy controls and comments
- Groups with membership management
- Events with RSVP functionality
- Real-time chat (private and group)
- WebSocket support for instant messaging
- Structured logging with color-coded levels

## Technology Stack

- Go 1.22+
- SQLite for database
- Native net/http for API endpoints
- Gorilla WebSocket for real-time communication
- Custom structured logger

## Project Structure

The backend follows a domain-driven architecture:

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

## Getting Started

### Prerequisites

- Go 1.22+
- SQLite

### Installation

1. Clone the repository
2. Navigate to the project directory
3. Install dependencies:
```
go mod download
```
4. Run database migrations:
```
goose -dir pkg/db/migrations/sqlite sqlite3 pkg/db/sqlite/main.db up
```

### Running the Server

Basic usage:
```
go run ./cmd/server
```

With custom logging configuration:
```
LOG_LEVEL=DEBUG LOG_COLORS=true go run ./cmd/server
```

Server will start on port 3000 by default.

## Logger Configuration

The application features a robust structured logging system with the following capabilities:

### Log Levels

Four severity levels are available:
- **DEBUG** (cyan): Detailed information for debugging purposes
- **INFO** (green): General information about application operation
- **WARN** (yellow): Potential issues that don't prevent operation
- **ERROR** (red): Critical errors that may impair functionality

### Environment Variables

Configure the logger using these environment variables:

| Variable | Description | Default | Options |
|----------|-------------|---------|---------|
| `LOG_LEVEL` | Sets the minimum log level to display | `INFO` | `DEBUG`, `INFO`, `WARN`, `ERROR` |
| `LOG_CONSOLE` | Enable/disable console output | `true` | `true`, `false` |
| `LOG_COLORS` | Enable/disable color-coded log levels | `true` | `true`, `false` |

### Log Output

- Console output: Enabled by default with color-coded levels
- File output: Logs are written to `./logs/app.log`

### Structured Logging

The logger supports structured key-value pairs for enhanced context:

```go
logger.WithField("user_id", userID).Info("User logged in")
logger.WithFields(map[string]interface{}{
    "user_id": userID,
    "ip": clientIP,
}).Info("User authenticated")
```

### Request Logging

All HTTP requests are automatically logged with:
- Request ID for tracking
- HTTP method and path
- Client IP and user agent
- Response status code
- Request duration

## API Endpoints

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

## API Testing

You can use the included test tools to verify the functionality:

1. Test REST API endpoints:
```
chmod +x test_api.sh
./test_api.sh
```

2. Test WebSocket functionality:
- Open `websocket_test.html` in a browser
- Enter your user ID, target user/group ID, and connect
- Send and receive messages in real-time

## Project History

The project underwent a significant restructuring to improve organization and maintainability:

### Package Migration

The codebase was migrated from an `/internal` to a `/pkg` structure to follow Go best practices. Key improvements included:

- Better organization of domain-driven components
- Resolution of import cycles using a utils package
- Centralized error handling and HTTP utilities
- Enhanced package visibility and reusability

### Logging Implementation

A comprehensive structured logging system was implemented with:

- Multiple severity levels with color coding
- Context-aware request logging
- Performance tracking for database operations
- File and console output options
- Package-specific loggers for targeted debugging

## Future Enhancements

Potential improvements for future consideration:

1. Enhanced logging features like log rotation
2. OpenAPI/Swagger documentation for the API
3. Redis or other caching mechanisms for frequently accessed data
4. Pagination support for endpoints returning large result sets
5. Advanced search capabilities across posts, groups, and profiles
6. Enhanced media file handling and storage

## License

This project is licensed under the MIT License.