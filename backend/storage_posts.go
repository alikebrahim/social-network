package main

import "log"

func (s *SQLiteStore) CreatePost(post Post) (id int64, err error) {
	res, err := s.db.Exec("INSERT INTO posts (user_id, content, image, privacy, created_at) VALUES (?, ?, ?, ?, ?)",
		post.UserID, post.Content, post.Image, post.Privacy, post.CreatedAt)
	if err != nil {
		log.Println("CreatePost error :", err)
		return 0, err
	}

	id, err = res.LastInsertId()

	return id, err
}
