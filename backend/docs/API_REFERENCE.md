# Social Network API Reference

This document provides a comprehensive reference for all API endpoints in the Social Network backend application. It includes details about HTTP methods, URL paths, authentication requirements, request/response formats, and status codes.

## Table of Contents

1. [Authentication](#authentication)
2. [Profiles](#profiles)
3. [Following](#following)
4. [Posts](#posts)
5. [Groups](#groups)
6. [Chat](#chat)
7. [Notifications](#notifications)
8. [WebSocket](#websocket)

## Authentication

Authentication endpoints handle user registration, login, and logout operations.

### POST /auth/register

Register a new user account.

- **Authentication Required**: No
- **Request Format**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "date_of_birth": "1990-01-01",
    "profile_type": "public",
    "nickname": "",
    "about_me": "",
    "avatar": ""
  }
  ```
  *Note: `nickname`, `about_me`, and `avatar` are optional.*

- **Response Format**:
  ```json
  {
    "message": "User registered successfully",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "date_of_birth": "1990-01-01",
      "avatar": "",
      "nickname": "",
      "about_me": "",
      "profile_type": "public"
    }
  }
  ```
- **Status Codes**:
  - `201 Created`: Registration successful
  - `400 Bad Request`: Invalid input data
  - `409 Conflict`: Email already exists

### POST /auth/login

Log in with existing credentials.

- **Authentication Required**: No
- **Request Format**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **Response Format**:
  ```json
  {
    "message": "Login successful",
    "user_id": 1
  }
  ```
- **Status Codes**:
  - `200 OK`: Login successful
  - `400 Bad Request`: Invalid input data
  - `401 Unauthorized`: Invalid credentials
- **Notes**: Sets a `session_token` cookie that must be included in subsequent requests

### POST /auth/logout

Log out and invalidate the current session.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: No body required
- **Response Format**:
  ```json
  {
    "message": "Logout successful"
  }
  ```
- **Status Codes**:
  - `200 OK`: Logout successful
- **Notes**: Clears the `session_token` cookie

## Profiles

Profile endpoints handle user profile data operations.

### GET /profiles/{id}

Get a user's profile data.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - user ID
- **Response Format**: 
  ```json
  {
    "user_id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "nickname": "",
    "about_me": "About me text",
    "avatar": "path/to/avatar.jpg",
    "date_of_birth": "1990-01-01",
    "profile_type": "public",
    "followers": 10,
    "following": 15,
    "post_count": 5,
    "created_at": "2023-01-01T00:00:00Z",
    "last_activity": "2023-01-02T00:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Profile retrieved successfully
  - `401 Unauthorized`: Not authenticated
  - `404 Not Found`: Profile not found
- **Notes**: Returns limited data if profile is private and requester is not following

### PUT /profiles/privacy

Update the privacy setting of the authenticated user's profile.

- **Authentication Required**: Yes (session cookie)
- **Request Format**:
  ```json
  {
    "profile_type": "private"
  }
  ```
- **Response Format**:
  ```json
  {
    "message": "Privacy setting updated successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Privacy setting updated
  - `400 Bad Request`: Invalid profile type
  - `401 Unauthorized`: Not authenticated

### GET /profiles/{id}/posts

Get posts created by a specific user.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: 
  - Path parameter: user ID
  - Query parameters (optional): 
    - `limit`: Number of posts to return
    - `offset`: Number of posts to skip
- **Response Format**:
  ```json
  {
    "posts": [
      {
        "id": 1,
        "user_id": 1,
        "content": "Post content",
        "image": "path/to/image.jpg",
        "privacy_level": "public",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z",
        "likes": 5,
        "comments": []
      }
    ],
    "count": 1
  }
  ```
- **Status Codes**:
  - `200 OK`: Posts retrieved successfully
  - `400 Bad Request`: Invalid parameters
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to view these posts
- **Notes**: Respects profile and post privacy settings

## Following

Following endpoints handle the follower relationships between users.

### POST /follow/{id}

Send a follow request to a user.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - user ID to follow
- **Response Format**:
  ```json
  {
    "message": "Follow request sent successfully",
    "status": "pending"
  }
  ```
  *Note: status will be "accepted" if the user has a public profile*
- **Status Codes**:
  - `201 Created`: Follow request sent or automatically accepted
  - `400 Bad Request`: Invalid user ID
  - `401 Unauthorized`: Not authenticated
  - `409 Conflict`: Already following or pending request exists

### GET /follow/requests

Get the list of pending follow requests for the authenticated user.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: None
- **Response Format**: 
  ```json
  {
    "requests": [
      {
        "id": 1,
        "follower_id": 2,
        "follower_first_name": "Jane",
        "follower_last_name": "Doe",
        "follower_avatar": "path/to/avatar.jpg",
        "created_at": "2023-01-01T00:00:00Z"
      }
    ]
  }
  ```
- **Status Codes**:
  - `200 OK`: Requests retrieved successfully
  - `401 Unauthorized`: Not authenticated

### POST /follow/{id}/accept

Accept a follow request.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - follower user ID
- **Response Format**:
  ```json
  {
    "message": "Follow request accepted",
    "status": "accepted"
  }
  ```
- **Status Codes**:
  - `200 OK`: Request accepted
  - `400 Bad Request`: Invalid user ID
  - `401 Unauthorized`: Not authenticated
  - `404 Not Found`: No pending request found

### DELETE /follow/{id}

Reject a follow request or unfollow a user.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - follower/following user ID
- **Response Format**:
  ```json
  {
    "message": "Follow request rejected",
    "status": "rejected"
  }
  ```
  *Note: message will be "User unfollowed" if unfollowing an existing follower*
- **Status Codes**:
  - `200 OK`: Request rejected or user unfollowed
  - `400 Bad Request`: Invalid user ID
  - `401 Unauthorized`: Not authenticated
  - `404 Not Found`: No follow relationship found

## Posts

Post endpoints handle creation, retrieval, and management of posts.

### POST /posts

Create a new post.

- **Authentication Required**: Yes (session cookie)
- **Request Format**:
  ```json
  {
    "content": "Post content",
    "image": "path/to/image.jpg",
    "privacy_level": "public",
    "allowed_followers": [1, 2, 3]
  }
  ```
  *Note: `image` and `allowed_followers` are optional. `allowed_followers` is only used when `privacy_level` is "private".*
- **Response Format**:
  ```json
  {
    "post_id": 1
  }
  ```
- **Status Codes**:
  - `201 Created`: Post created successfully
  - `400 Bad Request`: Invalid input data
  - `401 Unauthorized`: Not authenticated

### GET /posts/{id}

Get a specific post.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - post ID
- **Response Format**:
  ```json
  {
    "id": 1,
    "user_id": 1,
    "user_first_name": "John",
    "user_last_name": "Doe",
    "user_avatar": "path/to/avatar.jpg",
    "content": "Post content",
    "image": "path/to/image.jpg",
    "privacy_level": "public",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "likes": 5,
    "comments": [
      {
        "id": 1,
        "post_id": 1,
        "user_id": 2,
        "user_first_name": "Jane",
        "user_last_name": "Doe",
        "user_avatar": "path/to/avatar.jpg",
        "content": "Comment content",
        "created_at": "2023-01-02T00:00:00Z"
      }
    ],
    "liked_by_current_user": true
  }
  ```
- **Status Codes**:
  - `200 OK`: Post retrieved successfully
  - `401 Unauthorized`: Not authenticated 
  - `403 Forbidden`: Not authorized to view this post
  - `404 Not Found`: Post not found

### PUT /posts/{id}

Update an existing post.

- **Authentication Required**: Yes (session cookie)
- **Request Format**:
  ```json
  {
    "content": "Updated content",
    "image": "path/to/new_image.jpg",
    "privacy_level": "private",
    "allowed_followers": [1, 2, 3]
  }
  ```
  *Note: All fields are optional. Only provided fields will be updated.*
- **Response Format**:
  ```json
  {
    "message": "Post updated successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Post updated successfully
  - `400 Bad Request`: Invalid input data
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to edit this post
  - `404 Not Found`: Post not found

### DELETE /posts/{id}

Delete a post.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - post ID
- **Response Format**:
  ```json
  {
    "message": "Post deleted successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Post deleted successfully
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to delete this post
  - `404 Not Found`: Post not found

### POST /posts/{id}/comments

Add a comment to a post.

- **Authentication Required**: Yes (session cookie)
- **Request Format**:
  ```json
  {
    "content": "Comment content",
    "image": "path/to/image.jpg"
  }
  ```
  *Note: `image` is optional.*
- **Response Format**:
  ```json
  {
    "comment_id": 1,
    "message": "Comment added successfully"
  }
  ```
- **Status Codes**:
  - `201 Created`: Comment added successfully
  - `400 Bad Request`: Invalid input data
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to comment on this post
  - `404 Not Found`: Post not found

### POST /posts/{id}/likes

Like a post.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - post ID
- **Response Format**:
  ```json
  {
    "message": "Post liked successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Post liked successfully
  - `400 Bad Request`: Already liked
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to like this post
  - `404 Not Found`: Post not found

### DELETE /posts/{id}/likes

Unlike a post.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - post ID
- **Response Format**:
  ```json
  {
    "message": "Post unliked successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Post unliked successfully
  - `400 Bad Request`: Not liked yet
  - `401 Unauthorized`: Not authenticated
  - `404 Not Found`: Post not found

## Groups

Group endpoints handle group management, membership, events, and group posts.

### POST /groups

Create a new group.

- **Authentication Required**: Yes (session cookie)
- **Request Format**:
  ```json
  {
    "title": "Group Title",
    "description": "Group Description"
  }
  ```
- **Response Format**:
  ```json
  {
    "group_id": 1,
    "message": "Group created successfully"
  }
  ```
- **Status Codes**:
  - `201 Created`: Group created successfully
  - `400 Bad Request`: Invalid input data
  - `401 Unauthorized`: Not authenticated

### GET /groups

Get a list of all groups.

- **Authentication Required**: No
- **Request Format**: 
  - Query parameters (optional):
    - `limit`: Number of groups to return
    - `offset`: Number of groups to skip
- **Response Format**: 
  ```json
  {
    "groups": [
      {
        "id": 1,
        "creator_id": 1,
        "title": "Group Title",
        "description": "Group Description",
        "member_count": 5,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      }
    ],
    "count": 1
  }
  ```
- **Status Codes**:
  - `200 OK`: Groups retrieved successfully

### GET /groups/{id}

Get details of a specific group.

- **Authentication Required**: No
- **Request Format**: Path parameter - group ID
- **Response Format**:
  ```json
  {
    "id": 1,
    "creator_id": 1,
    "creator_first_name": "John",
    "creator_last_name": "Doe",
    "title": "Group Title",
    "description": "Group Description",
    "member_count": 5,
    "is_member": true,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Group retrieved successfully
  - `404 Not Found`: Group not found

### POST /groups/{id}/invite

Invite a user to join a group.

- **Authentication Required**: Yes (session cookie)
- **Request Format**:
  ```json
  {
    "invitee_id": 2
  }
  ```
- **Response Format**:
  ```json
  {
    "message": "User invited successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Invitation sent successfully
  - `400 Bad Request`: Invalid input data
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to invite users
  - `404 Not Found`: Group or user not found
  - `409 Conflict`: User already invited or a member

### POST /groups/invites/{id}/accept

Accept a group invitation.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - invitation ID
- **Response Format**:
  ```json
  {
    "message": "Invitation accepted successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Invitation accepted
  - `400 Bad Request`: Invalid invitation ID
  - `401 Unauthorized`: Not authenticated
  - `404 Not Found`: Invitation not found

### POST /groups/invites/{id}/reject

Reject a group invitation.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - invitation ID
- **Response Format**:
  ```json
  {
    "message": "Invitation rejected successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Invitation rejected
  - `400 Bad Request`: Invalid invitation ID
  - `401 Unauthorized`: Not authenticated
  - `404 Not Found`: Invitation not found

### GET /groups/search

Search for groups by title or description.

- **Authentication Required**: No
- **Request Format**: Query parameter - `q` (search term)
- **Response Format**: Same as GET /groups
- **Status Codes**:
  - `200 OK`: Search successful
  - `400 Bad Request`: Missing search term

### POST /groups/{id}/join

Request to join a group.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - group ID
- **Response Format**:
  ```json
  {
    "message": "Join request sent successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Request sent successfully
  - `400 Bad Request`: Invalid group ID
  - `401 Unauthorized`: Not authenticated
  - `404 Not Found`: Group not found
  - `409 Conflict`: Already a member or request pending

### GET /groups/{id}/requests

Get a list of pending join requests for a group.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - group ID
- **Response Format**:
  ```json
  {
    "requests": [
      {
        "id": 1,
        "user_id": 2,
        "user_first_name": "Jane",
        "user_last_name": "Doe",
        "user_avatar": "path/to/avatar.jpg",
        "created_at": "2023-01-01T00:00:00Z"
      }
    ]
  }
  ```
- **Status Codes**:
  - `200 OK`: Requests retrieved successfully
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to view requests
  - `404 Not Found`: Group not found

### POST /groups/{id}/requests/{userId}/accept

Accept a join request.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameters - group ID, user ID
- **Response Format**:
  ```json
  {
    "message": "Join request accepted"
  }
  ```
- **Status Codes**:
  - `200 OK`: Request accepted
  - `400 Bad Request`: Invalid parameters
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to accept requests
  - `404 Not Found`: Group, user, or request not found

### POST /groups/{id}/requests/{userId}/reject

Reject a join request.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameters - group ID, user ID
- **Response Format**:
  ```json
  {
    "message": "Join request rejected"
  }
  ```
- **Status Codes**:
  - `200 OK`: Request rejected
  - `400 Bad Request`: Invalid parameters
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to reject requests
  - `404 Not Found`: Group, user, or request not found

### GET /groups/{id}/members

Get a list of group members.

- **Authentication Required**: No
- **Request Format**: Path parameter - group ID
- **Response Format**:
  ```json
  {
    "members": [
      {
        "user_id": 1,
        "first_name": "John",
        "last_name": "Doe",
        "avatar": "path/to/avatar.jpg",
        "joined_at": "2023-01-01T00:00:00Z",
        "is_creator": true
      }
    ]
  }
  ```
- **Status Codes**:
  - `200 OK`: Members retrieved successfully
  - `404 Not Found`: Group not found

### POST /groups/{id}/events

Create a new event for a group.

- **Authentication Required**: Yes (session cookie)
- **Request Format**:
  ```json
  {
    "title": "Event Title",
    "description": "Event Description",
    "event_time": "2023-02-01T15:00:00Z"
  }
  ```
- **Response Format**:
  ```json
  {
    "event_id": 1,
    "message": "Event created successfully"
  }
  ```
- **Status Codes**:
  - `201 Created`: Event created successfully
  - `400 Bad Request`: Invalid input data
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to create events
  - `404 Not Found`: Group not found

### GET /groups/{id}/events

Get a list of events for a group.

- **Authentication Required**: No
- **Request Format**: Path parameter - group ID
- **Response Format**:
  ```json
  {
    "events": [
      {
        "id": 1,
        "group_id": 1,
        "creator_id": 1,
        "creator_first_name": "John",
        "creator_last_name": "Doe",
        "title": "Event Title",
        "description": "Event Description",
        "event_time": "2023-02-01T15:00:00Z",
        "created_at": "2023-01-01T00:00:00Z",
        "going_count": 3,
        "not_going_count": 1,
        "user_response": "going"
      }
    ]
  }
  ```
- **Status Codes**:
  - `200 OK`: Events retrieved successfully
  - `404 Not Found`: Group not found

### POST /groups/{id}/events/{eventId}/response

Respond to a group event.

- **Authentication Required**: Yes (session cookie)
- **Request Format**:
  ```json
  {
    "response": "going"
  }
  ```
  *Note: `response` must be either "going" or "not_going".*
- **Response Format**:
  ```json
  {
    "message": "Response recorded successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Response recorded successfully
  - `400 Bad Request`: Invalid input data
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not a group member
  - `404 Not Found`: Group or event not found

### POST /groups/{id}/posts

Create a post in a group.

- **Authentication Required**: Yes (session cookie)
- **Request Format**:
  ```json
  {
    "content": "Post content",
    "image": "path/to/image.jpg"
  }
  ```
  *Note: `image` is optional.*
- **Response Format**:
  ```json
  {
    "post_id": 1,
    "message": "Post created successfully"
  }
  ```
- **Status Codes**:
  - `201 Created`: Post created successfully
  - `400 Bad Request`: Invalid input data
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not a group member
  - `404 Not Found`: Group not found

### GET /groups/{id}/posts

Get posts from a group.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: 
  - Path parameter: group ID
  - Query parameters (optional):
    - `limit`: Number of posts to return
    - `offset`: Number of posts to skip
- **Response Format**:
  ```json
  {
    "posts": [
      {
        "id": 1,
        "user_id": 1,
        "user_first_name": "John",
        "user_last_name": "Doe",
        "user_avatar": "path/to/avatar.jpg",
        "content": "Post content",
        "image": "path/to/image.jpg",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z",
        "likes": 5,
        "comments": [],
        "liked_by_current_user": false
      }
    ],
    "count": 1
  }
  ```
- **Status Codes**:
  - `200 OK`: Posts retrieved successfully
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not a group member
  - `404 Not Found`: Group not found

## Chat

Chat endpoints handle direct messaging and group chat functionality.

### GET /chats

Get a list of the authenticated user's chat conversations.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: None
- **Response Format**:
  ```json
  {
    "chats": [
      {
        "user_id": 2,
        "first_name": "Jane",
        "last_name": "Doe",
        "avatar": "path/to/avatar.jpg",
        "last_message": "Hello there",
        "last_message_time": "2023-01-02T12:34:56Z",
        "unread_count": 2
      }
    ]
  }
  ```
- **Status Codes**:
  - `200 OK`: Chats retrieved successfully
  - `401 Unauthorized`: Not authenticated

### GET /chats/{userId}

Get chat history with a specific user.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: 
  - Path parameter: user ID
  - Query parameters (optional):
    - `limit`: Number of messages to return
    - `before`: Timestamp to get messages before
- **Response Format**:
  ```json
  {
    "messages": [
      {
        "id": 1,
        "sender_id": 1,
        "receiver_id": 2,
        "content": "Hello",
        "created_at": "2023-01-01T12:00:00Z"
      },
      {
        "id": 2,
        "sender_id": 2,
        "receiver_id": 1,
        "content": "Hi there",
        "created_at": "2023-01-01T12:01:00Z"
      }
    ],
    "user": {
      "id": 2,
      "first_name": "Jane",
      "last_name": "Doe",
      "avatar": "path/to/avatar.jpg"
    }
  }
  ```
- **Status Codes**:
  - `200 OK`: Messages retrieved successfully
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to chat with this user
  - `404 Not Found`: User not found

### GET /groups/{id}/chat

Get chat history for a group.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: 
  - Path parameter: group ID
  - Query parameters (optional):
    - `limit`: Number of messages to return
    - `before`: Timestamp to get messages before
- **Response Format**:
  ```json
  {
    "messages": [
      {
        "id": 1,
        "group_id": 1,
        "sender_id": 1,
        "sender_first_name": "John",
        "sender_last_name": "Doe",
        "sender_avatar": "path/to/avatar.jpg",
        "content": "Hello group",
        "created_at": "2023-01-01T12:00:00Z"
      }
    ],
    "group": {
      "id": 1,
      "title": "Group Title"
    }
  }
  ```
- **Status Codes**:
  - `200 OK`: Messages retrieved successfully
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not a group member
  - `404 Not Found`: Group not found

## Notifications

Notification endpoints handle user notifications for various events.

### GET /notifications

Get notifications for the authenticated user.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: 
  - Query parameters (optional):
    - `limit`: Number of notifications to return
    - `offset`: Number of notifications to skip
    - `unread_only`: Boolean flag to filter unread notifications
- **Response Format**:
  ```json
  {
    "notifications": [
      {
        "id": 1,
        "user_id": 1,
        "type": "follow_request",
        "content": "Jane Doe requested to follow you",
        "is_read": false,
        "related_id": 2,
        "created_at": "2023-01-01T00:00:00Z",
        "sender_id": 2,
        "sender_name": "Jane Doe",
        "sender_image": "path/to/avatar.jpg"
      }
    ],
    "count": 1
  }
  ```
- **Status Codes**:
  - `200 OK`: Notifications retrieved successfully
  - `401 Unauthorized`: Not authenticated

### PUT /notifications/{id}/read

Mark a notification as read.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: Path parameter - notification ID
- **Response Format**:
  ```json
  {
    "message": "Notification marked as read"
  }
  ```
- **Status Codes**:
  - `200 OK`: Notification marked as read
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to access this notification
  - `404 Not Found`: Notification not found

### PUT /notifications/read-all

Mark all notifications as read.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: None
- **Response Format**:
  ```json
  {
    "message": "All notifications marked as read",
    "count": 5
  }
  ```
- **Status Codes**:
  - `200 OK`: Notifications marked as read
  - `401 Unauthorized`: Not authenticated

### GET /notifications/unread-count

Get the count of unread notifications.

- **Authentication Required**: Yes (session cookie)
- **Request Format**: None
- **Response Format**:
  ```json
  {
    "count": 5
  }
  ```
- **Status Codes**:
  - `200 OK`: Count retrieved successfully
  - `401 Unauthorized`: Not authenticated

## WebSocket

WebSocket endpoints provide real-time communication for chat and notifications.

### WebSocket: /ws/chat/{userId}

Establish a WebSocket connection for direct chat with a user.

- **Authentication Required**: Yes (session cookie)
- **Connection Parameters**: Path parameter - user ID to chat with
- **Message Format (Client to Server)**:
  ```json
  {
    "type": "message",
    "content": "Hello there"
  }
  ```
- **Message Format (Server to Client)**:
  ```json
  {
    "type": "message",
    "sender_id": 2,
    "content": "Hi back",
    "created_at": "2023-01-01T12:01:00Z"
  }
  ```

### WebSocket: /ws/groups/{id}/chat

Establish a WebSocket connection for group chat.

- **Authentication Required**: Yes (session cookie)
- **Connection Parameters**: Path parameter - group ID
- **Message Format (Client to Server)**:
  ```json
  {
    "type": "message",
    "content": "Hello group"
  }
  ```
- **Message Format (Server to Client)**:
  ```json
  {
    "type": "message",
    "sender_id": 2,
    "sender_name": "Jane Doe",
    "content": "Hello everyone",
    "created_at": "2023-01-01T12:01:00Z"
  }
  ```

### WebSocket: /ws/notifications

Establish a WebSocket connection for real-time notifications.

- **Authentication Required**: Yes (session cookie)
- **Message Format (Server to Client)**:
  ```json
  {
    "id": 1,
    "type": "follow_request",
    "content": "Jane Doe requested to follow you",
    "sender_id": 2,
    "sender_name": "Jane Doe",
    "created_at": "2023-01-01T00:00:00Z"
  }
  ```
- **Notes**: This is a one-way channel from server to client