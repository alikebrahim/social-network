package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/domain/profile"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/logger"
	"socialNetwork/pkg/storage"
)

// ProfileHandler manages profile-related handlers
type ProfileHandler struct {
	store  storage.Storage
	logger logger.Logger
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(store storage.Storage) *ProfileHandler {
	return &ProfileHandler{
		store:  store,
		logger: logger.GetLogger("api/handlers/profile"),
	}
}

// HandleGetProfile handles retrieving a user's profile
func (h *ProfileHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)
	
	// Create a new profile instance
	var p profile.Profile
	
	// Get profile ID from request
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid profile ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}
	p.UserID = id
	
	// Get the requesting user's ID from the session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}
	
	requestID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}
	
	// Check if the requester is following the profile
	isFollowing, err := h.store.IsFollowing(requestID, p.UserID)
	if err != nil {
		log.Error("Error checking follow status", "error", err)
		return errors.ErrInternalServer
	}
	
	// Get profile data
	profileData, err := h.store.GetProfileData(&p)
	if err != nil {
		log.Error("Failed to get profile data", "error", err)
		return err
	}
	
	// If profile is private and requester is not following, hide sensitive information
	if profileData.ProfileType == "private" && !isFollowing && requestID != p.UserID {
		log.Info("Private profile accessed by non-follower", "profile_id", p.UserID, "requester_id", requestID)
		// Create a limited view of the profile
		limitedProfile := &profile.Profile{
			UserID:       profileData.UserID,
			FirstName:    profileData.FirstName,
			LastName:     profileData.LastName,
			Nickname:     profileData.Nickname,
			Avatar:       profileData.Avatar,
			ProfileType:  profileData.ProfileType,
			Followers:    profileData.Followers,
			Following:    profileData.Following,
			PostCount:    profileData.PostCount,
			CreatedAt:    profileData.CreatedAt,
			LastActivity: profileData.LastActivity,
		}
		// Don't include DateOfBirth or AboutMe for private profiles
		return utils.WriteJson(w, http.StatusOK, limitedProfile)
	}
	
	return utils.WriteJson(w, http.StatusOK, profileData)
}

// HandleSetProfilePrivacy handles updating a user's profile privacy settings
func (h *ProfileHandler) HandleSetProfilePrivacy(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)
	
	// Get the requesting user's ID from the session
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
	
	// Parse the request body
	var privacyReq profile.ProfilePrivacyRequest
	if err := json.NewDecoder(r.Body).Decode(&privacyReq); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return errors.ErrBadRequest
	}
	
	// Validate privacy type
	if privacyReq.ProfileType != "public" && privacyReq.ProfileType != "private" {
		log.Warn("Invalid profile type", "type", privacyReq.ProfileType)
		return errors.ErrInvalidInput
	}
	
	// Create profile object with user ID
	p := profile.Profile{
		UserID: userID,
	}
	
	// Get the current profile to check if privacy setting is already set
	currentProfile, err := h.store.GetProfileData(&p)
	if err != nil {
		log.Error("Failed to get current profile data", "error", err)
		return err
	}
	
	// Check if privacy is already set to the requested value
	if currentProfile.ProfileType == privacyReq.ProfileType {
		log.Info("Profile privacy already set", "type", privacyReq.ProfileType)
		return utils.WriteJson(w, http.StatusOK, map[string]string{
			"message": "Privacy setting unchanged",
		})
	}
	
	// Update the profile privacy
	profileToUpdate := profile.Profile{
		UserID: userID,
	}
	err = h.store.SetProfilePrivacy(profileToUpdate, privacyReq.ProfileType)
	if err != nil {
		log.Error("Failed to update profile privacy", "error", err)
		return err
	}
	
	log.Info("Profile privacy updated", "user_id", userID, "type", privacyReq.ProfileType)
	
	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Privacy setting updated successfully",
	})
}

// HandleGetProfileActivity handles retrieving a user's activity
func (h *ProfileHandler) HandleGetProfileActivity(w http.ResponseWriter, r *http.Request) error {
	// TODO: Implement profile activity functionality
	// This would include recent posts, comments, likes, etc.
	
	// For now, return a placeholder response
	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Profile activity feature coming soon",
	})
}