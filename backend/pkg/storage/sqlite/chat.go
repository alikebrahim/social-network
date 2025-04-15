package sqlite

import (
	"log"
	"time"

	"socialNetwork/pkg/domain/chat"
	"socialNetwork/pkg/errors"
)

// SaveChat saves a chat message between users
func (s *SQLiteStore) SaveChat(chat *chat.Chat) (int64, error) {
	// Insert the chat message
	query := `INSERT INTO chats (sender_id, receiver_id, content, created_at) 
              VALUES (?, ?, ?, ?)`

	if chat.CreatedAt.IsZero() {
		chat.CreatedAt = time.Now()
	}

	result, err := s.db.Exec(query, chat.SenderID, chat.ReceiverID, chat.Content, chat.CreatedAt)
	if err != nil {
		log.Print("Error saving chat message:", err)
		return 0, errors.ErrInternalServer
	}

	// Get the ID of the newly created message
	chatID, err := result.LastInsertId()
	if err != nil {
		log.Print("Error getting last insert ID:", err)
		return 0, errors.ErrInternalServer
	}

	return chatID, nil
}

// GetUserChats gets all chats for a user (conversations)
func (s *SQLiteStore) GetUserChats(userID int64) ([]*chat.Chat, error) {
	// This query gets the most recent message from each conversation
	query := `WITH recent_chats AS (
                SELECT 
                    c.id,
                    c.sender_id,
                    c.receiver_id,
                    c.content,
                    c.created_at,
                    ROW_NUMBER() OVER (
                        PARTITION BY 
                            CASE 
                                WHEN c.sender_id = ? THEN c.receiver_id 
                                ELSE c.sender_id 
                            END 
                        ORDER BY c.created_at DESC
                    ) as rn
                FROM chats c
                WHERE c.sender_id = ? OR c.receiver_id = ?
            )
            SELECT id, sender_id, receiver_id, content, created_at
            FROM recent_chats
            WHERE rn = 1
            ORDER BY created_at DESC`

	rows, err := s.db.Query(query, userID, userID, userID)
	if err != nil {
		log.Print("Error querying user chats:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	var chatsList []*chat.Chat
	for rows.Next() {
		var msg chat.Chat
		err := rows.Scan(
			&msg.ID,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.Content,
			&msg.CreatedAt,
		)
		if err != nil {
			log.Print("Error scanning chat row:", err)
			return nil, errors.ErrInternalServer
		}
		chatsList = append(chatsList, &msg)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating chat rows:", err)
		return nil, errors.ErrInternalServer
	}

	return chatsList, nil
}

// GetChatHistory gets the chat history between two users
func (s *SQLiteStore) GetChatHistory(userID1, userID2 int64) ([]*chat.Chat, error) {
	query := `SELECT id, sender_id, receiver_id, content, created_at
              FROM chats
              WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
              ORDER BY created_at`

	rows, err := s.db.Query(query, userID1, userID2, userID2, userID1)
	if err != nil {
		log.Print("Error querying chat history:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	var chatsList []*chat.Chat
	for rows.Next() {
		var msg chat.Chat
		err := rows.Scan(
			&msg.ID,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.Content,
			&msg.CreatedAt,
		)
		if err != nil {
			log.Print("Error scanning chat history row:", err)
			return nil, errors.ErrInternalServer
		}
		chatsList = append(chatsList, &msg)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating chat history rows:", err)
		return nil, errors.ErrInternalServer
	}

	return chatsList, nil
}

// SaveGroupChat saves a chat message in a group
func (s *SQLiteStore) SaveGroupChat(chat *chat.GroupChat) (int64, error) {
	// Check if the sender is a member of the group
	isMember, err := s.IsGroupMember(chat.GroupID, chat.SenderID)
	if err != nil {
		return 0, err
	}
	if !isMember {
		return 0, errors.ErrUnauthorized
	}

	// Insert the group chat message
	query := `INSERT INTO group_chats (group_id, sender_id, content, created_at) 
              VALUES (?, ?, ?, ?)`

	if chat.CreatedAt.IsZero() {
		chat.CreatedAt = time.Now()
	}

	result, err := s.db.Exec(query, chat.GroupID, chat.SenderID, chat.Content, chat.CreatedAt)
	if err != nil {
		log.Print("Error saving group chat message:", err)
		return 0, errors.ErrInternalServer
	}

	// Get the ID of the newly created message
	chatID, err := result.LastInsertId()
	if err != nil {
		log.Print("Error getting last insert ID:", err)
		return 0, errors.ErrInternalServer
	}

	return chatID, nil
}

// GetGroupChatHistory gets the chat history for a group
func (s *SQLiteStore) GetGroupChatHistory(groupID int64) ([]*chat.GroupChat, error) {
	query := `SELECT id, group_id, sender_id, content, created_at
              FROM group_chats
              WHERE group_id = ?
              ORDER BY created_at`

	rows, err := s.db.Query(query, groupID)
	if err != nil {
		log.Print("Error querying group chat history:", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	var chatsList []*chat.GroupChat
	for rows.Next() {
		var msg chat.GroupChat
		err := rows.Scan(
			&msg.ID,
			&msg.GroupID,
			&msg.SenderID,
			&msg.Content,
			&msg.CreatedAt,
		)
		if err != nil {
			log.Print("Error scanning group chat row:", err)
			return nil, errors.ErrInternalServer
		}
		chatsList = append(chatsList, &msg)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating group chat rows:", err)
		return nil, errors.ErrInternalServer
	}

	return chatsList, nil
}