package main

// NOTE: *Request structs similar to this one
// could be used for follow requests
// type CreateAccountRequest struct {
// 	Name           string `json:"name"`
// 	Affiliation    string `json:"affiliation"`
// 	PrimaryContact string `json:"primaryContact"`
// }

type UserAccount struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	First_name    string `json:"first_name"`
	Last_name     string `json:"last_name"`
	Date_of_birth string `json:"date_of_birth"`
	Avatar        string `json:"avatar"`
	Nickname      string `json:"nickname"`
	About_me      string `json:"about_me"`
	Profile_type  string `json:"profile_type"`
}

// func NewAccount(email, pass, fName, lName, dob, avatar, nic, abt, pType string) *UserAccount {
// 	return &UserAccount{
// 		Email:         email,
// 		Password:      pass,
// 		First_name:    fName,
// 		Last_name:     lName,
// 		Date_of_birth: dob,
// 		Avatar:        avatar,
// 		Nickname:      nic,
// 		About_me:      abt,
// 		Profile_type:  pType,
// 	}
// }
