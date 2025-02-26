package main

import "net/http"

// PROFILES HANDLERS
// GET /profiles/{id}
func (s *APIServer) HandleProfileGet(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// PUT /profiles/privacy
func (s *APIServer) HandleProfileSetPrivacy(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// GET /profiles/{id}/activity
func (s *APIServer) HandleProfileGetActivity(w http.ResponseWriter, r *http.Request) error {
	return nil

}
