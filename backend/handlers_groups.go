package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

// GROUPS HANDLERS
// POST /groups
func (s *APIServer) HandleGroupCreate(w http.ResponseWriter, r *http.Request) error {
	// Parse request body
	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Validate input
	if req.Title == "" {
		return ErrInvalidInput
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Create group
	groupID, err := s.store.CreateGroup(userID, &req)
	if err != nil {
		return err
	}

	// Get the created group
	group, err := s.store.GetGroup(groupID)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusCreated, group)
}

// GET /groups
func (s *APIServer) HandleGroupList(w http.ResponseWriter, r *http.Request) error {
	// Get all groups
	groups, err := s.store.GetGroups()
	if err != nil {
		return err
	}

	// If empty, return empty array instead of null
	if groups == nil {
		groups = []*Group{}
	}

	return WriteJson(w, http.StatusOK, groups)
}

// GET /groups/{id}
func (s *APIServer) HandleGroupGet(w http.ResponseWriter, r *http.Request) error {
	// Parse group ID from URL
	groupID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}

	// Get group
	group, err := s.store.GetGroup(groupID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	return WriteJson(w, http.StatusOK, group)
}

// POST /groups/{id}/invite
func (s *APIServer) HandleGroupInvite(w http.ResponseWriter, r *http.Request) error {
	// Parse group ID from URL
	groupID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}

	// Parse request body
	var invite GroupInvite
	if err := json.NewDecoder(r.Body).Decode(&invite); err != nil {
		return err
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Validate that inviter is the authenticated user
	if invite.InviterID != userID {
		return ErrUnauthorized
	}

	// Invite user to group
	err = s.store.InviteToGroup(groupID, userID, invite.InviteeID)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, map[string]string{"status": "invited"})
}

// GET /groups/search
func (s *APIServer) HandleGroupSearch(w http.ResponseWriter, r *http.Request) error {
	// Get search query from URL parameters
	query := r.URL.Query().Get("q")
	if query == "" {
		return ErrInvalidInput
	}

	// Search groups
	groups, err := s.store.SearchGroups(query)
	if err != nil {
		return err
	}

	// If empty, return empty array instead of null
	if groups == nil {
		groups = []*Group{}
	}

	return WriteJson(w, http.StatusOK, groups)
}

// POST /groups/{id}/join
func (s *APIServer) HandleGroupRequest(w http.ResponseWriter, r *http.Request) error {
	// Parse group ID from URL
	groupID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Request to join group
	err = s.store.RequestJoinGroup(groupID, userID)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, map[string]string{"status": "requested"})
}

// GET /groups/{id}/requests
func (s *APIServer) HandleGroupListRequest(w http.ResponseWriter, r *http.Request) error {
	// Parse group ID from URL
	groupID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Check if user is group creator
	isCreator, err := s.store.IsGroupCreator(groupID, userID)
	if err != nil {
		return err
	}
	if !isCreator {
		return ErrUnauthorized
	}

	// Get group join requests
	requests, err := s.store.GetGroupJoinRequests(groupID)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, requests)
}

// POST /groups/{id}/requests/{userId}/accept
func (s *APIServer) HandleGroupRecAccept(w http.ResponseWriter, r *http.Request) error {
	// Parse group ID and user ID from URL
	groupID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}
	requestUserID, err := strconv.ParseInt(r.PathValue("userId"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Check if user is group creator
	isCreator, err := s.store.IsGroupCreator(groupID, userID)
	if err != nil {
		return err
	}
	if !isCreator {
		return ErrUnauthorized
	}

	// Accept join request
	err = s.store.AcceptGroupJoinRequest(groupID, requestUserID)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, map[string]string{"status": "accepted"})
}

// POST /groups/{id}/requests/{userId}/reject
func (s *APIServer) HandleGroupRecReject(w http.ResponseWriter, r *http.Request) error {
	// Parse group ID and user ID from URL
	groupID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}
	requestUserID, err := strconv.ParseInt(r.PathValue("userId"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Check if user is group creator
	isCreator, err := s.store.IsGroupCreator(groupID, userID)
	if err != nil {
		return err
	}
	if !isCreator {
		return ErrUnauthorized
	}

	// Reject join request
	err = s.store.RejectGroupJoinRequest(groupID, requestUserID)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, map[string]string{"status": "rejected"})
}

// GET /groups/{id}/members
func (s *APIServer) HandleGroupListMembers(w http.ResponseWriter, r *http.Request) error {
	// Parse group ID from URL
	groupID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}

	// Get group members
	members, err := s.store.GetGroupMembers(groupID)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, members)
}

// POST /groups/{id}/events
func (s *APIServer) HandleGroupCreateEvent(w http.ResponseWriter, r *http.Request) error {
	// Parse group ID from URL
	groupID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}

	// Parse request body
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Validate input
	if req.Title == "" || req.EventTime.IsZero() {
		return ErrInvalidInput
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Check if user is a member of the group
	isMember, err := s.store.IsGroupMember(groupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrUnauthorized
	}

	// Create event
	event := &GroupEvent{
		GroupID:     groupID,
		CreatorID:   userID,
		Title:       req.Title,
		Description: req.Description,
		EventTime:   req.EventTime,
		CreatedAt:   time.Now(),
	}

	eventID, err := s.store.CreateGroupEvent(event)
	if err != nil {
		return err
	}

	// Get the created event
	event.ID = eventID
	return WriteJson(w, http.StatusCreated, event)
}

// GET /groups/{id}/events
func (s *APIServer) HandleGroupListEvents(w http.ResponseWriter, r *http.Request) error {
	// Parse group ID from URL
	groupID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Check if user is a member of the group
	isMember, err := s.store.IsGroupMember(groupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrUnauthorized
	}

	// Get group events
	events, err := s.store.GetGroupEvents(groupID)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, events)
}

// POST /groups/{id}/events/{eventId}/response
func (s *APIServer) HandleEventResponse(w http.ResponseWriter, r *http.Request) error {
	// Parse group ID and event ID from URL
	groupID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}
	eventID, err := strconv.ParseInt(r.PathValue("eventId"), 10, 64)
	if err != nil {
		return ErrInvalidInput
	}

	// Parse request body
	var req EventResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Validate input
	if req.Response != "going" && req.Response != "not_going" {
		return ErrInvalidInput
	}

	// Get user ID from session
	userID, err := s.getUserIDFromSession(r)
	if err != nil {
		return ErrUnauthorized
	}

	// Check if user is a member of the group
	isMember, err := s.store.IsGroupMember(groupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrUnauthorized
	}

	// Respond to event
	err = s.store.RespondToEvent(eventID, userID, req.Response)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, map[string]string{"status": "responded"})
}

// Helper function to get user ID from session
func (s *APIServer) getUserIDFromSession(r *http.Request) (int64, error) {
	//NOTE: fix once session managemnt is implemented
	// Place holder for now
	userID := int64(1) // Assuming user ID 1 for development
	return userID, nil
}
