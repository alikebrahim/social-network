package main

import (
	"time"
)

// SaveChat saves a chat message
func (s *SQLiteStore) SaveChat(chat *Chat) (int64, error) {
	query := `INSERT INTO chats (sender_id, receiver_id, content, image, created_at) 
              VALUES (?, ?, ?, ?, ?)`

	result, err := s.db.Exec(
		query,
		chat.SenderID,
		chat.ReceiverID,
		chat.Content,
		chat.Image,
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	chatID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

// GetUserChats retrieves a user's chats (latest message from each conversation)
func (s *SQLiteStore) GetUserChats(userID int64) ([]*Chat, error) {
	// This query gets the latest message from each conversation
	query := `
        WITH LatestMessages AS (
            SELECT c1.*, 
                ROW_NUMBER() OVER (
                    PARTITION BY 
                        CASE 
                            WHEN c1.sender_id = ? THEN c1.receiver_id 
                            ELSE c1.sender_id 
                        END
                    ORDER BY c1.created_at DESC
                ) as rn
            FROM chats c1
            WHERE c1.sender_id = ? OR c1.receiver_id = ?
        )
        SELECT id, sender_id, receiver_id, content, image, created_at 
        FROM LatestMessages
        WHERE rn = 1
        ORDER BY created_at DESC
    `

	rows, err := s.db.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*Chat
	for rows.Next() {
		var chat Chat
		err := rows.Scan(
			&chat.ID,
			&chat.SenderID,
			&chat.ReceiverID,
			&chat.Content,
			&chat.Image,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, &chat)
	}

	return chats, nil
}

// GetChatHistory retrieves the chat history between two users
func (s *SQLiteStore) GetChatHistory(userID1, userID2 int64) ([]*Chat, error) {
	query := `SELECT id, sender_id, receiver_id, content, image, created_at 
              FROM chats 
              WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) 
              ORDER BY created_at ASC`

	rows, err := s.db.Query(query, userID1, userID2, userID2, userID1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*Chat
	for rows.Next() {
		var chat Chat
		err := rows.Scan(
			&chat.ID,
			&chat.SenderID,
			&chat.ReceiverID,
			&chat.Content,
			&chat.Image,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, &chat)
	}

	return chats, nil
}

// SaveGroupChat saves a message in a group chat
func (s *SQLiteStore) SaveGroupChat(chat *GroupChat) (int64, error) {
	// Check if user is a member of the group
	isMember, err := s.IsGroupMember(chat.GroupID, chat.SenderID)
	if err != nil {
		return 0, err
	}
	if !isMember {
		return 0, ErrUnauthorized
	}

	query := `INSERT INTO group_chats (group_id, sender_id, content, image, created_at) 
              VALUES (?, ?, ?, ?, ?)`

	result, err := s.db.Exec(
		query,
		chat.GroupID,
		chat.SenderID,
		chat.Content,
		chat.Image,
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	chatID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

// GetGroupChatHistory retrieves the chat history for a group
func (s *SQLiteStore) GetGroupChatHistory(groupID int64) ([]*GroupChat, error) {
	query := `SELECT id, group_id, sender_id, content, image, created_at 
              FROM group_chats 
              WHERE group_id = ? 
              ORDER BY created_at ASC`

	rows, err := s.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*GroupChat
	for rows.Next() {
		var chat GroupChat
		err := rows.Scan(
			&chat.ID,
			&chat.GroupID,
			&chat.SenderID,
			&chat.Content,
			&chat.Image,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, &chat)
	}

	return chats, nil
}
