package main

type follow struct {
	FollowerID int64         `json:"follower_id"`
	FollowedID int64         `json:"followed_id"`
	Status     string        `json:"status"`
	Followers  []UserAccount `json:"followers"`
}
