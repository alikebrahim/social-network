package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/domain/notifications"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/logger"
	"socialNetwork/pkg/storage"
)

// NotificationHandler handles notification-related endpoints
type NotificationHandler struct {
	store storage.Storage
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(store storage.Storage) *NotificationHandler {
	return &NotificationHandler{
		store: store,
	}
}

// HandleGetNotifications gets a user's notifications
func (h *NotificationHandler) HandleGetNotifications(w http.ResponseWriter, r *http.Request) error {
	// Get user ID from context after JWT auth
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		return err
	}

	// Get limit and offset from query params, default to 20 and 0
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit < 1 {
			limit = 20
		}
	}

	if offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil || offset < 0 {
			offset = 0
		}
	}

	// Get notifications
	notificationsList, err := h.store.GetNotifications(userID, limit, offset)
	if err != nil {
		log := logger.FromRequest(r)
		log.Error("Failed to get notifications", "error", err)
		return err
	}

	return utils.WriteJson(w, http.StatusOK, notificationsList)
}

// HandleMarkNotificationRead marks a notification as read
func (h *NotificationHandler) HandleMarkNotificationRead(w http.ResponseWriter, r *http.Request) error {
	// Get user ID from context after JWT auth
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		return err
	}

	// Get notification ID from URL path
	notificationIDStr := utils.GetParam(r, "id")
	notificationID, err := strconv.ParseInt(notificationIDStr, 10, 64)
	if err != nil {
		return errors.ErrBadRequest
	}

	// Mark as read
	err = h.store.MarkNotificationAsRead(notificationID, userID)
	if err != nil {
		log := logger.FromRequest(r)
		log.Error("Failed to mark notification as read", "error", err)
		return err
	}

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Notification marked as read",
	})
}

// HandleMarkAllNotificationsRead marks all notifications as read
func (h *NotificationHandler) HandleMarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) error {
	// Get user ID from context after JWT auth
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		return err
	}

	// Mark all as read
	err = h.store.MarkAllNotificationsAsRead(userID)
	if err != nil {
		log := logger.FromRequest(r)
		log.Error("Failed to mark all notifications as read", "error", err)
		return err
	}

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "All notifications marked as read",
	})
}

// HandleGetUnreadCount gets the count of unread notifications
func (h *NotificationHandler) HandleGetUnreadCount(w http.ResponseWriter, r *http.Request) error {
	// Get user ID from context after JWT auth
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		return err
	}

	// Get unread count
	count, err := h.store.GetUnreadNotificationCount(userID)
	if err != nil {
		log := logger.FromRequest(r)
		log.Error("Failed to get unread notification count", "error", err)
		return err
	}

	return utils.WriteJson(w, http.StatusOK, map[string]int{
		"count": count,
	})
}