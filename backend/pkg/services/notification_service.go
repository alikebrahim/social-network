package services

import (
	"backend/pkg/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreateGroup creates a new group and adds the creator as the first member.
func CreateGroup(db *sql.DB, userID, title, description string, isPublic bool) (string, error) {
	// Generate a new group ID
	groupID := uuid.New().String()

	// Insert the new group into the database
	query := `INSERT INTO groups (id, title, description, owner_id, created_at, is_public)
			  VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, groupID, title, description, userID, time.Now(), isPublic)
	if err != nil {
		return "", fmt.Errorf("failed to create group: %w", err)
	}

	// Add the creator as the first member and admin
	query = `INSERT INTO group_members (group_id, user_id, joined_at, is_admin)
			  VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, groupID, userID, time.Now(), true)
	if err != nil {
		return "", fmt.Errorf("failed to add creator as first member: %w", err)
	}

	return groupID, nil
}

// AddGroupMember adds a user to an existing group as a member.
func AddGroupMember(db *sql.DB, groupID, userID string) error {
	// Ensure the user isn't already a member of the group
	var existingMemberID string
	query := `SELECT user_id FROM group_members WHERE group_id = ? AND user_id = ?`
	err := db.QueryRow(query, groupID, userID).Scan(&existingMemberID)
	if err == nil {
		return fmt.Errorf("user is already a member of the group")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check group membership: %w", err)
	}

	// Add the user to the group
	query = `INSERT INTO group_members (group_id, user_id, joined_at, is_admin)
			  VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, groupID, userID, time.Now(), false)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}

	return nil
}

// RemoveGroupMember removes a user from a group.
func RemoveGroupMember(db *sql.DB, groupID, userID string) error {
	// Ensure the user is part of the group
	var existingMemberID string
	query := `SELECT user_id FROM group_members WHERE group_id = ? AND user_id = ?`
	err := db.QueryRow(query, groupID, userID).Scan(&existingMemberID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user is not a member of the group")
	}
	if err != nil {
		return fmt.Errorf("failed to check group membership: %w", err)
	}

	// Remove the user from the group
	query = `DELETE FROM group_members WHERE group_id = ? AND user_id = ?`
	_, err = db.Exec(query, groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}

	return nil
}

// SendGroupMessage sends a message to a group chat.
func SendGroupMessage(db *sql.DB, groupID, senderID, message string) (string, error) {
	// Generate a new message ID and insert the group message into the database
	messageID := uuid.New().String()
	query := `INSERT INTO group_chat_messages (id, group_id, sender_id, message, created_at)
			  VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, messageID, groupID, senderID, message, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to send group message: %w", err)
	}

	return messageID, nil
}

// GetGroupMessages retrieves all messages for a group chat.
func GetGroupMessages(db *sql.DB, groupID string) ([]models.GroupChatMessage, error) {
	// Retrieve messages for the group chat
	query := `SELECT id, group_id, sender_id, message, created_at 
			  FROM group_chat_messages 
			  WHERE group_id = ? 
			  ORDER BY created_at ASC`
	rows, err := db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve group messages: %w", err)
	}
	defer rows.Close()

	var messages []models.GroupChatMessage
	for rows.Next() {
		var msg models.GroupChatMessage
		err := rows.Scan(&msg.ID, &msg.GroupID, &msg.SenderID, &msg.Message, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse group message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// CreateGroupEvent creates an event for a group.
func CreateGroupEvent(db *sql.DB, groupID, title, description, dayTime string) (string, error) {
	// Generate a new event ID and insert the event into the database
	eventID := uuid.New().String()
	query := `INSERT INTO group_events (id, group_id, title, description, day_time, created_at)
			  VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, eventID, groupID, title, description, dayTime, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to create group event: %w", err)
	}

	return eventID, nil
}

// GetGroupEvents retrieves all events for a group.
func GetGroupEvents(db *sql.DB, groupID string) ([]models.GroupEvent, error) {
	// Retrieve events for the group
	query := `SELECT id, group_id, title, description, day_time, created_at 
			  FROM group_events 
			  WHERE group_id = ? 
			  ORDER BY created_at ASC`
	rows, err := db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve group events: %w", err)
	}
	defer rows.Close()

	var events []models.GroupEvent
	for rows.Next() {
		var event models.GroupEvent
		err := rows.Scan(&event.ID, &event.GroupID, &event.Title, &event.Description, &event.DayTime, &event.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse group event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}
