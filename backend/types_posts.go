package main

type Post struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	Privacy   string `json:"privacy"`
	CreatedAt string `json:"created_at"`
}
