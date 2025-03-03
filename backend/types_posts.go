package main

type Post struct {
	PostID int    `json:"post_id"`
	UserID int    `json:"user_id"`
	Text   string `json:"text"`
}
