package main

import "net/http"

// GROUPS HANDLERS
// POST /groups
func (s *APIServer) HandleGroupCreate(w http.ResponseWriter, r *http.Request) error {

	return nil

}

// POST /groups/{id}/invite
func (s *APIServer) HandleGroupInvite(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// GET /groups/search
func (s *APIServer) HandleGroupSearch(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// POST /groups/{id}/join
func (s *APIServer) HandleGroupRequest(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// GET /groups/{id}/requests
func (s *APIServer) HandleGroupListRequest(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// POST /groups/{id}/requests/{userId}/accept
func (s *APIServer) HandleGroupRecAccept(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// POST /groups/{id}/events
func (s *APIServer) HandleGroupCreateEvent(w http.ResponseWriter, r *http.Request) error {
	return nil

}
