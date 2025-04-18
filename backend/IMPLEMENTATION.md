# IMPLEMENTATION PLAN

This document outlines a phased implementation plan to address the gaps identified in the SRS.md document for the social network backend.

## Phase 1: Authentication Improvements

**Objective**: Enhance authentication security and validation.

**Tasks**:
1. Implement explicit email uniqueness check in `CreateUserAccount`
   - File: `/pkg/storage/sqlite/auth.go`
   - Add function to check if email exists before insertion

2. Replace timestamp-based session tokens with UUID
   - File: `/pkg/storage/sqlite/auth.go`
   - Update `generateSessionToken` to use google/uuid

**Expected Outcome**: More secure session management and proper validation during registration.

## Phase 2: Followers and Notifications Integration

**Objective**: Connect follower system with notifications.

**Tasks**:
1. Create notification storage implementation
   - File: `/pkg/storage/sqlite/notifications.go` (new)
   - Implement functions for creating, retrieving, and marking notifications

2. Integrate follow request notifications
   - File: `/pkg/storage/sqlite/following.go`
   - Update `CreateFollowRequest` to generate notifications for private profiles

3. Create notification handlers
   - File: `/pkg/api/handlers/notification_handlers.go` (new)
   - Implement endpoints for retrieving and managing notifications

4. Add notification routes
   - File: `/pkg/api/router.go`
   - Register notification endpoints

**Expected Outcome**: Follow requests generate notifications that users can view and manage.

## Phase 3: Profile Visibility and User Posts

**Objective**: Enforce profile privacy and create dedicated post endpoints.

**Tasks**:
1. Implement profile visibility check
   - File: `/pkg/storage/sqlite/profile.go`
   - Update `GetProfileData` to check if requester can view profile

2. Create profile-specific post endpoint
   - File: `/pkg/storage/sqlite/profile.go`
   - Add `GetUserPosts` function
   - File: `/pkg/api/handlers/profile_handlers.go`
   - Add handler for retrieving user posts in profile context

**Expected Outcome**: Private profiles are properly protected, and users can efficiently view posts by profile.

## Phase 4: Enhanced Post Privacy

**Objective**: Implement granular post privacy controls.

**Tasks**:
1. Modify posts table schema
   - File: `/pkg/db/migrations/sqlite/0012_update_posts_table.sql` (new)
   - Add privacy field with values (public, almost_private, private)

2. Update post creation and retrieval logic
   - File: `/pkg/storage/sqlite/posts.go`
   - Modify `CreatePost` to include privacy setting
   - Update `CanUserSeePost` to check for privacy level
   - Add support for private posts with selected followers

3. Implement image format validation
   - File: `/pkg/api/handlers/posts_handlers.go`
   - Add validation for JPEG, PNG, and GIF formats

**Expected Outcome**: Posts have three privacy levels (public, almost private, private) with appropriate visibility rules and image validation.

## Phase 5: Group Posts Implementation

**Objective**: Enable post creation and management within groups.

**Tasks**:
1. Implement group post functionality
   - File: `/pkg/storage/sqlite/groups.go`
   - Add `CreateGroupPost` and `GetGroupPosts` functions

2. Add group post handlers
   - File: `/pkg/api/handlers/groups_handlers.go`
   - Implement endpoints for creating and retrieving group posts

3. Integrate group notifications
   - Update group invitation and join request functions to generate notifications
   - Add event creation notifications

**Expected Outcome**: Users can create posts within groups, and group activities generate appropriate notifications.

## Phase 6: Chat System Enhancements

**Objective**: Improve chat security and feature support.

**Tasks**:
1. Implement follower/public profile check for private chats
   - File: `/pkg/websocket/websocket.go`
   - Update connection handling to validate user relationships

2. Enforce group membership in WebSocket messages
   - File: `/pkg/websocket/websocket.go`
   - Validate group membership before allowing group chat messages

3. Verify emoji support in chat messages
   - Ensure JSON encoding/decoding preserves emoji characters

**Expected Outcome**: Chat system respects visibility rules and properly supports all required features.

## Phase 7: Comprehensive Notification System

**Objective**: Complete notification system for all required events.

**Tasks**:
1. Implement remaining notification types
   - Group invitations
   - Group join requests
   - Event creation
   - Comment notifications

2. Create WebSocket notification delivery
   - File: `/pkg/websocket/websocket.go`
   - Add real-time notification support

3. Build notification management UI endpoints
   - Mark as read functionality
   - Deletion support
   - Filtering options

**Expected Outcome**: Full notification system supporting all required triggers and real-time delivery.

## Implementation Details

### Database Updates

The following tables/columns will need modification:

1. **posts table**: Add privacy field (enum: public, almost_private, private)
2. **notifications table**: Ensure it has all required fields:
   - id
   - user_id (recipient)
   - type (follow_request, group_invite, join_request, event_created)
   - content
   - related_id (follow request ID, group ID, etc.)
   - created_at
   - read (boolean)

### New Files to Create

1. `/pkg/storage/sqlite/notifications.go`
2. `/pkg/api/handlers/notification_handlers.go`
3. `/pkg/db/migrations/sqlite/0012_update_posts_table.sql`

### Files to Modify

1. `/pkg/storage/sqlite/auth.go`: Email check, UUID implementation
2. `/pkg/storage/sqlite/following.go`: Notification integration
3. `/pkg/storage/sqlite/profile.go`: Visibility check, user posts
4. `/pkg/storage/sqlite/posts.go`: Privacy controls
5. `/pkg/api/handlers/posts_handlers.go`: Image validation
6. `/pkg/storage/sqlite/groups.go`: Group posts, notifications
7. `/pkg/api/handlers/groups_handlers.go`: Group post endpoints
8. `/pkg/websocket/websocket.go`: Chat validation, notifications
9. `/pkg/api/router.go`: New endpoints registration

## Testing Strategy

For each phase, implement appropriate tests to verify functionality:

1. Unit tests for new storage functions
2. API tests for new endpoints
3. Integration tests for notification delivery
4. WebSocket tests for real-time features

## Priority Order

Implementation should follow this priority order:

1. Authentication improvements (security-critical)
2. Profile visibility and post privacy (core functionality)
3. Notification system (enables many features)
4. Group posts (commonly used feature)
5. Chat enhancements (improves existing functionality)