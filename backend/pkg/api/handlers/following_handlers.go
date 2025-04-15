package handlers

import (
	"net/http"
	"strconv"

	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/domain/following"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/logger"
	"socialNetwork/pkg/storage"
)

// FollowingHandler manages following-related handlers
type FollowingHandler struct {
	store  storage.Storage
	logger logger.Logger
}

// NewFollowingHandler creates a new following handler
func NewFollowingHandler(store storage.Storage) *FollowingHandler {
	return &FollowingHandler{
		store:  store,
		logger: logger.GetLogger("api/handlers/following"),
	}
}

// HandleFollowing handles creating a follow request
func (h *FollowingHandler) HandleFollowing(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get the ID of the user to follow from the URL
	idStr := r.PathValue("id")
	followedID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid user ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get the follower's ID from the session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	followerID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Prevent following oneself
	if followerID == followedID {
		log.Warn("User attempted to follow themselves", "user_id", followerID)
		return errors.ErrInvalidInput
	}

	// Create the follow request
	followReq := following.FollowRequest{
		FollowerID: followerID,
		FollowedID: followedID,
	}

	err = h.store.CreateFollowRequest(followReq)
	if err != nil {
		// Check for already following error
		if err == errors.ErrAlreadyExists {
			log.Info("Already following or request pending", "follower", followerID, "followed", followedID)
			return utils.WriteJson(w, http.StatusOK, map[string]string{
				"message": "Already following or request pending",
			})
		}
		log.Error("Failed to create follow request", "error", err)
		return err
	}

	log.Info("Follow request created", "follower", followerID, "followed", followedID)

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Follow request sent successfully",
	})
}

// HandleFollowingRequests handles retrieving pending follow requests
func (h *FollowingHandler) HandleFollowingRequests(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get the user's ID from the session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	userID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Get pending follow requests
	requests, err := h.store.GetFollowRequests(userID)
	if err != nil {
		log.Error("Failed to get follow requests", "error", err)
		return err
	}

	log.Info("Fetched follow requests", "user_id", userID, "count", len(requests))

	return utils.WriteJson(w, http.StatusOK, requests)
}

// HandleFollowingAccept handles accepting a follow request
func (h *FollowingHandler) HandleFollowingAccept(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get the follower's ID from the URL
	idStr := r.PathValue("id")
	followerID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid user ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get the followed user's ID from the session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	followedID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Accept the follow request
	err = h.store.AcceptFollowRequest(followerID, followedID)
	if err != nil {
		if err == errors.ErrNotFound {
			log.Warn("Follow request not found", "follower", followerID, "followed", followedID)
			return errors.ErrNotFound
		}
		log.Error("Failed to accept follow request", "error", err)
		return err
	}

	log.Info("Follow request accepted", "follower", followerID, "followed", followedID)

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Follow request accepted",
	})
}

// HandleFollowingReject handles rejecting a follow request
func (h *FollowingHandler) HandleFollowingReject(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get the follower's ID from the URL
	idStr := r.PathValue("id")
	followerID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid user ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get the followed user's ID from the session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	followedID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Delete the follow request
	err = h.store.DeleteFollowRequest(followerID, followedID)
	if err != nil {
		if err == errors.ErrNotFound {
			log.Warn("Follow request not found", "follower", followerID, "followed", followedID)
			return errors.ErrNotFound
		}
		log.Error("Failed to reject follow request", "error", err)
		return err
	}

	log.Info("Follow request rejected", "follower", followerID, "followed", followedID)

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Follow request rejected",
	})
}