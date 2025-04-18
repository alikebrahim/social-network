package api

import (
	"net/http"

	"socialNetwork/pkg/api/handlers"
	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/logger"
)

// Router sets up all API routes and handlers
func (s *APIServer) setupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	log := logger.GetLogger("api")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(s.Store)
	chatHandler := handlers.NewChatHandler(s.Store, s.Hub)
	profileHandler := handlers.NewProfileHandler(s.Store)
	followingHandler := handlers.NewFollowingHandler(s.Store)
	postsHandler := handlers.NewPostsHandler(s.Store)
	groupsHandler := handlers.NewGroupsHandler(s.Store)
	notificationHandler := handlers.NewNotificationHandler(s.Store)

	// Add logger middleware
	logMiddleware := logger.RequestMiddleware()

	// Helper function for routes that require JWT authentication
	withJWTAuth := func(handlerFunc func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) error {
		return func(w http.ResponseWriter, r *http.Request) error {
			// Get logger from request context
			reqLog := logger.FromRequest(r)
			
			// Get session token from cookie
			cookie, err := r.Cookie("session_token")
			if err != nil {
				reqLog.Warn("Authentication required but no session token found")
				return utils.WriteJson(w, http.StatusUnauthorized, map[string]string{
					"error": "Authentication required",
				})
			}

			// Validate the session
			userId, err := s.Store.GetUserIdBySession(cookie.Value)
			if err != nil {
				reqLog.Warn("Invalid or expired session", "error", err)
				return utils.WriteJson(w, http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired session",
				})
			}
			
			// Add user ID to logger context
			reqLog = reqLog.WithField("user_id", userId)
			ctx := logger.WithContext(r.Context(), reqLog)
			r = r.WithContext(ctx)
			
			reqLog.Debug("User authenticated successfully")

			// Call the original handler
			return handlerFunc(w, r)
		}
	}

	// Wrap HandleFunc to apply logger middleware
	handleWithLog := func(pattern string, handler http.HandlerFunc) {
		mux.Handle(pattern, logMiddleware(handler))
	}

	// AUTH ROUTES
	handleWithLog("POST /auth/register", utils.MakeHTTPHandleFunc(authHandler.HandleRegister))
	handleWithLog("POST /auth/login", utils.MakeHTTPHandleFunc(authHandler.HandleLogin))
	handleWithLog("POST /auth/logout", utils.MakeHTTPHandleFunc(authHandler.HandleLogout))

	// PROFILE ROUTES
	handleWithLog("GET /profiles/{id}", utils.MakeHTTPHandleFunc(profileHandler.HandleGetProfile))
	handleWithLog("PUT /profiles/privacy", utils.MakeHTTPHandleFunc(withJWTAuth(profileHandler.HandleSetProfilePrivacy)))
	handleWithLog("GET /profiles/{id}/activity", utils.MakeHTTPHandleFunc(withJWTAuth(profileHandler.HandleGetProfileActivity)))

	// FOLLOWING ROUTES
	handleWithLog("POST /follow/{id}", utils.MakeHTTPHandleFunc(withJWTAuth(followingHandler.HandleFollowing)))
	handleWithLog("GET /follow/requests", utils.MakeHTTPHandleFunc(withJWTAuth(followingHandler.HandleFollowingRequests)))
	handleWithLog("POST /follow/{id}/accept", utils.MakeHTTPHandleFunc(withJWTAuth(followingHandler.HandleFollowingAccept)))
	handleWithLog("DELETE /follow/{id}", utils.MakeHTTPHandleFunc(withJWTAuth(followingHandler.HandleFollowingReject)))

	// POSTS ROUTES
	handleWithLog("POST /posts", utils.MakeHTTPHandleFunc(withJWTAuth(postsHandler.HandlePostCreate)))
	handleWithLog("GET /posts/{id}", utils.MakeHTTPHandleFunc(postsHandler.HandlePostsGet))
	handleWithLog("PUT /posts/{id}", utils.MakeHTTPHandleFunc(withJWTAuth(postsHandler.HandlePostEdit)))
	handleWithLog("DELETE /posts/{id}", utils.MakeHTTPHandleFunc(withJWTAuth(postsHandler.HandlePostDelete)))
	handleWithLog("POST /posts/{id}/comments", utils.MakeHTTPHandleFunc(withJWTAuth(postsHandler.HandlePostComment)))
	handleWithLog("POST /posts/{id}/likes", utils.MakeHTTPHandleFunc(withJWTAuth(postsHandler.HandlePostLike)))
	handleWithLog("DELETE /posts/{id}/likes", utils.MakeHTTPHandleFunc(withJWTAuth(postsHandler.HandlePostUnlike)))

	// GROUPS ROUTES
	handleWithLog("POST /groups", utils.MakeHTTPHandleFunc(withJWTAuth(groupsHandler.HandleGroupCreate)))
	handleWithLog("GET /groups", utils.MakeHTTPHandleFunc(groupsHandler.HandleGroupList))
	handleWithLog("GET /groups/{id}", utils.MakeHTTPHandleFunc(groupsHandler.HandleGroupGet))
	handleWithLog("POST /groups/{id}/invite", utils.MakeHTTPHandleFunc(withJWTAuth(groupsHandler.HandleGroupInvite)))
	handleWithLog("GET /groups/search", utils.MakeHTTPHandleFunc(groupsHandler.HandleGroupSearch))
	handleWithLog("POST /groups/{id}/join", utils.MakeHTTPHandleFunc(withJWTAuth(groupsHandler.HandleGroupRequest)))
	handleWithLog("GET /groups/{id}/requests", utils.MakeHTTPHandleFunc(withJWTAuth(groupsHandler.HandleGroupListRequest)))
	handleWithLog("POST /groups/{id}/requests/{userId}/accept", utils.MakeHTTPHandleFunc(withJWTAuth(groupsHandler.HandleGroupReqAccept)))
	handleWithLog("POST /groups/{id}/requests/{userId}/reject", utils.MakeHTTPHandleFunc(withJWTAuth(groupsHandler.HandleGroupReqReject)))
	handleWithLog("GET /groups/{id}/members", utils.MakeHTTPHandleFunc(groupsHandler.HandleGroupListMembers))
	handleWithLog("POST /groups/{id}/events", utils.MakeHTTPHandleFunc(withJWTAuth(groupsHandler.HandleGroupCreateEvent)))
	handleWithLog("GET /groups/{id}/events", utils.MakeHTTPHandleFunc(groupsHandler.HandleGroupListEvents))
	handleWithLog("POST /groups/{id}/events/{eventId}/response", utils.MakeHTTPHandleFunc(withJWTAuth(groupsHandler.HandleEventResponse)))

	// CHAT ROUTES
	handleWithLog("GET /chats", utils.MakeHTTPHandleFunc(withJWTAuth(chatHandler.HandleGetChats)))
	handleWithLog("GET /chats/{userId}", utils.MakeHTTPHandleFunc(withJWTAuth(chatHandler.HandleGetChatHistory)))
	handleWithLog("GET /groups/{id}/chat", utils.MakeHTTPHandleFunc(withJWTAuth(chatHandler.HandleGetGroupChatHistory)))
	
	// WebSocket handlers (these don't use makeHTTPHandleFunc because they handle their own responses)
	handleWithLog("/ws/chat/{userId}", chatHandler.HandleChatWebSocket)
	handleWithLog("/ws/groups/{id}/chat", chatHandler.HandleGroupChatWebSocket)
	
	// NOTIFICATION ROUTES
	handleWithLog("GET /notifications", utils.MakeHTTPHandleFunc(withJWTAuth(notificationHandler.HandleGetNotifications)))
	handleWithLog("PUT /notifications/{id}/read", utils.MakeHTTPHandleFunc(withJWTAuth(notificationHandler.HandleMarkNotificationRead)))
	handleWithLog("PUT /notifications/read-all", utils.MakeHTTPHandleFunc(withJWTAuth(notificationHandler.HandleMarkAllNotificationsRead)))
	handleWithLog("GET /notifications/unread-count", utils.MakeHTTPHandleFunc(withJWTAuth(notificationHandler.HandleGetUnreadCount)))

	log.Info("Router setup complete", "routes", 37)
	return mux
}

// Run starts the API server
func (s *APIServer) Run() {
	log := logger.GetLogger("api")
	
	// Set up the router
	mux := s.setupRouter()

	// Start the WebSocket hub in a goroutine
	go s.Hub.Run()

	// Start the server
	log.Info("Server running on port", "address", s.ListenAddr)
	err := http.ListenAndServe(s.ListenAddr, mux)
	log.Error("Server error", "error", err)
}