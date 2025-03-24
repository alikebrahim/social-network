package main

type Profile struct {
	RequestID   int64
	IsFollowing bool
	UserData    UserAccount `json:"user_data"`
	Privacy     string      `json:"privacy"`
	Posts       []Post      `json:"posts"`
}
