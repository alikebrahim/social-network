package api

import (
	"database/sql"
	"net/http"

	"backend/pkg/api/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, db *sql.DB) {
	// Authentication routes
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", handlers.RegisterHandler(db)).Methods("POST")
	authRouter.HandleFunc("/login", handlers.LoginHandler(db)).Methods("POST")

	// Profile routes
	profileRouter := router.PathPrefix("/profile").Subrouter()
	profileRouter.HandleFunc("/{userID}", handlers.GetProfileHandler(db)).Methods("GET")
	profileRouter.HandleFunc("/{userID}/follow", handlers.FollowHandler(db)).Methods("POST")
	profileRouter.HandleFunc("/{userID}/unfollow", handlers.UnfollowHandler(db)).Methods("POST")

	// Post routes
	postRouter := router.PathPrefix("/posts").Subrouter()
	postRouter.HandleFunc("", handlers.CreatePostHandler(db)).Methods("POST")
	postRouter.HandleFunc("/{postID}", handlers.GetPostHandler(db)).Methods("GET")
	postRouter.HandleFunc("/{postID}/comments", handlers.AddCommentHandler(db)).Methods("POST")

	// Group routes
	groupRouter := router.PathPrefix("/groups").Subrouter()
	groupRouter.HandleFunc("", handlers.CreateGroupHandler(db)).Methods("POST")
	groupRouter.HandleFunc("/{groupID}", handlers.GetGroupHandler(db)).Methods("GET")
	groupRouter.HandleFunc("/{groupID}/invite", handlers.InviteToGroupHandler(db)).Methods("POST")
	groupRouter.HandleFunc("/{groupID}/join", handlers.JoinGroupHandler(db)).Methods("POST")
	groupRouter.HandleFunc("/{groupID}/events", handlers.CreateEventHandler(db)).Methods("POST")

	// Chat routes
	chatRouter := router.PathPrefix("/chat").Subrouter()
	chatRouter.HandleFunc("/private/{userID}", handlers.PrivateChatHandler(db)).Methods("GET")
	chatRouter.HandleFunc("/group/{groupID}", handlers.GroupChatHandler(db)).Methods("GET")

	// Notification routes
	router.HandleFunc("/notifications", handlers.GetNotificationsHandler(db)).Methods("GET")

	// Static files (e.g., images)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
}
