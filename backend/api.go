package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

// json
func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// server API
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()

	// AUTH HANDLERS
	// POST /auth/register
	mux.HandleFunc("POST /auth/register", makeHTTPHandleFunc(s.HandleRegister))
	// POST /auth/login
	mux.HandleFunc("POST /auth/login", makeHTTPHandleFunc(s.HandleLogin))
	// POST /auth/logout
	mux.HandleFunc("POST /auth/logout", makeHTTPHandleFunc(s.HandleLogout))

	// FOLLOWING HANDLERS
	// POST /follow/{id}
	mux.HandleFunc("POST /follow/{id}", makeHTTPHandleFunc(s.HandleFollowing))
	// GET /follow/requests
	mux.HandleFunc("GET /follow/requests", makeHTTPHandleFunc(s.HandleFollowingRequests))
	// POST /follow/{id}/accept
	mux.HandleFunc("POST /follow/{id}/accept", makeHTTPHandleFunc(s.HandleFollowingAccept))
	// DELETE /follow/{id}
	mux.HandleFunc("DELETE /follow/{id}", makeHTTPHandleFunc(s.HandleFollowingReject))

	// GROUPS HANDLERS
	// POST /groups
	mux.HandleFunc("POST /groups", makeHTTPHandleFunc(s.HandleGroupCreate))
	// GET /groups
	mux.HandleFunc("GET /groups", makeHTTPHandleFunc(s.HandleGroupList))
	// GET /groups/{id}
	mux.HandleFunc("GET /groups/{id}", makeHTTPHandleFunc(s.HandleGroupGet))
	// POST /groups/{id}/invite
	mux.HandleFunc("POST /groups/{id}/invite", makeHTTPHandleFunc(s.HandleGroupInvite))
	// GET /groups/search
	mux.HandleFunc("GET /groups/search", makeHTTPHandleFunc(s.HandleGroupSearch))
	// POST /groups/{id}/join
	mux.HandleFunc("POST /groups/{id}/join", makeHTTPHandleFunc(s.HandleGroupRequest))
	// GET /groups/{id}/requests
	mux.HandleFunc("GET /groups/{id}/requests", makeHTTPHandleFunc(s.HandleGroupListRequest))
	// POST /groups/{id}/requests/{userId}/accept
	mux.HandleFunc("POST /groups/{id}/requests/{userId}/accept", makeHTTPHandleFunc(s.HandleGroupRecAccept))
	// POST /groups/{id}/requests/{userId}/reject
	mux.HandleFunc("POST /groups/{id}/requests/{userId}/reject", makeHTTPHandleFunc(s.HandleGroupRecReject))
	// GET /groups/{id}/members
	mux.HandleFunc("GET /groups/{id}/members", makeHTTPHandleFunc(s.HandleGroupListMembers))
	// POST /groups/{id}/events
	mux.HandleFunc("POST /groups/{id}/events", makeHTTPHandleFunc(s.HandleGroupCreateEvent))
	// GET /groups/{id}/events
	mux.HandleFunc("GET /groups/{id}/events", makeHTTPHandleFunc(s.HandleGroupListEvents))
	// POST /groups/{id}/events/{eventId}/response
	mux.HandleFunc("POST /groups/{id}/events/{eventId}/response", makeHTTPHandleFunc(s.HandleEventResponse))
	// GET /groups/{id}/chat
	mux.HandleFunc("GET /groups/{id}/chat", makeHTTPHandleFunc(s.HandleGetGroupChatHistory))

	// POSTS HANDLERS
	// POST /posts
	mux.HandleFunc("POST /posts", makeHTTPHandleFunc(s.HandlePostCreate))
	// GET /posts/{id}
	mux.HandleFunc("GET /posts/{id}", makeHTTPHandleFunc(s.HandlePostsGet))
	// PUT /posts/{id}
	mux.HandleFunc("PUT /posts/{id}", makeHTTPHandleFunc(s.HandlePostEdit))
	// DELETE /posts/{id}
	mux.HandleFunc("DELETE /posts/{id}", makeHTTPHandleFunc(s.HandlePostDelete))
	// POST /posts/{id}/comments
	mux.HandleFunc("POST /posts/{id}/comments", makeHTTPHandleFunc(s.HandlePostComment))
	// POST /posts/{id}/likes
	mux.HandleFunc("POST /posts/{id}/likes", makeHTTPHandleFunc(s.HandlePostLike))

	// PROFILES HANDLERS
	// GET /profiles/{id}
	mux.HandleFunc("GET /profiles/{id}", makeHTTPHandleFunc(s.HandleProfileGet))
	// PUT /profiles/privacy
	mux.HandleFunc("PUT /profiles/privacy", makeHTTPHandleFunc(s.HandleProfileSetPrivacy))
	// GET /profiles/{id}/activity
	mux.HandleFunc("GET /profiles/{id}/activity", makeHTTPHandleFunc(s.HandleProfileGetActivity))

	// CHAT HANDLERS
	// GET /chats
	mux.HandleFunc("GET /chats", makeHTTPHandleFunc(s.HandleGetChats))
	// GET /chats/{userId}
	mux.HandleFunc("GET /chats/{userId}", makeHTTPHandleFunc(s.HandleGetChatHistory))
	
	// WebSocket handlers (these do not use makeHTTPHandleFunc)
	// WebSocket /ws/chat/{userId}
	mux.HandleFunc("/ws/chat/{userId}", s.HandleChatWebSocket)
	// WebSocket /ws/groups/{id}/chat
	mux.HandleFunc("/ws/groups/{id}/chat", s.HandleGroupChatWebSocket)

	// Start the WebSocket hub in a goroutine
	go hub.Run()

	log.Println("Server running on port: ", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, mux))
}
