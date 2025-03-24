package main

import (
	"database/sql"
	"time"
)

// GetGroupMembers retrieves all members of a group
func (s *SQLiteStore) GetGroupMembers(groupID int64) ([]*GroupMember, error) {
	query := `SELECT id, group_id, user_id, status, created_at 
              FROM group_members 
              WHERE group_id = ?`

	rows, err := s.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*GroupMember
	for rows.Next() {
		var member GroupMember
		err := rows.Scan(
			&member.ID,
			&member.GroupID,
			&member.UserID,
			&member.Status,
			&member.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, &member)
	}

	return members, nil
}

// GetGroupInvites retrieves all groups that a user has been invited to
func (s *SQLiteStore) GetGroupInvites(userID int64) ([]*Group, error) {
	query := `SELECT g.id, g.creator_id, g.title, g.description, g.created_at, g.updated_at 
              FROM groups g 
              JOIN group_members gm ON g.id = gm.group_id 
              WHERE gm.user_id = ? AND gm.status = 'pending'`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*Group
	for rows.Next() {
		var group Group
		err := rows.Scan(
			&group.ID,
			&group.CreatorID,
			&group.Title,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}

// AcceptGroupInvite accepts a group invitation
func (s *SQLiteStore) AcceptGroupInvite(groupID, userID int64) error {
	query := `UPDATE group_members 
              SET status = 'accepted' 
              WHERE group_id = ? AND user_id = ? AND status = 'pending'`

	result, err := s.db.Exec(query, groupID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// RejectGroupInvite rejects a group invitation
func (s *SQLiteStore) RejectGroupInvite(groupID, userID int64) error {
	query := `DELETE FROM group_members 
              WHERE group_id = ? AND user_id = ? AND status = 'pending'`

	result, err := s.db.Exec(query, groupID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// RequestJoinGroup creates a request to join a group
func (s *SQLiteStore) RequestJoinGroup(groupID, userID int64) error {
	// Check if user already has a pending invitation or is a member
	query := `SELECT COUNT(*) FROM group_members 
              WHERE group_id = ? AND user_id = ?`

	var count int
	err := s.db.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrAlreadyExists
	}

	// Create join request
	insertQuery := `INSERT INTO group_members (group_id, user_id, status, created_at) 
                   VALUES (?, ?, ?, ?)`
	_, err = s.db.Exec(insertQuery, groupID, userID, "pending", time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetGroupJoinRequests retrieves all pending join requests for a group
func (s *SQLiteStore) GetGroupJoinRequests(groupID int64) ([]*GroupMember, error) {
	query := `SELECT id, group_id, user_id, status, created_at 
              FROM group_members 
              WHERE group_id = ? AND status = 'pending'`

	rows, err := s.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*GroupMember
	for rows.Next() {
		var req GroupMember
		err := rows.Scan(
			&req.ID,
			&req.GroupID,
			&req.UserID,
			&req.Status,
			&req.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &req)
	}

	return requests, nil
}

// AcceptGroupJoinRequest accepts a request to join a group
func (s *SQLiteStore) AcceptGroupJoinRequest(groupID, userID int64) error {
	query := `UPDATE group_members 
              SET status = 'accepted' 
              WHERE group_id = ? AND user_id = ? AND status = 'pending'`

	result, err := s.db.Exec(query, groupID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// RejectGroupJoinRequest rejects a request to join a group
func (s *SQLiteStore) RejectGroupJoinRequest(groupID, userID int64) error {
	query := `DELETE FROM group_members 
              WHERE group_id = ? AND user_id = ? AND status = 'pending'`

	result, err := s.db.Exec(query, groupID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// CreateGroupEvent creates a new event in a group
func (s *SQLiteStore) CreateGroupEvent(event *GroupEvent) (int64, error) {
	query := `INSERT INTO events (group_id, creator_id, title, description, event_time, created_at) 
              VALUES (?, ?, ?, ?, ?, ?)`

	result, err := s.db.Exec(
		query,
		event.GroupID,
		event.CreatorID,
		event.Title,
		event.Description,
		event.EventTime,
		event.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return eventID, nil
}

// GetGroupEvents retrieves all events in a group
func (s *SQLiteStore) GetGroupEvents(groupID int64) ([]*GroupEvent, error) {
	query := `SELECT id, group_id, creator_id, title, description, event_time, created_at 
              FROM events 
              WHERE group_id = ? 
              ORDER BY event_time DESC`

	rows, err := s.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*GroupEvent
	for rows.Next() {
		var event GroupEvent
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
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

// RespondToEvent records a user's response to an event
func (s *SQLiteStore) RespondToEvent(eventID, userID int64, response string) error {
	// Check if user has already responded
	checkQuery := `SELECT id FROM event_responses 
                 WHERE event_id = ? AND user_id = ?`

	var responseID int64
	err := s.db.QueryRow(checkQuery, eventID, userID).Scan(&responseID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// If response exists, update it
	if err == nil {
		updateQuery := `UPDATE event_responses 
                     SET response = ?, created_at = ? 
                     WHERE id = ?`
		_, err = s.db.Exec(updateQuery, response, time.Now(), responseID)
		return err
	}

	// Otherwise, create a new response
	insertQuery := `INSERT INTO event_responses (event_id, user_id, response, created_at) 
                   VALUES (?, ?, ?, ?)`
	_, err = s.db.Exec(insertQuery, eventID, userID, response, time.Now())
	return err
}

// GetEventResponses retrieves all responses to an event
func (s *SQLiteStore) GetEventResponses(eventID int64) ([]*EventResponse, error) {
	query := `SELECT id, event_id, user_id, response, created_at 
              FROM event_responses 
              WHERE event_id = ?`

	rows, err := s.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []*EventResponse
	for rows.Next() {
		var resp EventResponse
		err := rows.Scan(
			&resp.ID,
			&resp.EventID,
			&resp.UserID,
			&resp.Response,
			&resp.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		responses = append(responses, &resp)
	}

	return responses, nil
}

// CreateGroupPost creates a new post in a group
func (s *SQLiteStore) CreateGroupPost(groupID, userID int64, post *Post) (int64, error) {
	// Check if user is a member of the group
	isMember, err := s.IsGroupMember(groupID, userID)
	if err != nil {
		return 0, err
	}
	if !isMember {
		return 0, ErrUnauthorized
	}

	// Set post attributes
	post.UserID = userID
	post.GroupID = groupID
	post.CreatedAt = time.Now()

	// Insert post
	query := `INSERT INTO posts (user_id, group_id, content, image, created_at) 
              VALUES (?, ?, ?, ?, ?)`

	result, err := s.db.Exec(
		query,
		post.UserID,
		post.GroupID,
		post.Content,
		post.Image,
		post.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

// GetGroupPosts retrieves all posts in a group
func (s *SQLiteStore) GetGroupPosts(groupID int64) ([]*Post, error) {
	query := `SELECT id, user_id, group_id, content, image, created_at 
              FROM posts 
              WHERE group_id = ? 
              ORDER BY created_at DESC`

	rows, err := s.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.GroupID,
			&post.Content,
			&post.Image,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}
