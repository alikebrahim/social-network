# Software Requirements Specification (SRS) for Social Network Backend

## 1. Introduction

### 1.1 Purpose
This Software Requirements Specification (SRS) outlines the functional and non-functional requirements for the backend of a social network application, designed as an educational project. The backend handles server-side logic, database interactions, and real-time communication to support features like user authentication, social connections, posts, groups, notifications, and chats.

### 1.2 Scope
The backend provides a RESTful API and WebSocket-based real-time communication. Key features include:
- User authentication with sessions and cookies
- Follower system with public/private profile support
- User profiles with privacy settings
- Post creation with comments, likes, and granular privacy levels
- Group management with membership, events, and internal posts
- Real-time private and group chats with emoji support
- Notification system with real-time delivery
- Containerized deployment (Docker planned)

Backend is written in Go using SQLite and follows a domain-driven architecture with modular packages and a clear separation between server, app logic, and storage.

### 1.3 Definitions
(unchanged)

### 1.4 References
- `README.md`, `AUDIT.md`, `API Endpoints Doc`, backend codebase

## 2. Overall Description

### 2.1 User Needs
Users expect:
- Secure login with persistent sessions
- Privacy-aware profile and post visibility
- Engagement via follow, comment, like
- Private and group chat
- Real-time notifications

### 2.2 Assumptions and Dependencies
- Accessed via browser-based frontend
- SQLite sufficient for educational scope
- Docker for containerization (to be implemented)

## 3. System Features

### 3.1 Authentication
- ✔ Email uniqueness checked explicitly
- ✔ bcrypt used for hashing
- ✔ Sessions stored with UUID tokens
- ✔ Login, logout, session management via cookies

### 3.2 Followers
- ✔ Public: auto-follow
- ✔ Private: request/accept workflow
- ✔ Notification on follow request implemented

### 3.3 Profile
- ✔ User info viewable by others with visibility restrictions
- ✔ Users can toggle privacy
- ✔ Posts per profile viewable via `/profiles/{id}/posts`

### 3.4 Posts
- ✔ Posts support JPEG/PNG/GIF validation
- ✔ `privacy_level`: public, almost_private, private
- ✔ `allowed_followers` field implemented
- ✔ Edit/delete supported

### 3.5 Groups
- ✔ Create, invite, request join
- ✔ Internal posts, events
- ✔ Notifications for invites, requests, and events

### 3.6 Chat
- ✔ WebSocket-based messaging for direct and group chat
- ✔ Validates follower/public profile status
- ✔ Emoji support handled

### 3.7 Notifications
- ✔ Created for all expected triggers
- ✔ Stored, fetchable, markable as read
- ✔ Delivered via WebSocket

### 3.8 Database and Migrations
- ✔ SQLite used
- ✔ Goose migration system applied on startup via `Init()`

### 3.9 Docker
- ❌ No Dockerfile yet
- ❌ No container orchestration implemented

## 4. External Interface Requirements

### 4.1 User Interfaces
- API only, consumed by frontend

### 4.2 Hardware Interfaces
- Server with file access for SQLite

### 4.3 Software Interfaces
- Go 1.22+, SQLite, Gorilla WebSocket, bcrypt, uuid

### 4.4 Communication Interfaces
- HTTP/REST (port 3000)
- WebSocket (port 3000)

## 5. Non-Functional Requirements

### 5.1 Performance
- API ≤ 500ms, WebSocket ≤ 100ms

### 5.2 Security
- bcrypt, secure UUID sessions, input validation

### 5.3 Quality Attributes
- Maintainable, modular, reliable

### 5.4 Compliance
- Meets allowed packages, OWASP session guidelines

## 6. Evaluation Against AUDIT.md
- ✅ Backend structure, session management, follower logic, and group chat validated
- ✅ Notifications and profile visibility implemented
- ❌ Docker containers missing

## 7. Recommendations
- [x] Dockerfile + container setup
- [x] Optionally normalize `allowed_followers`

## 8. Conclusion
The backend meets all major project criteria except Docker. The architecture is robust, real-time features are implemented, and all SRS gaps previously noted have been addressed in the codebase.

