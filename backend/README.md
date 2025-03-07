# Social Network Backend

This is a Go-based backend for a social network application. It provides REST API endpoints and WebSocket functionality for user interaction, group management, and real-time chat.

## Features

- User authentication and profile management
- Social connections (following/followers)
- Groups with membership management
- Group posts and comments
- Events with RSVP functionality
- Real-time chat (private and group)
- WebSocket support for instant messaging

## Technology Stack

- Go 1.22+
- SQLite for database
- Native net/http for API endpoints
- Gorilla WebSocket for real-time communication

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

### Running the Server

```
go run .
```

Server will start on port 3000 by default.

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

## API Endpoints

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

## License

This project is licensed under the MIT License - see the LICENSE file for details.