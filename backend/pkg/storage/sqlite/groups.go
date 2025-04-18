package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"socialNetwork/pkg/domain/groups"
	"socialNetwork/pkg/domain/notifications"
	"socialNetwork/pkg/domain/posts"
	"socialNetwork/pkg/errors"
)

// CreateGroup creates a new group
func (s *SQLiteStore) CreateGroup(userID int64, req *groups.CreateGroupRequest) (int64, error) {
	query := `INSERT INTO groups (creator_id, title, description, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := s.db.Exec(query, userID, req.Title, req.Description, now, now)
	if err != nil {
		log.Print("Error creating group:", err)
		return 0, errors.ErrInternalServer
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		log.Print("Error getting last insert ID:", err)
		return 0, errors.ErrInternalServer
	}

	// Add creator as a member with accepted status
	memberQuery := `INSERT INTO group_members (group_id, user_id, status, created_at) 
                   VALUES (?, ?, ?, ?)`
	_, err = s.db.Exec(memberQuery, groupID, userID, "accepted", now)
	if err != nil {
		log.Print("Error adding creator as group member:", err)
		return 0, errors.ErrInternalServer
	}

	return groupID, nil
}

// GetGroup retrieves a group by ID
func (s *SQLiteStore) GetGroup(groupID int64) (*groups.Group, error) {
	query := `SELECT id, creator_id, title, description, created_at, updated_at 
              FROM groups WHERE id = ?`

	var group groups.Group
	err := s.db.QueryRow(query, groupID).Scan(
		&group.ID,
		&group.CreatorID,
		&group.Title,
		&group.Description,
		&group.CreatedAt,
		&group.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		log.Print("Error retrieving group:", err)
		return nil, errors.ErrInternalServer
	}

	return &group, nil
}

// GetGroups retrieves all groups
func (s *SQLiteStore) GetGroups() ([]*groups.Group, error) {
	// Check if the database has groups table
	var count int
	err := s.db.QueryRow(`SELECT count(name) FROM sqlite_master 
                          WHERE type='table' AND name='groups'`).Scan(&count)
	if err != nil {
		log.Print("Error checking for groups table:", err)
		return nil, errors.ErrInternalServer
	}

	// If groups table doesn't exist, return empty slice
	if count == 0 {
		return []*groups.Group{}, nil
	}

	// Check if there are any groups
	err = s.db.QueryRow(`SELECT COUNT(*) FROM groups`).Scan(&count)
	if err != nil {
		log.Print("Error counting groups:", err)
		return nil, errors.ErrInternalServer
	}

	// If no groups exist, return empty slice
	if count == 0 {
		return []*groups.Group{}, nil
	}

	// Query all groups
	query := `SELECT id, creator_id, title, description, created_at, updated_at
              FROM groups 
              WHERE id IS NOT NULL
              ORDER BY created_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		log.Print("Error querying groups:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	var groupsList []*groups.Group
	for rows.Next() {
		var group groups.Group
		err := rows.Scan(
			&group.ID,
			&group.CreatorID,
			&group.Title,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			log.Print("Error scanning group row:", err)
			return nil, errors.ErrInternalServer
		}
		groupsList = append(groupsList, &group)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating group rows:", err)
		return nil, errors.ErrInternalServer
	}

	return groupsList, nil
}

// SearchGroups searches groups by title or description
func (s *SQLiteStore) SearchGroups(query string) ([]*groups.Group, error) {
	// Check if the database has groups table
	var count int
	err := s.db.QueryRow(`SELECT count(name) FROM sqlite_master 
                          WHERE type='table' AND name='groups'`).Scan(&count)
	if err != nil {
		log.Print("Error checking for groups table:", err)
		return nil, errors.ErrInternalServer
	}

	// If groups table doesn't exist, return empty slice
	if count == 0 {
		return []*groups.Group{}, nil
	}

	// Check if there are any groups
	err = s.db.QueryRow(`SELECT COUNT(*) FROM groups`).Scan(&count)
	if err != nil {
		log.Print("Error counting groups:", err)
		return nil, errors.ErrInternalServer
	}

	// If no groups exist, return empty slice
	if count == 0 {
		return []*groups.Group{}, nil
	}

	// Search groups
	sqlQuery := `SELECT id, creator_id, title, description, created_at, updated_at
                FROM groups 
                WHERE (title LIKE ? OR description LIKE ?) AND id IS NOT NULL
                ORDER BY created_at DESC`

	searchParam := "%" + query + "%"
	rows, err := s.db.Query(sqlQuery, searchParam, searchParam)
	if err != nil {
		log.Print("Error searching groups:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	var groupsList []*groups.Group
	for rows.Next() {
		var group groups.Group
		err := rows.Scan(
			&group.ID,
			&group.CreatorID,
			&group.Title,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			log.Print("Error scanning search result:", err)
			return nil, errors.ErrInternalServer
		}
		groupsList = append(groupsList, &group)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating search results:", err)
		return nil, errors.ErrInternalServer
	}

	return groupsList, nil
}

// Group Membership

// InviteToGroup invites a user to a group
func (s *SQLiteStore) InviteToGroup(groupID, inviterID, inviteeID int64) error {
	// Check if inviter is a member of the group
	isMember, err := s.IsGroupMember(groupID, inviterID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.ErrUnauthorized
	}

	// Check if invitee is already a member or has a pending invitation
	query := `SELECT COUNT(*) FROM group_members 
              WHERE group_id = ? AND user_id = ?`

	var count int
	err = s.db.QueryRow(query, groupID, inviteeID).Scan(&count)
	if err != nil {
		log.Print("Error checking existing membership:", err)
		return errors.ErrInternalServer
	}

	if count > 0 {
		return errors.ErrConflict
	}

	// Create invitation (pending membership)
	insertQuery := `INSERT INTO group_members (group_id, user_id, status, created_at) 
                   VALUES (?, ?, ?, ?)`
	_, err = s.db.Exec(insertQuery, groupID, inviteeID, "pending", time.Now())
	if err != nil {
		log.Print("Error creating group invitation:", err)
		return errors.ErrInternalServer
	}
	
	// Create a notification for the invitee
	// Get group info
	group, err := s.GetGroup(groupID)
	if err != nil {
		log.Print("Error getting group info for notification:", err)
		// Continue even if notification creation fails
	} else {
		// Get inviter name
		var firstName, lastName string
		userQuery := `SELECT first_name, last_name FROM users WHERE id = ?`
		err = s.db.QueryRow(userQuery, inviterID).Scan(&firstName, &lastName)
		if err != nil {
			log.Print("Error getting inviter name:", err)
			// Continue even if notification creation fails
		} else {
			// Create notification
			notif := &notifications.Notification{
				UserID:    inviteeID,
				Type:      notifications.TypeGroupInvite,
				Content:   fmt.Sprintf("%s %s invited you to join group: %s", firstName, lastName, group.Title),
				RelatedID: groupID,
				SenderID:  inviterID,
				CreatedAt: time.Now(),
			}
			
			err = s.CreateNotification(notif)
			if err != nil {
				log.Print("Error creating group invitation notification:", err)
				// Continue even if notification creation fails
			}
		}
	}

	return nil
}

// GetGroupInvites gets all groups a user has been invited to
func (s *SQLiteStore) GetGroupInvites(userID int64) ([]*groups.Group, error) {
	query := `SELECT g.id, g.creator_id, g.title, g.description, g.created_at, g.updated_at
              FROM groups g
              JOIN group_members gm ON g.id = gm.group_id
              WHERE gm.user_id = ? AND gm.status = 'pending'
              ORDER BY gm.created_at DESC`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		log.Print("Error querying group invites:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	var groupsList []*groups.Group
	for rows.Next() {
		var group groups.Group
		err := rows.Scan(
			&group.ID,
			&group.CreatorID,
			&group.Title,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			log.Print("Error scanning invite row:", err)
			return nil, errors.ErrInternalServer
		}
		groupsList = append(groupsList, &group)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating invite rows:", err)
		return nil, errors.ErrInternalServer
	}

	return groupsList, nil
}

// AcceptGroupInvite accepts a group invitation
func (s *SQLiteStore) AcceptGroupInvite(groupID, userID int64) error {
	// Check if the invitation exists
	query := `SELECT COUNT(*) FROM group_members 
              WHERE group_id = ? AND user_id = ? AND status = 'pending'`

	var count int
	err := s.db.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		log.Print("Error checking invitation:", err)
		return errors.ErrInternalServer
	}

	if count == 0 {
		return errors.ErrNotFound
	}

	// Update to accepted status
	updateQuery := `UPDATE group_members SET status = 'accepted' 
                   WHERE group_id = ? AND user_id = ?`
	_, err = s.db.Exec(updateQuery, groupID, userID)
	if err != nil {
		log.Print("Error accepting invitation:", err)
		return errors.ErrInternalServer
	}

	return nil
}

// RejectGroupInvite rejects a group invitation
func (s *SQLiteStore) RejectGroupInvite(groupID, userID int64) error {
	// Check if the invitation exists
	query := `SELECT COUNT(*) FROM group_members 
              WHERE group_id = ? AND user_id = ? AND status = 'pending'`

	var count int
	err := s.db.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		log.Print("Error checking invitation:", err)
		return errors.ErrInternalServer
	}

	if count == 0 {
		return errors.ErrNotFound
	}

	// Delete the membership record
	deleteQuery := `DELETE FROM group_members 
                   WHERE group_id = ? AND user_id = ?`
	_, err = s.db.Exec(deleteQuery, groupID, userID)
	if err != nil {
		log.Print("Error rejecting invitation:", err)
		return errors.ErrInternalServer
	}

	return nil
}

// RequestJoinGroup creates a request to join a group
func (s *SQLiteStore) RequestJoinGroup(groupID, userID int64) error {
	// Check if user is already a member or has a pending request
	query := `SELECT COUNT(*) FROM group_members 
              WHERE group_id = ? AND user_id = ?`

	var count int
	err := s.db.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		log.Print("Error checking existing membership:", err)
		return errors.ErrInternalServer
	}

	if count > 0 {
		return errors.ErrConflict
	}

	// Create join request
	insertQuery := `INSERT INTO group_members (group_id, user_id, status, created_at) 
                   VALUES (?, ?, ?, ?)`
	_, err = s.db.Exec(insertQuery, groupID, userID, "requested", time.Now())
	if err != nil {
		log.Print("Error creating join request:", err)
		return errors.ErrInternalServer
	}
	
	// Create a notification for the group creator
	// Get group info
	group, err := s.GetGroup(groupID)
	if err != nil {
		log.Print("Error getting group info for notification:", err)
		// Continue even if notification creation fails
	} else {
		// Get requester name
		var firstName, lastName string
		userQuery := `SELECT first_name, last_name FROM users WHERE id = ?`
		err = s.db.QueryRow(userQuery, userID).Scan(&firstName, &lastName)
		if err != nil {
			log.Print("Error getting requester name:", err)
			// Continue even if notification creation fails
		} else {
			// Create notification
			notif := &notifications.Notification{
				UserID:    group.CreatorID,
				Type:      notifications.TypeJoinRequest,
				Content:   fmt.Sprintf("%s %s wants to join your group: %s", firstName, lastName, group.Title),
				RelatedID: groupID,
				SenderID:  userID,
				CreatedAt: time.Now(),
			}
			
			err = s.CreateNotification(notif)
			if err != nil {
				log.Print("Error creating join request notification:", err)
				// Continue even if notification creation fails
			}
		}
	}

	return nil
}

// GetGroupJoinRequests gets all join requests for a group
func (s *SQLiteStore) GetGroupJoinRequests(groupID int64) ([]*groups.GroupMember, error) {
	query := `SELECT gm.id, gm.group_id, gm.user_id, gm.status, gm.created_at
              FROM group_members gm
              WHERE gm.group_id = ? AND gm.status = 'requested'
              ORDER BY gm.created_at DESC`

	rows, err := s.db.Query(query, groupID)
	if err != nil {
		log.Print("Error querying join requests:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	var membersList []*groups.GroupMember
	for rows.Next() {
		var member groups.GroupMember
		err := rows.Scan(
			&member.ID,
			&member.GroupID,
			&member.UserID,
			&member.Status,
			&member.CreatedAt,
		)
		if err != nil {
			log.Print("Error scanning request row:", err)
			return nil, errors.ErrInternalServer
		}
		membersList = append(membersList, &member)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating request rows:", err)
		return nil, errors.ErrInternalServer
	}

	return membersList, nil
}

// AcceptGroupJoinRequest accepts a request to join a group
func (s *SQLiteStore) AcceptGroupJoinRequest(groupID, userID int64) error {
	// Check if the request exists
	query := `SELECT COUNT(*) FROM group_members 
              WHERE group_id = ? AND user_id = ? AND status = 'requested'`

	var count int
	err := s.db.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		log.Print("Error checking join request:", err)
		return errors.ErrInternalServer
	}

	if count == 0 {
		return errors.ErrNotFound
	}

	// Update to accepted status
	updateQuery := `UPDATE group_members SET status = 'accepted' 
                   WHERE group_id = ? AND user_id = ?`
	_, err = s.db.Exec(updateQuery, groupID, userID)
	if err != nil {
		log.Print("Error accepting join request:", err)
		return errors.ErrInternalServer
	}
	
	// Create a notification for the user
	// Get group info
	group, err := s.GetGroup(groupID)
	if err != nil {
		log.Print("Error getting group info for notification:", err)
		// Continue even if notification creation fails
	} else {
		// Create notification
		notif := &notifications.Notification{
			UserID:    userID,
			Type:      "join_request_accepted",
			Content:   fmt.Sprintf("Your request to join %s has been accepted", group.Title),
			RelatedID: groupID,
			SenderID:  group.CreatorID,
			CreatedAt: time.Now(),
		}
		
		err = s.CreateNotification(notif)
		if err != nil {
			log.Print("Error creating join acceptance notification:", err)
			// Continue even if notification creation fails
		}
	}

	return nil
}

// RejectGroupJoinRequest rejects a request to join a group
func (s *SQLiteStore) RejectGroupJoinRequest(groupID, userID int64) error {
	// Check if the request exists
	query := `SELECT COUNT(*) FROM group_members 
              WHERE group_id = ? AND user_id = ? AND status = 'requested'`

	var count int
	err := s.db.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		log.Print("Error checking join request:", err)
		return errors.ErrInternalServer
	}

	if count == 0 {
		return errors.ErrNotFound
	}

	// Delete the membership record
	deleteQuery := `DELETE FROM group_members 
                   WHERE group_id = ? AND user_id = ?`
	_, err = s.db.Exec(deleteQuery, groupID, userID)
	if err != nil {
		log.Print("Error rejecting join request:", err)
		return errors.ErrInternalServer
	}

	return nil
}

// IsGroupMember checks if a user is a member of a group
func (s *SQLiteStore) IsGroupMember(groupID, userID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM group_members 
              WHERE group_id = ? AND user_id = ? AND status = 'accepted'`

	var count int
	err := s.db.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		log.Print("Error checking group membership:", err)
		return false, errors.ErrInternalServer
	}

	return count > 0, nil
}

// IsGroupCreator checks if a user is the creator of a group
func (s *SQLiteStore) IsGroupCreator(groupID, userID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM groups WHERE id = ? AND creator_id = ?`

	var count int
	err := s.db.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		log.Print("Error checking group creator:", err)
		return false, errors.ErrInternalServer
	}

	return count > 0, nil
}

// GetGroupMembers gets all members of a group
func (s *SQLiteStore) GetGroupMembers(groupID int64) ([]*groups.GroupMember, error) {
	query := `SELECT gm.id, gm.group_id, gm.user_id, gm.status, gm.created_at
              FROM group_members gm
              WHERE gm.group_id = ? AND gm.status = 'accepted'
              ORDER BY gm.created_at`

	rows, err := s.db.Query(query, groupID)
	if err != nil {
		log.Print("Error querying group members:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	var membersList []*groups.GroupMember
	for rows.Next() {
		var member groups.GroupMember
		err := rows.Scan(
			&member.ID,
			&member.GroupID,
			&member.UserID,
			&member.Status,
			&member.CreatedAt,
		)
		if err != nil {
			log.Print("Error scanning member row:", err)
			return nil, errors.ErrInternalServer
		}
		membersList = append(membersList, &member)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating member rows:", err)
		return nil, errors.ErrInternalServer
	}

	return membersList, nil
}

// Group Events

// CreateGroupEvent creates a new event in a group
func (s *SQLiteStore) CreateGroupEvent(event *groups.GroupEvent) (int64, error) {
	// Check if the creator is a member of the group
	isMember, err := s.IsGroupMember(event.GroupID, event.CreatorID)
	if err != nil {
		return 0, err
	}
	if !isMember {
		return 0, errors.ErrUnauthorized
	}

	query := `INSERT INTO events (group_id, creator_id, title, description, event_time, created_at) 
              VALUES (?, ?, ?, ?, ?, ?)`

	result, err := s.db.Exec(
		query,
		event.GroupID,
		event.CreatorID,
		event.Title,
		event.Description,
		event.EventTime,
		time.Now(),
	)
	if err != nil {
		log.Print("Error creating event:", err)
		return 0, errors.ErrInternalServer
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		log.Print("Error getting last insert ID:", err)
		return 0, errors.ErrInternalServer
	}
	
	// Create notifications for all group members
	// Get group info
	group, err := s.GetGroup(event.GroupID)
	if err != nil {
		log.Print("Error getting group info for notification:", err)
		// Continue even if notification creation fails
	} else {
		// Get all group members
		members, err := s.GetGroupMembers(event.GroupID)
		if err != nil {
			log.Print("Error getting group members for notifications:", err)
			// Continue even if notification creation fails
		} else {
			// Get creator name
			var firstName, lastName string
			userQuery := `SELECT first_name, last_name FROM users WHERE id = ?`
			err = s.db.QueryRow(userQuery, event.CreatorID).Scan(&firstName, &lastName)
			if err != nil {
				log.Print("Error getting creator name:", err)
				// Continue even if notification creation fails
			} else {
				// Create notifications for each member except the creator
				for _, member := range members {
					if member.UserID != event.CreatorID {
						notif := &notifications.Notification{
							UserID:    member.UserID,
							Type:      notifications.TypeEventCreated,
							Content:   fmt.Sprintf("%s %s created a new event in %s: %s", firstName, lastName, group.Title, event.Title),
							RelatedID: eventID,
							SenderID:  event.CreatorID,
							CreatedAt: time.Now(),
						}
						
						// We don't check for errors here - if one notification fails, we still want to create others
						s.CreateNotification(notif)
					}
				}
			}
		}
	}

	return eventID, nil
}

// GetGroupEvents gets all events in a group
func (s *SQLiteStore) GetGroupEvents(groupID int64) ([]*groups.GroupEvent, error) {
	query := `SELECT id, group_id, creator_id, title, description, event_time, created_at
              FROM events
              WHERE group_id = ?
              ORDER BY event_time`

	rows, err := s.db.Query(query, groupID)
	if err != nil {
		log.Print("Error querying group events:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	var eventsList []*groups.GroupEvent
	for rows.Next() {
		var event groups.GroupEvent
		err := rows.Scan(
			&event.ID,
			&event.GroupID,
			&event.CreatorID,
			&event.Title,
			&event.Description,
			&event.EventTime,
			&event.CreatedAt,
		)
		if err != nil {
			log.Print("Error scanning event row:", err)
			return nil, errors.ErrInternalServer
		}
		eventsList = append(eventsList, &event)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating event rows:", err)
		return nil, errors.ErrInternalServer
	}

	return eventsList, nil
}

// RespondToEvent records a user's response to an event
func (s *SQLiteStore) RespondToEvent(eventID, userID int64, response string) error {
	// Check if the user has already responded to this event
	query := `SELECT COUNT(*) FROM event_responses
              WHERE event_id = ? AND user_id = ?`

	var count int
	err := s.db.QueryRow(query, eventID, userID).Scan(&count)
	if err != nil {
		log.Print("Error checking event response:", err)
		return errors.ErrInternalServer
	}

	if count > 0 {
		// Update existing response
		updateQuery := `UPDATE event_responses
                       SET response = ?
                       WHERE event_id = ? AND user_id = ?`
		_, err = s.db.Exec(updateQuery, response, eventID, userID)
		if err != nil {
			log.Print("Error updating event response:", err)
			return errors.ErrInternalServer
		}
	} else {
		// Create new response
		insertQuery := `INSERT INTO event_responses (event_id, user_id, response, created_at)
                       VALUES (?, ?, ?, ?)`
		_, err = s.db.Exec(insertQuery, eventID, userID, response, time.Now())
		if err != nil {
			log.Print("Error creating event response:", err)
			return errors.ErrInternalServer
		}
	}

	return nil
}

// GetEventResponses gets all responses to an event
func (s *SQLiteStore) GetEventResponses(eventID int64) ([]*groups.EventResponse, error) {
	query := `SELECT id, event_id, user_id, response, created_at
              FROM event_responses
              WHERE event_id = ?
              ORDER BY created_at`

	rows, err := s.db.Query(query, eventID)
	if err != nil {
		log.Print("Error querying event responses:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	var responsesList []*groups.EventResponse
	for rows.Next() {
		var eventResponse groups.EventResponse
		err := rows.Scan(
			&eventResponse.ID,
			&eventResponse.EventID,
			&eventResponse.UserID,
			&eventResponse.Response,
			&eventResponse.CreatedAt,
		)
		if err != nil {
			log.Print("Error scanning response row:", err)
			return nil, errors.ErrInternalServer
		}
		responsesList = append(responsesList, &eventResponse)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating response rows:", err)
		return nil, errors.ErrInternalServer
	}

	return responsesList, nil
}

// Group Posts

// CreateGroupPost creates a new post in a group
func (s *SQLiteStore) CreateGroupPost(groupID, userID int64, post *posts.Post) (int64, error) {
	// First, check if the user is a member of the group
	isMember, err := s.IsGroupMember(groupID, userID)
	if err != nil {
		return 0, err
	}
	if !isMember {
		return 0, errors.ErrUnauthorized
	}
	
	// Validate post
	if post.Content == "" && post.Image == "" {
		return 0, errors.ErrInvalidInput
	}
	
	// Validate privacy level
	if post.PrivacyLevel == "" {
		post.PrivacyLevel = posts.PrivacyPublic // Default to public
	} else if post.PrivacyLevel != posts.PrivacyPublic && 
		post.PrivacyLevel != posts.PrivacyAlmostPrivate && 
		post.PrivacyLevel != posts.PrivacyPrivate {
		return 0, errors.ErrInvalidInput
	}
	
	// For group posts, override some privacy settings
	// All group posts are private to the group by default
	post.PrivacyLevel = posts.PrivacyPrivate
	
	// Insert the post with a reference to the group
	query := `INSERT INTO posts (user_id, group_id, content, image, privacy_level, allowed_followers, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	result, err := s.db.Exec(
		query,
		userID,
		groupID,
		post.Content,
		post.Image,
		post.PrivacyLevel,
		post.AllowedFollowers,
		now,
		now,
	)
	if err != nil {
		log.Print("Error creating group post:", err)
		return 0, errors.ErrInternalServer
	}
	
	postID, err := result.LastInsertId()
	if err != nil {
		log.Print("Error getting last insert ID:", err)
		return 0, errors.ErrInternalServer
	}
	
	return postID, nil
}

// GetGroupPosts gets all posts in a group
func (s *SQLiteStore) GetGroupPosts(groupID int64) ([]*posts.Post, error) {
	// First, check if the group exists
	_, err := s.GetGroup(groupID)
	if err != nil {
		return nil, err
	}
	
	// Query posts
	query := `SELECT p.id, p.user_id, p.content, p.image, p.privacy_level, p.allowed_followers,
			  p.created_at, p.updated_at,
              (SELECT COUNT(*) FROM likes WHERE post_id = p.id) as likes_count
              FROM posts p
              WHERE p.group_id = ?
              ORDER BY p.created_at DESC`
	
	rows, err := s.db.Query(query, groupID)
	if err != nil {
		log.Print("Error querying group posts:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()
	
	var result []*posts.Post
	for rows.Next() {
		p := &posts.Post{}
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Content,
			&p.Image,
			&p.PrivacyLevel,
			&p.AllowedFollowers,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Likes,
		)
		if err != nil {
			log.Print("Error scanning group post row:", err)
			return nil, errors.ErrInternalServer
		}
		
		// Get comments for this post
		comments, err := s.getPostComments(p.ID)
		if err != nil {
			log.Print("Error getting comments:", err)
			return nil, errors.ErrInternalServer
		}
		p.Comments = comments
		
		result = append(result, p)
	}
	
	if err = rows.Err(); err != nil {
		log.Print("Error iterating group post rows:", err)
		return nil, errors.ErrInternalServer
	}
	
	return result, nil
}