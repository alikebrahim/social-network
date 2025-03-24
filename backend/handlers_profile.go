package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
)

// PROFILES HANDLERS
// GET /profiles/{id}
func (s *APIServer) HandleProfileGet(w http.ResponseWriter, r *http.Request) error {
	var profile Profile
	// get profile id from request
	id := r.PathValue("id")
	var err error
	profile.UserData.ID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("strconv error :", err)
		return err
	}
	// get how requsteing the profile 
	sessionToken, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}
	profile.RequestID, err = s.store.GetUserIdBySession(sessionToken)
	if err != nil {
		log.Println("store.GetUserIDBySession error :", err)
		return err
	}
	// check if the requster is following the profile
	isFollowing, err := s.store.isFollowing(profile.RequestID, profile.UserData.ID)
	if err != nil {
		log.Println("store.isFollowing error :", err)
		return err
	}
	

	//log.Println("isFollowing:", isFollowing)
	
	// get profile data
	_, err = s.store.GetProfileData(&profile)
	if err != nil {
		log.Println("store.GetProfile error :", err)
		return err
	}
	// hide email and date of birth if the profile is private and the requster is not following
	if profile.UserData.Profile_type == "private" && !isFollowing {
		profile.UserData.Email = ""
		profile.UserData.Date_of_birth = ""
		profile.Posts = nil
	}
	// write profile data to response
	err = WriteJson(w, http.StatusOK, profile)
	if err != nil {
		log.Println("WriteJson error :", err)
		return err
	}
	return nil
}

// PUT /profiles/privacy
func (s *APIServer) HandleProfileSetPrivacy(w http.ResponseWriter, r *http.Request) error {
	var profile Profile
	sessionToken, err := getSessionToken(r)
	if err != nil {
		log.Println("getSessionToken error :", err)
		return err
	}
	profile.UserData.ID, err = s.store.GetUserIdBySession(sessionToken)
	if err != nil {
		log.Println("store.GetUserIDBySession error :", err)
		return err
	}

	// get privacy from request
	NewPrivacy := r.FormValue("privacy")

	_, err = s.store.GetProfileData(&profile)
	if err != nil {
		log.Println("store.GetProfileData error :", err)
		return err
	}

	if profile.Privacy == NewPrivacy {
		log.Println("Privacy is already set to", NewPrivacy)
		return errors.New("Privacy is already set to " + NewPrivacy)
	}


	err = s.store.SetProfilePrivacy(profile, NewPrivacy)
	if err != nil {
		log.Println("store.SetProfilePrivacy error :", err)
		return err
	}

	return nil

}

// GET /profiles/{id}/activity
func (s *APIServer) HandleProfileGetActivity(w http.ResponseWriter, r *http.Request) error {
	return nil

}
