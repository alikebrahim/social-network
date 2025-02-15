package routes

import "net/http"

// GROUPS HANDLERS
// POST /groups
func GroupCreateHandler(w http.ResponseWriter, r *http.Request) {

}

// POST /groups/{id}/invite
func GroupInviteHandler(w http.ResponseWriter, r *http.Request) {

}

// GET /groups/search
func GroupSearchHandler(w http.ResponseWriter, r *http.Request) {

}

// POST /groups/{id}/join
func GroupRequestHandler(w http.ResponseWriter, r *http.Request) {

}

// GET /groups/{id}/requests
func GroupListRequestsHandler(w http.ResponseWriter, r *http.Request) {

}

// POST /groups/{id}/requests/{userId}/accept
func GroupAcceptRequestHandler(w http.ResponseWriter, r *http.Request) {

}

// POST /groups/{id}/events
func GroupCreateEventHandler(w http.ResponseWriter, r *http.Request) {

}
