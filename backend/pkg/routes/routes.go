package routes

import (
	"database/sql"
	"net/http"
)


func SetupRoutes(database *sql.DB) *http.ServeMux {
	DB = database 

	mux := http.NewServeMux()


	mux.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r, DB)
	})
	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		LoginHandler(w, r, DB)
	})
	mux.HandleFunc("/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		LogoutHandler(w, r, DB)
	})

	
	mux.HandleFunc("/follow/", func(w http.ResponseWriter, r *http.Request) {
		FollowingHandler(w, r, DB)
	})
	mux.HandleFunc("/follow/requests", func(w http.ResponseWriter, r *http.Request) {
		FollowingRequestsHandler(w, r, DB)
	})
	mux.HandleFunc("/follow/accept/", func(w http.ResponseWriter, r *http.Request) {
		FollowAcceptRequestHandler(w, r, DB)
	})
	mux.HandleFunc("/follow/reject/", func(w http.ResponseWriter, r *http.Request) {
		FollowRejectRequestHandler(w, r, DB)
	})

	
	mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		GroupCreateHandler(w, r, DB)
	})
	mux.HandleFunc("/groups/invite/", func(w http.ResponseWriter, r *http.Request) {
		GroupInviteHandler(w, r, DB)
	})
	mux.HandleFunc("/groups/search", func(w http.ResponseWriter, r *http.Request) {
		GroupSearchHandler(w, r, DB)
	})
	mux.HandleFunc("/groups/join/", func(w http.ResponseWriter, r *http.Request) {
		GroupRequestHandler(w, r, DB)
	})
	mux.HandleFunc("/groups/requests/", func(w http.ResponseWriter, r *http.Request) {
		GroupListRequestsHandler(w, r, DB)
	})
	mux.HandleFunc("/groups/requests/accept/", func(w http.ResponseWriter, r *http.Request) {
		GroupAcceptRequestHandler(w, r, DB)
	})
	mux.HandleFunc("/groups/events/", func(w http.ResponseWriter, r *http.Request) {
		GroupCreateEventHandler(w, r, DB)
	})

	
	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		PostCreateHandler(w, r, DB)
	})
	mux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		PostsGetHandler(w, r, DB)
	})
	mux.HandleFunc("/posts/edit/", func(w http.ResponseWriter, r *http.Request) {
		PostEditHandler(w, r, DB)
	})
	mux.HandleFunc("/posts/delete/", func(w http.ResponseWriter, r *http.Request) {
		PostDeleteHandler(w, r, DB)
	})
	mux.HandleFunc("/posts/comments/", func(w http.ResponseWriter, r *http.Request) {
		PostCommentHandler(w, r, DB)
	})
	mux.HandleFunc("/posts/likes/", func(w http.ResponseWriter, r *http.Request) {
		PostLikeHandler(w, r, DB)
	})

	
	mux.HandleFunc("/profiles/", func(w http.ResponseWriter, r *http.Request) {
		ProfileGetHandler(w, r, DB)
	})
	mux.HandleFunc("/profiles/privacy", func(w http.ResponseWriter, r *http.Request) {
		ProfileSetPrivacyHandler(w, r, DB)
	})
	mux.HandleFunc("/profiles/activity/", func(w http.ResponseWriter, r *http.Request) {
		ProfileGetActivitiHandler(w, r, DB)
	})

	return mux
}
