package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Authentication HANDLERS
// POST /auth/register
func (s *APIServer) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	var user UserAccount

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("json decode ", err)
	}

	id, err := s.store.CreateUserAccount(&user)
	if err != nil {
		fmt.Println("createUserAccount err: ", err)
	}

	WriteJson(w, http.StatusCreated, id)
	return nil
}

// POST /auth/login
func (s *APIServer) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// POST /auth/logout
func (s *APIServer) HandleLogout(w http.ResponseWriter, r *http.Request) error {
	return nil
}
