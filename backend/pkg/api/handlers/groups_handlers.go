package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/domain/groups"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/logger"
	"socialNetwork/pkg/storage"
)

// GroupsHandler manages groups-related handlers
type GroupsHandler struct {
	store  storage.Storage
	logger logger.Logger
}

// NewGroupsHandler creates a new groups handler
func NewGroupsHandler(store storage.Storage) *GroupsHandler {
	return &GroupsHandler{
		store:  store,
		logger: logger.GetLogger("api/handlers/groups"),
	}
}

// HandleGroupCreate handles creating a new group
func (h *GroupsHandler) HandleGroupCreate(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Parse request body
	var req groups.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return errors.ErrBadRequest
	}

	// Validate group details
	if req.Title == "" {
		log.Warn("Invalid group: empty title")
		return errors.ErrInvalidInput
	}

	// Get the user ID from the session
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

	// Create the group
	groupID, err := h.store.CreateGroup(userID, &req)
	if err != nil {
		log.Error("Failed to create group", "error", err)
		return err
	}

	log.Info("Group created", "group_id", groupID, "user_id", userID)

	// Return the created group ID
	return utils.WriteJson(w, http.StatusCreated, map[string]int64{
		"group_id": groupID,
	})
}

// HandleGroupList handles listing all groups
func (h *GroupsHandler) HandleGroupList(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get all groups
	groups, err := h.store.GetGroups()
	if err != nil {
		log.Error("Failed to get groups", "error", err)
		return err
	}

	log.Info("Groups retrieved", "count", len(groups))

	return utils.WriteJson(w, http.StatusOK, groups)
}

// HandleGroupGet handles retrieving a specific group
func (h *GroupsHandler) HandleGroupGet(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get group ID from the URL
	idStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid group ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get the group
	group, err := h.store.GetGroup(groupID)
	if err != nil {
		if err == errors.ErrNotFound {
			log.Warn("Group not found", "group_id", groupID)
			return errors.ErrNotFound
		}
		log.Error("Failed to get group", "error", err)
		return err
	}

	log.Info("Group retrieved", "group_id", groupID)

	return utils.WriteJson(w, http.StatusOK, group)
}

// HandleGroupInvite handles inviting a user to a group
func (h *GroupsHandler) HandleGroupInvite(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get group ID from the URL
	idStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid group ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Parse request body
	var invite groups.GroupInvite
	if err := json.NewDecoder(r.Body).Decode(&invite); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return errors.ErrBadRequest
	}

	// Get the inviter's ID from the session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	inviterID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Set the group ID and inviter ID
	invite.GroupID = groupID
	invite.InviterID = inviterID

	// Invite the user to the group
	err = h.store.InviteToGroup(groupID, inviterID, invite.InviteeID)
	if err != nil {
		if err == errors.ErrAlreadyExists {
			log.Info("User already invited or member", "group_id", groupID, "user_id", invite.InviteeID)
			return utils.WriteJson(w, http.StatusOK, map[string]string{
				"message": "User is already invited or a member of the group",
			})
		}
		if err == errors.ErrUnauthorized {
			log.Warn("Unauthorized invite attempt", "group_id", groupID, "inviter_id", inviterID)
			return errors.ErrUnauthorized
		}
		log.Error("Failed to invite user", "error", err)
		return err
	}

	log.Info("User invited to group", "group_id", groupID, "inviter_id", inviterID, "invitee_id", invite.InviteeID)

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "User invited successfully",
	})
}

// HandleGroupSearch handles searching for groups
func (h *GroupsHandler) HandleGroupSearch(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get search query from URL parameters
	query := r.URL.Query().Get("q")
	if query == "" {
		log.Warn("Empty search query")
		return errors.ErrBadRequest
	}

	// Search for groups
	groups, err := h.store.SearchGroups(query)
	if err != nil {
		log.Error("Failed to search groups", "error", err)
		return err
	}

	log.Info("Groups search results", "query", query, "count", len(groups))

	return utils.WriteJson(w, http.StatusOK, groups)
}

// HandleGroupRequest handles requesting to join a group
func (h *GroupsHandler) HandleGroupRequest(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get group ID from the URL
	idStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid group ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get the requester's ID from the session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	requesterID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Request to join the group
	err = h.store.RequestJoinGroup(groupID, requesterID)
	if err != nil {
		if err == errors.ErrAlreadyExists {
			log.Info("Already requested or member", "group_id", groupID, "user_id", requesterID)
			return utils.WriteJson(w, http.StatusOK, map[string]string{
				"message": "Already requested to join or a member of the group",
			})
		}
		log.Error("Failed to request joining group", "error", err)
		return err
	}

	log.Info("Join request created", "group_id", groupID, "user_id", requesterID)

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Join request sent successfully",
	})
}

// HandleGroupListRequest handles listing join requests for a group
func (h *GroupsHandler) HandleGroupListRequest(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get group ID from the URL
	idStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid group ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

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

	// Check if the user is the group creator
	isCreator, err := h.store.IsGroupCreator(groupID, userID)
	if err != nil {
		log.Error("Failed to check group creator", "error", err)
		return err
	}

	if !isCreator {
		log.Warn("Unauthorized access to join requests", "group_id", groupID, "user_id", userID)
		return errors.ErrUnauthorized
	}

	// Get join requests
	requests, err := h.store.GetGroupJoinRequests(groupID)
	if err != nil {
		log.Error("Failed to get join requests", "error", err)
		return err
	}

	log.Info("Join requests retrieved", "group_id", groupID, "count", len(requests))

	return utils.WriteJson(w, http.StatusOK, requests)
}

// HandleGroupReqAccept handles accepting a join request
func (h *GroupsHandler) HandleGroupReqAccept(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get group ID and user ID from the URL
	groupIDStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		log.Error("Invalid group ID", "id", groupIDStr, "error", err)
		return errors.ErrBadRequest
	}

	userIDStr := r.PathValue("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Error("Invalid user ID", "id", userIDStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get the current user's ID from the session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	currentUserID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Check if the current user is the group creator
	isCreator, err := h.store.IsGroupCreator(groupID, currentUserID)
	if err != nil {
		log.Error("Failed to check group creator", "error", err)
		return err
	}

	if !isCreator {
		log.Warn("Unauthorized accept request attempt", "group_id", groupID, "user_id", currentUserID)
		return errors.ErrUnauthorized
	}

	// Accept the join request
	err = h.store.AcceptGroupJoinRequest(groupID, userID)
	if err != nil {
		if err == errors.ErrNotFound {
			log.Warn("Join request not found", "group_id", groupID, "user_id", userID)
			return errors.ErrNotFound
		}
		log.Error("Failed to accept join request", "error", err)
		return err
	}

	log.Info("Join request accepted", "group_id", groupID, "user_id", userID)

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Join request accepted",
	})
}

// HandleGroupReqReject handles rejecting a join request
func (h *GroupsHandler) HandleGroupReqReject(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get group ID and user ID from the URL
	groupIDStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		log.Error("Invalid group ID", "id", groupIDStr, "error", err)
		return errors.ErrBadRequest
	}

	userIDStr := r.PathValue("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Error("Invalid user ID", "id", userIDStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get the current user's ID from the session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	currentUserID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Check if the current user is the group creator
	isCreator, err := h.store.IsGroupCreator(groupID, currentUserID)
	if err != nil {
		log.Error("Failed to check group creator", "error", err)
		return err
	}

	if !isCreator {
		log.Warn("Unauthorized reject request attempt", "group_id", groupID, "user_id", currentUserID)
		return errors.ErrUnauthorized
	}

	// Reject the join request
	err = h.store.RejectGroupJoinRequest(groupID, userID)
	if err != nil {
		if err == errors.ErrNotFound {
			log.Warn("Join request not found", "group_id", groupID, "user_id", userID)
			return errors.ErrNotFound
		}
		log.Error("Failed to reject join request", "error", err)
		return err
	}

	log.Info("Join request rejected", "group_id", groupID, "user_id", userID)

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Join request rejected",
	})
}

// HandleGroupListMembers handles listing all members of a group
func (h *GroupsHandler) HandleGroupListMembers(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get group ID from the URL
	idStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid group ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get the group members
	members, err := h.store.GetGroupMembers(groupID)
	if err != nil {
		log.Error("Failed to get group members", "error", err)
		return err
	}

	log.Info("Group members retrieved", "group_id", groupID, "count", len(members))

	return utils.WriteJson(w, http.StatusOK, members)
}

// HandleGroupCreateEvent handles creating a new event in a group
func (h *GroupsHandler) HandleGroupCreateEvent(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get group ID from the URL
	idStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid group ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Parse request body
	var eventReq groups.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&eventReq); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return errors.ErrBadRequest
	}

	// Validate event details
	if eventReq.Title == "" {
		log.Warn("Invalid event: empty title")
		return errors.ErrInvalidInput
	}

	// Get the creator's ID from the session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	creatorID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Create the event
	event := &groups.GroupEvent{
		GroupID:     groupID,
		CreatorID:   creatorID,
		Title:       eventReq.Title,
		Description: eventReq.Description,
		EventTime:   eventReq.EventTime,
		CreatedAt:   time.Now(),
	}

	eventID, err := h.store.CreateGroupEvent(event)
	if err != nil {
		if err == errors.ErrUnauthorized {
			log.Warn("Unauthorized event creation attempt", "group_id", groupID, "user_id", creatorID)
			return errors.ErrUnauthorized
		}
		log.Error("Failed to create event", "error", err)
		return err
	}

	log.Info("Event created", "group_id", groupID, "event_id", eventID, "creator_id", creatorID)

	return utils.WriteJson(w, http.StatusCreated, map[string]int64{
		"event_id": eventID,
	})
}

// HandleGroupListEvents handles listing all events in a group
func (h *GroupsHandler) HandleGroupListEvents(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get group ID from the URL
	idStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid group ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get the group events
	events, err := h.store.GetGroupEvents(groupID)
	if err != nil {
		log.Error("Failed to get group events", "error", err)
		return err
	}

	log.Info("Group events retrieved", "group_id", groupID, "count", len(events))

	return utils.WriteJson(w, http.StatusOK, events)
}

// HandleGroupCreatePost handles creating a new post in a group
func (h *GroupsHandler) HandleGroupCreatePost(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)
	
	// Get group ID from the URL
	idStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid group ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}
	
	// Parse request body
	var post posts.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return errors.ErrBadRequest
	}
	
	// Validate post content
	if post.Content == "" && post.Image == "" {
		log.Warn("Invalid post: empty content and image")
		return errors.ErrInvalidInput
	}
	
	// Validate image format if provided
	if post.Image != "" {
		isValidFormat := false
		validFormats := []string{".jpg", ".jpeg", ".png", ".gif"}
		
		// Simple extension check - in a real app, you would validate the actual file contents
		for _, format := range validFormats {
			if len(post.Image) > len(format) && post.Image[len(post.Image)-len(format):] == format {
				isValidFormat = true
				break
			}
		}
		
		if !isValidFormat {
			log.Warn("Invalid image format", "image", post.Image)
			return utils.WriteJson(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid image format. Supported formats: JPEG, PNG, GIF",
			})
		}
	}
	
	// Get the user ID from the session
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
	
	// Create the post
	post.UserID = userID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	
	postID, err := h.store.CreateGroupPost(groupID, userID, &post)
	if err != nil {
		if err == errors.ErrUnauthorized {
			log.Warn("Unauthorized post creation attempt", "group_id", groupID, "user_id", userID)
			return errors.ErrUnauthorized
		}
		log.Error("Failed to create group post", "error", err)
		return err
	}
	
	log.Info("Group post created", "group_id", groupID, "post_id", postID, "user_id", userID)
	
	return utils.WriteJson(w, http.StatusCreated, map[string]int64{
		"post_id": postID,
	})
}

// HandleGroupListPosts handles listing all posts in a group
func (h *GroupsHandler) HandleGroupListPosts(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)
	
	// Get group ID from the URL
	idStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid group ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}
	
	// Get the user ID from the session
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
	
	// Check if the user is a member of the group
	isMember, err := h.store.IsGroupMember(groupID, userID)
	if err != nil {
		log.Error("Failed to check group membership", "error", err)
		return err
	}
	
	if !isMember {
		log.Warn("Unauthorized access to group posts", "group_id", groupID, "user_id", userID)
		return errors.ErrUnauthorized
	}
	
	// Get the group posts
	posts, err := h.store.GetGroupPosts(groupID)
	if err != nil {
		log.Error("Failed to get group posts", "error", err)
		return err
	}
	
	log.Info("Group posts retrieved", "group_id", groupID, "count", len(posts))
	
	return utils.WriteJson(w, http.StatusOK, posts)
}

// HandleEventResponse handles responding to an event
func (h *GroupsHandler) HandleEventResponse(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get group ID and event ID from the URL
	groupIDStr := r.PathValue("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		log.Error("Invalid group ID", "id", groupIDStr, "error", err)
		return errors.ErrBadRequest
	}

	eventIDStr := r.PathValue("eventId")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		log.Error("Invalid event ID", "id", eventIDStr, "error", err)
		return errors.ErrBadRequest
	}

	// Parse request body
	var responseReq groups.EventResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&responseReq); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return errors.ErrBadRequest
	}

	// Validate response
	if responseReq.Response != "going" && responseReq.Response != "not_going" {
		log.Warn("Invalid response value", "response", responseReq.Response)
		return errors.ErrInvalidInput
	}

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

	// Check if the user is a member of the group
	isMember, err := h.store.IsGroupMember(groupID, userID)
	if err != nil {
		log.Error("Failed to check group membership", "error", err)
		return err
	}

	if !isMember {
		log.Warn("Unauthorized event response attempt", "group_id", groupID, "user_id", userID)
		return errors.ErrUnauthorized
	}

	// Record the response
	err = h.store.RespondToEvent(eventID, userID, responseReq.Response)
	if err != nil {
		log.Error("Failed to record event response", "error", err)
		return err
	}

	log.Info("Event response recorded", "event_id", eventID, "user_id", userID, "response", responseReq.Response)

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Response recorded successfully",
	})
}