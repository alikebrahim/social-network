package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
)

// Authentication HANDLERS
// POST /auth/register
func (s *APIServer) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	var user UserAccount

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print("json.NewDecoder: ", err)
		return err
	}

	id, err := s.store.CreateUserAccount(&user)
	if err != nil {
		log.Print("store.CreateUserAccount: ", err)
		return err
	}
	// set cookie
	sessionToken, err := s.store.GetSeesionToken(id)
	if err != nil {
		log.Print("store.GetSeesionToken: ", err)
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HttpOnly: true,
		Secure:   true,
	})

	WriteJson(w, http.StatusCreated, id)
	return nil
}

// POST /auth/login
func (s *APIServer) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	var user UserAccount

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print("json.NewDecoder: ", err)
		return err
	}
	regUser, err := s.store.AuthenticateUser(user.Email, user.Password)
	if err != nil {
		log.Print("store.AuthenticateUser: ", err)
		return err
	}
	//log.Print("regUser: ", regUser)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    regUser,
		HttpOnly: true,
		Secure:   true,
	})
	return nil
}

// POST /auth/logout
func (s *APIServer) HandleLogout(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Print("r.Cookies: ", err)
		return err
	}
	sessionToken := cookie.Value
	err = s.store.DeleteSession(sessionToken)
	if err != nil {
		log.Print("store.DeleteSession: ", err)
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
	return nil
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getSessionToken(r *http.Request) (string, error) {
	log.Print("r.Cookies: ", r.Cookies())
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}
	log.Print("cookie: ", cookie)
	return cookie.Value, nil
}
