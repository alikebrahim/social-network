package sqlite

import (
	"log"
	"time"

	"socialNetwork/pkg/domain/notifications"
	"socialNetwork/pkg/errors"
)

// Create a new notification
func (s *SQLiteStore) CreateNotification(notification *notifications.Notification) error {
	query := `INSERT INTO notifications 
		(user_id, type, content, related_id, sender_id, created_at) 
		VALUES (?, ?, ?, ?, ?, ?)`

	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = time.Now()
	}

	_, err := s.db.Exec(
		query,
		notification.UserID,
		notification.Type,
		notification.Content,
		notification.RelatedID,
		notification.SenderID,
		notification.CreatedAt,
	)

	if err != nil {
		log.Print("Error creating notification:", err)
		return errors.ErrInternalServer
	}

	return nil
}

// Get notifications for user
func (s *SQLiteStore) GetNotifications(userID int64, limit, offset int) (notifications.NotificationsList, error) {
	query := `
		SELECT n.id, n.user_id, n.type, n.content, n.is_read, n.related_id, n.created_at, 
		       n.sender_id, u.first_name || ' ' || u.last_name as sender_name, u.avatar as sender_image
		FROM notifications n
		LEFT JOIN users u ON n.sender_id = u.id
		WHERE n.user_id = ?
		ORDER BY n.created_at DESC
		LIMIT ? OFFSET ?
	`

	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = ?`

	rows, err := s.db.Query(query, userID, limit, offset)
	if err != nil {
		log.Print("Error querying notifications:", err)
		return notifications.NotificationsList{}, errors.ErrInternalServer
	}
	defer rows.Close()

	var notifs []notifications.Notification
	for rows.Next() {
		var notif notifications.Notification
		var createdAtStr string

		err := rows.Scan(
			&notif.ID, 
			&notif.UserID,
			&notif.Type,
			&notif.Content,
			&notif.IsRead,
			&notif.RelatedID,
			&createdAtStr,
			&notif.SenderID,
			&notif.SenderName,
			&notif.SenderImage,
		)
		if err != nil {
			log.Print("Error scanning notification row:", err)
			return notifications.NotificationsList{}, errors.ErrInternalServer
		}

		// Parse created_at string to time.Time
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			log.Print("Error parsing notification timestamp:", err)
			notif.CreatedAt = time.Time{} // Zero time if parsing fails
		} else {
			notif.CreatedAt = createdAt
		}

		notifs = append(notifs, notif)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error iterating notification rows:", err)
		return notifications.NotificationsList{}, errors.ErrInternalServer
	}

	// Get total count
	var count int
	err = s.db.QueryRow(countQuery, userID).Scan(&count)
	if err != nil {
		log.Print("Error counting notifications:", err)
		return notifications.NotificationsList{}, errors.ErrInternalServer
	}

	return notifications.NotificationsList{
		Notifications: notifs,
		Count:         count,
	}, nil
}

// Mark notification as read
func (s *SQLiteStore) MarkNotificationAsRead(notificationID, userID int64) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = ? AND user_id = ?`

	result, err := s.db.Exec(query, notificationID, userID)
	if err != nil {
		log.Print("Error marking notification as read:", err)
		return errors.ErrInternalServer
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Print("Error getting rows affected:", err)
		return errors.ErrInternalServer
	}

	if rows == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Mark all notifications as read
func (s *SQLiteStore) MarkAllNotificationsAsRead(userID int64) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE user_id = ?`

	_, err := s.db.Exec(query, userID)
	if err != nil {
		log.Print("Error marking all notifications as read:", err)
		return errors.ErrInternalServer
	}

	return nil
}

// Get unread notification count
func (s *SQLiteStore) GetUnreadNotificationCount(userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = FALSE`

	var count int
	err := s.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		log.Print("Error counting unread notifications:", err)
		return 0, errors.ErrInternalServer
	}

	return count, nil
}