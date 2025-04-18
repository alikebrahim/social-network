# Software Requirements Specification (SRS) for Social Network Backend

## 1. Introduction

### 1.1 Purpose
This Software Requirements Specification (SRS) outlines the functional and non-functional requirements for the backend of a social network application, designed as an educational project. The backend is responsible for handling server-side logic, database interactions, and real-time communication to support features like user authentication, social connections, posts, groups, notifications, and chats. This document evaluates the current implementation against the project expectations outlined in the provided README.md and AUDIT.md, ensuring all required functionalities are addressed.

### 1.2 Scope
The social network backend provides a RESTful API and WebSocket-based real-time communication to support a client-side frontend. Key features include:
- User authentication with sessions and cookies
- Follower system with public/private profile support
- User profiles with privacy settings
- Post creation with comments, likes, and privacy controls
- Group management with membership and events
- Real-time private and group chats
- Notification system for user interactions
- Containerized deployment using Docker

The backend is implemented in Go with SQLite as the database, adhering to the allowed packages specified in the project requirements.

### 1.3 Definitions, Acronyms, and Abbreviations
- **API**: Application Programming Interface
- **REST**: Representational State Transfer
- **WebSocket**: Protocol for real-time, bidirectional communication
- **SQLite**: Lightweight, file-based relational database
- **Docker**: Platform for containerizing applications
- **SRS**: Software Requirements Specification
- **UUID**: Universally Unique Identifier
- **Bcrypt**: Password hashing algorithm
- **Middleware**: Software layer that processes requests/responses

### 1.4 References
- **README.md**: Project objectives and instructions
- **AUDIT.md**: Evaluation criteria for the implementation
- **backend/ directory**: Current backend implementation files

## 2. Overall Description

### 2.1 User Needs
The backend serves users who interact with the social network via a frontend interface. Users expect:
- Secure registration and login with persistent sessions
- Ability to follow/unfollow users with request/accept mechanics for private profiles
- Customizable profiles with privacy settings
- Creation and interaction with posts (comments, likes) with privacy controls
- Group creation, membership management, and event organization
- Real-time private and group chats with emoji support
- Notifications for follow requests, group invitations, join requests, and events
- Reliable performance and data persistence

### 2.2 Assumptions and Dependencies
- **Assumptions**:
  - Users access the application through a compatible web browser.
  - The frontend is responsible for rendering UI and making API calls.
  - SQLite is sufficient for the educational scope; production would require a more robust database.
- **Dependencies**:
  - Go 1.22+ for backend development
  - SQLite for data storage
  - External packages: gorilla/websocket, mattn/go-sqlite3, golang.org/x/crypto/bcrypt, google/uuid
  - Docker for containerization
  - Frontend application to consume backend APIs

## 3. System Features

### 3.1 Authentication
#### 3.1.1 Description and Priority
Users must register and log in to access the social network. Sessions and cookies maintain user authentication state. Logout functionality is available. Priority: High.

#### 3.1.2 Functional Requirements
- **FR1.1**: Register a user with email, password, first name, last name, date of birth, and optional fields (avatar, nickname, about me).
- **FR1.2**: Validate unique email during registration.
- **FR1.3**: Hash passwords using bcrypt before storage.
- **FR1.4**: Log in with email and password, generating a session token stored in a cookie.
- **FR1.5**: Reject login attempts with incorrect credentials.
- **FR1.6**: Maintain session state across browser refreshes for logged-in users.
- **FR1.7**: Allow logout, invalidating the session.
- **FR1.8**: Ensure separate sessions for different users in different browsers.

#### 3.1.3 Current Implementation Evaluation
- **Implemented**:
  - Registration (`pkg/storage/sqlite/auth.go:CreateUserAccount`) supports all required fields and uses bcrypt for password hashing.
  - Login (`AuthenticateUser`) generates a session token stored in the sessions table.
  - Logout (`DeleteSession`) removes the session.
  - Session validation (`GetUserIdBySession`) checks for expiry.
  - Unique email validation is implied via database constraints.
- **Gaps**:
  - No explicit validation for duplicate emails in `CreateUserAccount` before insertion; relies on database error handling.
  - Session token generation (`generateSessionToken`) uses timestamp-based strings, which is insecure; should use a secure random generator (e.g., google/uuid).
- **Recommendations**:
  - Add explicit email uniqueness check before insertion.
  - Use google/uuid for session tokens to enhance security.

#### 3.1.4 Non-Functional Requirements
- **NFR1.1**: Registration and login must complete within 1 second under normal load.
- **NFR1.2**: Password hashing must use bcrypt with a secure cost factor.
- **NFR1.3**: Session tokens must be cryptographically secure.

### 3.2 Followers
#### 3.2.1 Description and Priority
Users can follow others, with private profiles requiring a request/accept process. Public profiles allow immediate following. Priority: High.

#### 3.2.2 Functional Requirements
- **FR2.1**: Send a follow request to a private user, pending acceptance.
- **FR2.2**: Automatically follow a public user without a request.
- **FR2.3**: Accept or decline follow requests for private profiles.
- **FR2.4**: Unfollow a user if already following.
- **FR2.5**: List pending follow requests for a user.

#### 3.2.3 Current Implementation Evaluation
- **Implemented**:
  - Follow request creation (`pkg/storage/sqlite/following.go:CreateFollowRequest`) checks profile type and sets status (pending for private, accepted for public).
  - Accept (`AcceptFollowRequest`) and delete (`DeleteFollowRequest`) requests.
  - Retrieve pending requests (`GetFollowRequests`).
  - Unfollow via `DeleteFollowRequest`.
- **Gaps**:
  - No notification generated on follow request creation (required per Notifications section).
- **Recommendations**:
  - Integrate notification creation in `CreateFollowRequest` for private profiles.

#### 3.2.4 Non-Functional Requirements
- **NFR2.1**: Follow operations must complete within 500ms.
- **NFR2.2**: Follower lists must be retrieved efficiently using indexed queries.

### 3.3 Profile
#### 3.3.1 Description and Priority
Users have profiles displaying registration information (except password), posts, followers, and following users. Profiles can be public or private. Priority: High.

#### 3.3.2 Functional Requirements
- **FR3.1**: Display user information, posts, followers, and following users.
- **FR3.2**: Restrict private profile visibility to accepted followers.
- **FR3.3**: Allow public profile visibility to all users.
- **FR3.4**: Enable users to toggle profile privacy (public/private).
- **FR3.5**: Show all profile data to the profile owner.

#### 3.3.3 Current Implementation Evaluation
- **Implemented**:
  - Profile data retrieval (`pkg/storage/sqlite/profile.go:GetProfileData`) includes all required fields, follower/following counts, and post count.
  - Privacy toggle (`SetProfilePrivacy`) updates profile type.
  - Visibility logic implied in post retrieval (`CanUserSeePost`), which checks profile type and follower status.
- **Gaps**:
  - No explicit profile visibility check in `GetProfileData`; assumes frontend handles visibility.
  - No endpoint to retrieve user posts directly in profile context (handled in posts).
- **Recommendations**:
  - Add visibility check in `GetProfileData` to enforce private profile restrictions.
  - Implement a profile-specific endpoint to fetch user posts.

#### 3.3.4 Non-Functional Requirements
- **NFR3.1**: Profile data retrieval must complete within 500ms.
- **NFR3.2**: Privacy changes must be reflected immediately.

### 3.4 Posts
#### 3.4.1 Description and Priority
Users can create posts and comments with optional images (JPEG, PNG, GIF) and specify privacy (public, almost private, private). Priority: High.

#### 3.4.2 Functional Requirements
- **FR4.1**: Create posts with content, optional image, and privacy setting (public, almost private, private).
- **FR4.2**: Create comments on posts with content and optional image.
- **FR4.3**: Like/unlike posts.
- **FR4.4**: Edit/delete own posts.
- **FR4.5**: Restrict post visibility based on privacy:
  - Public: All users
  - Almost private: Followers only
  - Private: Selected followers
- **FR4.6**: Retrieve post details, including comments and like count.

#### 3.4.3 Current Implementation Evaluation
- **Implemented**:
  - Post creation (`pkg/storage/sqlite/posts.go:CreatePost`) supports content and image.
  - Comment creation (`CreateComment`), like/unlike (`CreateLike`, `RemoveLikes`).
  - Edit/delete posts (`EditPost`, `DeletePost`).
  - Post retrieval (`GetPostByID`) includes comments and likes.
  - Visibility check (`CanUserSeePost`) enforces public/private profile rules.
- **Gaps**:
  - No support for "almost private" or "private" (selected followers) privacy levels; only public/private profile-based visibility.
  - Image handling supports storage but does not validate JPEG, PNG, GIF formats explicitly.
- **Recommendations**:
  - Add privacy field to posts table and implement logic for almost private/private settings.
  - Validate image formats in `CreatePost` and `CreateComment`.

#### 3.4.4 Non-Functional Requirements
- **NFR4.1**: Post operations must complete within 500ms.
- **NFR4.2**: Image uploads must support files up to 5MB.

### 3.5 Groups
#### 3.5.1 Description and Priority
Users can create groups, invite members, manage join requests, create posts, and organize events. Priority: High.

#### 3.5.2 Functional Requirements
- **FR5.1**: Create groups with title and description.
- **FR5.2**: Invite users to groups; invited users accept/reject invitations.
- **FR5.3**: Allow users to request to join groups; creators accept/reject requests.
- **FR5.4**: List all groups and search by title/description.
- **FR5.5**: Create posts within groups, visible only to members.
- **FR5.6**: Create events with title, description, date/time, and options (Going, Not Going).
- **FR5.7**: Allow members to respond to events.

#### 3.5.3 Current Implementation Evaluation
- **Implemented**:
  - Group creation (`pkg/storage/sqlite/groups.go:CreateGroup`).
  - Invite/accept/reject (`InviteToGroup`, `AcceptGroupInvite`, `RejectGroupInvite`).
  - Join requests (`RequestJoinGroup`, `AcceptGroupJoinRequest`, `RejectGroupJoinRequest`).
  - Group listing and search (`GetGroups`, `SearchGroups`).
  - Event creation and responses (`CreateGroupEvent`, `RespondToEvent`, `GetEventResponses`).
- **Gaps**:
  - Group post creation (`CreateGroupPost`, `GetGroupPosts`) not implemented.
  - No notification for group invitations or join requests.
- **Recommendations**:
  - Implement group post functionality.
  - Add notification generation for invitations and join requests.

#### 3.5.4 Non-Functional Requirements
- **NFR5.1**: Group operations must complete within 500ms.
- **NFR5.2**: Group search must handle up to 1000 groups efficiently.

### 3.6 Chat
#### 3.6.1 Description and Priority
Users can send private messages to followers or public-profile users and participate in group chats. Priority: High.

#### 3.6.2 Functional Requirements
- **FR6.1**: Send/receive private messages in real-time via WebSocket if users follow each other or the recipient has a public profile.
- **FR6.2**: Prevent chats between non-following private-profile users.
- **FR6.3**: Support group chats for members, with real-time messaging.
- **FR6.4**: Support emoji in messages.
- **FR6.5**: Retrieve chat history for private and group chats.

#### 3.6.3 Current Implementation Evaluation
- **Implemented**:
  - WebSocket hub (`pkg/websocket/websocket.go`) manages client connections.
  - Private chat (`SaveChat`, `GetChatHistory`) and group chat (`SaveGroupChat`, `GetGroupChatHistory`) storage.
  - Real-time messaging via `SendMessageToUser` and `SendMessageToGroup`.
  - Chat history retrieval (`GetUserChats`, `GetChatHistory`).
- **Gaps**:
  - No explicit check for follower/public profile status before allowing private chats.
  - Emoji support not explicitly validated (assumed to work via JSON).
  - Group chat messages lack validation for group membership (handled in `SaveGroupChat` but not in WebSocket).
- **Recommendations**:
  - Add follower/public profile check in WebSocket connection setup.
  - Validate emoji support in message processing.
  - Enforce group membership in WebSocket message sending.

#### 3.6.4 Non-Functional Requirements
- **NFR6.1**: Messages must be delivered within 100ms.
- **NFR6.2**: WebSocket connections must support up to 100 concurrent users.

### 3.7 Notifications
#### 3.7.1 Description and Priority
Users receive notifications for follow requests, group invitations, join requests, and event creation. Priority: Medium.

#### 3.7.2 Functional Requirements
- **FR7.1**: Notify users of follow requests to their private profile.
- **FR7.2**: Notify users of group invitations.
- **FR7.3**: Notify group creators of join requests.
- **FR7.4**: Notify group members of new events.
- **FR7.5**: Display notifications distinctly from messages.

#### 3.7.3 Current Implementation Evaluation
- **Implemented**:
  - Notifications table exists (`pkg/db/migrations/sqlite/0011_create_notifications_table.sql`).
- **Gaps**:
  - No implementation for notification creation or retrieval in storage or handlers.
  - No API endpoints for notifications.
- **Recommendations**:
  - Implement notification creation in relevant handlers (follow, group, event).
  - Add endpoints to retrieve and mark notifications as read.

#### 3.7.4 Non-Functional Requirements
- **NFR7.1**: Notifications must be generated within 200ms of the triggering event.
- **NFR7.2**: Notification retrieval must complete within 500ms.

### 3.8 Database and Migrations
#### 3.8.1 Description and Priority
The backend uses SQLite with a migration system to manage schema changes. Priority: High.

#### 3.8.2 Functional Requirements
- **FR8.1**: Use SQLite as the database.
- **FR8.2**: Implement migrations to create tables for users, sessions, followers, posts, comments, likes, groups, group members, events, event responses, chats, group chats, and notifications.
- **FR8.3**: Apply migrations automatically on application start.
- **FR8.4**: Organize migrations in a dedicated folder.

#### 3.8.3 Current Implementation Evaluation
- **Implemented**:
  - SQLite database (`pkg/db/sqlite/main.db`).
  - Migration files (`pkg/db/migrations/sqlite/`) cover all required tables.
  - Migration application via goose (per README.md instructions).
  - Organized migration folder structure.
- **Gaps**:
  - No explicit migration application code in `pkg/storage/sqlite/sqlite.go`.
- **Recommendations**:
  - Add migration application logic in `NewSQLiteStore` using golang-migrate or goose.

#### 3.8.4 Non-Functional Requirements
- **NFR8.1**: Database queries must complete within 100ms under normal load.
- **NFR8.2**: Migrations must apply within 5 seconds on startup.

### 3.9 Docker
#### 3.9.1 Description and Priority
The backend must be containerized using Docker, exposing necessary ports. Priority: High.

#### 3.9.2 Functional Requirements
- **FR9.1**: Create a Docker image for the backend.
- **FR9.2**: Expose port 3000 for API and WebSocket communication.
- **FR9.3**: Ensure the container runs the server and applies migrations.

#### 3.9.3 Current Implementation Evaluation
- **Implemented**:
  - README.md implies Docker support (runs on port 3000).
- **Gaps**:
  - No Dockerfile or Docker-related files provided.
  - No evidence of container configuration or port exposure.
- **Recommendations**:
  - Create a Dockerfile to build and run the Go application.
  - Expose port 3000 and ensure migrations are applied on container start.

#### 3.9.4 Non-Functional Requirements
- **NFR9.1**: Docker container must start within 10 seconds.
- **NFR9.2**: Container size must be optimized (under 500MB).

## 4. External Interface Requirements

### 4.1 User Interfaces
- The backend has no direct user interface; it exposes APIs and WebSocket endpoints consumed by the frontend.

### 4.2 Hardware Interfaces
- Runs on standard server hardware or Docker containers.
- Requires storage for SQLite database file.

### 4.3 Software Interfaces
- **Go 1.22+**: For application runtime.
- **SQLite**: For data persistence.
- **External Packages**:
  - github.com/gorilla/websocket
  - github.com/mattn/go-sqlite3
  - golang.org/x/crypto/bcrypt
  - github.com/google/uuid
- **Docker**: For containerization.

### 4.4 Communications Interfaces
- **HTTP/REST**: For API endpoints (port 3000).
- **WebSocket**: For real-time chat (port 3000).
- **SQLite**: File-based database access.

## 5. Non-Functional Requirements

### 5.1 Performance Requirements
- API responses must complete within 500ms under normal load.
- WebSocket messages must deliver within 100ms.
- Database queries must complete within 100ms.
- Application startup (including migrations) must complete within 10 seconds.

### 5.2 Security Requirements
- Passwords must be hashed with bcrypt.
- Session tokens must be cryptographically secure.
- Private data (profiles, posts) must be restricted based on privacy settings.
- Input validation to prevent SQL injection and XSS.

### 5.3 Quality Attributes
- **Maintainability**: Code must follow Go best practices and be well-documented.
- **Scalability**: Handle up to 100 concurrent users for educational purposes.
- **Reliability**: Ensure database transactions are atomic and consistent.

### 5.4 Compliance Requirements
- Adhere to allowed packages specified in README.md.
- Follow OWASP guidelines for session management and cookies.

## 6. Evaluation Against AUDIT.md

### 6.1 Functional Checks
- **Allowed Packages**: Uses only permitted packages (gorilla/websocket, mattn/go-sqlite3, golang.org/x/crypto, google/uuid).
- **File System Organization**: Well-organized with clear separation of packages and migrations (`pkg/db/migrations/sqlite`).

### 6.2 Backend Checks
- **Separation of Responsibilities**: Clear division with server (`cmd/server/main.go`), app (`pkg/api`), and database (`pkg/storage/sqlite`).
- **Server**: Receives requests via net/http (`pkg/api/server.go`).
- **App**: Handles requests and database interactions via handlers and storage.
- **Core Logic**: Implemented in handlers and storage packages.

### 6.3 Database Checks
- **SQLite Usage**: Confirmed with `pkg/db/sqlite/main.db`.
- **Client Interactions**: Storage methods support data retrieval and submission.
- **Migration System**: Migration files exist; application logic needed.
- **Migration Organization**: Matches example structure.

### 6.4 Authentication Checks
- **Sessions**: Implemented with session tokens and cookies.
- **Form Elements**: All required fields supported.
- **Registration/Login**: Functional with error handling.
- **Duplicate Users**: Handled implicitly via database constraints.
- **Browser Sessions**: Separate sessions maintained (assumed via unique tokens).

### 6.5 Followers Checks
- **Private User Follow**: Request-based system implemented.
- **Public User Follow**: Automatic following for public profiles.
- **Accept/Decline**: Supported.
- **Unfollow**: Supported.

### 6.6 Profile Checks
- **Information Display**: All fields except password supported.
- **Posts/Followers**: Counts implemented; post retrieval needs endpoint.
- **Privacy Toggle**: Supported.
- **Visibility**: Partially implemented (needs profile-level check).

### 6.7 Posts Checks
- **Create/Comment**: Supported with image handling.
- **Privacy**: Missing almost private/private options.
- **Image Formats**: Not explicitly validated.

### 6.8 Groups Checks
- **Create/Invite**: Supported.
- **Join Requests**: Supported.
- **Posts**: Not implemented.
- **Events**: Fully supported.

### 6.9 Chat Checks
- **Private/Group Chat**: Real-time via WebSocket.
- **Follower Check**: Missing.
- **Emojis**: Assumed supported.
- **Stability**: No crash indicators.

### 6.10 Notifications Checks
- **Visibility**: Not implemented.
- **Triggers**: Not implemented.

### 6.11 Docker Checks
- **Containers**: Missing Dockerfile.
- **Access**: Not verifiable without Docker setup.

## 7. Recommendations for Completion

1. **Authentication**:
   - Validate email uniqueness explicitly.
   - Use google/uuid for session tokens.
2. **Followers**:
   - Add notification generation for follow requests.
3. **Profile**:
   - Implement visibility check in profile retrieval.
   - Add endpoint for user posts.
4. **Posts**:
   - Support almost private/private privacy levels.
   - Validate image formats (JPEG, PNG, GIF).
5. **Groups**:
   - Implement group post functionality.
   - Add notifications for invitations and join requests.
6. **Chat**:
   - Enforce follower/public profile checks for private chats.
   - Validate emoji support.
7. **Notifications**:
   - Implement notification creation and retrieval endpoints.
8. **Database**:
   - Add migration application logic in `NewSQLiteStore`.
9. **Docker**:
   - Create a Dockerfile exposing port 3000 and running migrations.

## 8. Conclusion
The current backend implementation covers most required features but has gaps in notifications, group posts, post privacy levels, and Docker support. Addressing these gaps will ensure compliance with the project objectives and audit criteria, providing a robust learning experience in backend development, authentication, real-time communication, and containerization.

