package main

type Comment struct {
	ComenterID int64  `json:"comenter_id"`
	PostID     int64  `json:"post_id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}
