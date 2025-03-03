package main

func (s *SQLiteStore) CreatePost(post Post) (id int64, err error) {
	res, err := s.db.Exec("INSERT INTO posts (user_id, text) VALUES (?, ?)", post.UserID, post.Text)
	if err != nil {
		return 0, err
	}

	id, err = res.LastInsertId()

	return id, err
}
