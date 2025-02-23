# API Endpoints:
## Authentication (/auth)
---
#### - POST /auth/register
    - Process: Register a new user with email, password, name, date of birth, and optional fields like avatar. Hash the password using bcrypt, generate a UUID for user ID, and store in the database. Return session token.
#### - POST /auth/login
    - Process: Authenticate user with email and password. If valid, create a session and set a secure cookie. Return session token.
#### - POST /auth/logout
    - Process: Clear the session cookie, effectively logging the user out.

## Profile  (/profiles)
---
#### - GET /profiles/{id}
    - Process: Retrieve user details. Middleware checks if the profile is public or if the requester is a follower. Return user info, activity, and follow lists accordingly.
#### - PUT /profiles/privacy
    - Process: Toggle profile visibility between public and private. Only accessible by the profile owner.
#### - GET /profiles/{id}/activity
    - Process: Fetch and return paginated posts and comments by the user. Visibility depends on the profile's privacy settings.

## Posts  (/posts)
---
#### - POST /posts
    - Process: Create a new post with privacy settings (public, almost private, private). Handle image/GIF uploads if included.
#### - GET /posts/{id}
    - Process: Retrieve a post. Middleware ensures the requester has permission to view based on post privacy.
#### - PUT /posts/{id}
    - Process: Edit a post. Only the post owner can perform this action.
#### - DELETE /posts/{id}
    - Process: Delete a post. Similar to edit, only the owner can delete.
#### - POST /posts/{id}/comments
    - Process: Add a comment to a post. Ensure the user has permission to comment based on post visibility.
#### - POST /posts/{id}/likes
    - Process: Like a post. Handle notifications for this action via WebSocket.
        
## Following  (/following)
---
#### - POST /follow/{id}
    - Process: Request to follow another user. If the profile is public, automatically approve; otherwise, create a pending request.
#### - GET /follow/requests
    - Process: List all pending follow requests for the authenticated user.
#### - POST /follow/{id}/accept
    - Process: Accept a follow request for a private profile.
#### - POST /follow/{id}
    - Process: Unfollow or reject a follow request.
        
## Groups (/groups)
---
#### - POST /groups
    - Process: Create a new group with title and description.
#### - POST /groups/{id}/invite
    - Process: Send invitations to join the group. Notifications are sent to invitees.
#### - GET /groups/search
    - Process: Search for groups by name or description.
#### - POST /groups/{id}/join
    - Process: Request to join a group. Admin will see this in notifications.
#### - GET /groups/{id}/requests
    - Process: List join requests for the group (admin view).
#### - POST /groups/{id}/requests/{userId}/accept
    - Process: Accept a user's request to join the group.
#### - POST /groups/{id}/events
    - Process: Create an event within the group with details like title, description, date/time, and options for attendance.
        
## Chat  (/ws)
---
#### - WebSocket /ws
    - Process: Handle both direct messages and group chat. Messages are categorized by type (chat, notification, comment, etc.). Use separate handlers for different message types.
## Notifications (/ws)
---
#### - WebSocket /ws
    - Process: Real-time notifications for various actions like likes, comments, messages, group invitations, etc. Use the same WebSocket connection for efficiency but with distinct message types.
