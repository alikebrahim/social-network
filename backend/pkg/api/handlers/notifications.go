package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Notification struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// GetNotificationsHandler retrieves notifications for a user
func GetNotificationsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id").(string) // Assuming user ID is stored in the context

		query := `SELECT id, user_id, message, is_read, created_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC`
		rows, err := db.Query(query, userID)
		if err != nil {
			http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var notifications []Notification
		for rows.Next() {
			var notification Notification
			err := rows.Scan(&notification.ID, &notification.UserID, &notification.Message, &notification.IsRead, &notification.CreatedAt)
			if err != nil {
				http.Error(w, "Failed to parse notifications", http.StatusInternalServerError)
				return
			}
			notifications = append(notifications, notification)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(notifications)
	}
}

// MarkNotificationAsReadHandler marks a notification as read
func MarkNotificationAsReadHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notificationID := r.URL.Query().Get("id")
		if notificationID == "" {
			http.Error(w, "Notification ID is required", http.StatusBadRequest)
			return
		}

		query := `UPDATE notifications SET is_read = 1 WHERE id = ?`
		result, err := db.Exec(query, notificationID)
		if err != nil {
			http.Error(w, "Failed to mark notification as read", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Notification marked as read"})
	}
}

// CreateNotification creates a notification for a user
func CreateNotification(db *sql.DB, userID string, message string) error {
	query := `INSERT INTO notifications (user_id, message, is_read, created_at) VALUES (?, ?, 0, ?)`
	_, err := db.Exec(query, userID, message, time.Now())
	return err
}
