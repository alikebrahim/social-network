package services

import (
	"backend/pkg/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SendPrivateMessage sends a private message from one user to another.
func SendPrivateMessage(db *sql.DB, senderID, receiverID, message string) (string, error) {
	// Validate that the sender and receiver are not the same user
	if senderID == receiverID {
		return "", fmt.Errorf("sender and receiver cannot be the same user")
	}

	// Generate a new message ID and insert the private message into the database
	messageID := uuid.New().String()
	query := `INSERT INTO private_chat_messages (id, sender_id, receiver_id, message, created_at)
			  VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, messageID, senderID, receiverID, message, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to send private message: %w", err)
	}

	return messageID, nil
}

// GetPrivateMessages retrieves all private messages between two users.
func GetPrivateMessages(db *sql.DB, user1ID, user2ID string) ([]models.PrivateChatMessage, error) {
	// Retrieve messages between the two users
	query := `SELECT id, sender_id, receiver_id, message, created_at 
			  FROM private_chat_messages 
			  WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) 
			  ORDER BY created_at ASC`
	rows, err := db.Query(query, user1ID, user2ID, user2ID, user1ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private messages: %w", err)
	}
	defer rows.Close()

	var messages []models.PrivateChatMessage
	for rows.Next() {
		var msg models.PrivateChatMessage
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Message, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
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

// GetGroupMessages retrieves all messages for a group.
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

// CreateGroupChat creates a new group chat and adds the creator as the first member.
func CreateGroupChat(db *sql.DB, userID, groupTitle, groupDescription string, isPublic bool) (string, error) {
	// Generate a new group ID and insert the new group into the database
	groupID := uuid.New().String()
	query := `INSERT INTO group_chats (id, group_id, created_at) VALUES (?, ?, ?)`
	_, err := db.Exec(query, groupID, groupTitle, groupDescription, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to create group chat: %w", err)
	}

	// Add the creator as the first member and as an admin
	query = `INSERT INTO group_members (group_id, user_id, joined_at, is_admin)
			  VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, groupID, userID, time.Now(), true)
	if err != nil {
		return "", fmt.Errorf("failed to add user to group chat: %w", err)
	}

	return groupID, nil
}

// AddGroupMember adds a user to an existing group.
func AddGroupMember(db *sql.DB, groupID, userID string) (string, error) {
	// Ensure the user isn't already a member
	var existingMemberID string
	query := `SELECT user_id FROM group_members WHERE group_id = ? AND user_id = ?`
	err := db.QueryRow(query, groupID, userID).Scan(&existingMemberID)
	if err == nil {
		return "", fmt.Errorf("user is already a member of the group")
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to check group membership: %w", err)
	}

	// Add the user to the group
	query = `INSERT INTO group_members (group_id, user_id, joined_at) VALUES (?, ?, ?)`
	_, err = db.Exec(query, groupID, userID, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to add user to group: %w", err)
	}

	return userID, nil
}
