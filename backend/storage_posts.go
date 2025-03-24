package main

import (
	"errors"
	"log"
)

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

func (s *SQLiteStore) GetPostByID(userId int64, followerId int64) (posts []Post, err error) {
	rows, err := s.db.Query(`
	SELECT posts.id, posts.user_id, users.profile_type, posts.content, posts.image, posts.privacy, posts.created_at 
	FROM posts 
	JOIN users ON posts.user_id = users.id
	WHERE posts.privacy = 'public' 
	OR (posts.privacy = 'private' AND posts.user_id IN 
		(SELECT followed_id FROM followers WHERE follower_id = ? AND status = 'accepted'))
	ORDER BY posts.created_at DESC`, userId)

	if err != nil {
		log.Print(err)
		return posts, err
	}
	defer rows.Close()

	
	for rows.Next() {
		var post Post
		var profileType string

		err := rows.Scan(&post.ID, &post.UserID, &profileType, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt)
		if err != nil {
			log.Println("Error scanning post:", err)
			continue
		}

		if profileType == "private" {
			var exists bool
			err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'accepted')", followerId, post.UserID).Scan(&exists)
			if err != nil {
				log.Println("Error checking profile visibility:", err)
				continue
			}
			log.Println("Profile visibility exists:", exists)
			if !exists {
				continue
			}
		}

		if post.Privacy == "private" {
			var exists bool
			err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM post_visibility WHERE post_id = ? AND viewer_id = ?)", post.ID, followerId).Scan(&exists)
			if err != nil {
				log.Println("Error checking post visibility:", err)
				continue
			}
			log.Println("Post visibility exists:", exists)
			if !exists {
				continue
			}
		}

	

		log.Println(post)

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *SQLiteStore) IsPostOwner(userId int64, followerId int64) (IsPostOwner bool, err error) {
	var isOwner = false
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE user_id = ? AND id = ?)", userId, followerId).Scan(&isOwner)
	if err != nil {
		log.Println("Error checking post owner:", err)
		return isOwner, err
	}
	log.Println("Is post owner:", isOwner)
	return isOwner, nil
}

func (s *SQLiteStore) EditPost(post Post) (err error) {
	log.Println("Edit post:", post)
	if post.Privacy == "" {
		post.Privacy = "public"
	} else if post.Content == "" && post.Image == "" {
		err = errors.New("Content or image is required")
		log.Println("Content or image is required")
		return err
	}
	_, err = s.db.Exec("UPDATE posts SET content = ?, image = ?, privacy = ? WHERE id = ?", post.Content, post.Image, post.Privacy, post.ID)
	if err != nil {
		log.Println("editPost error :", err)
		return err
	}
	return nil
}

func (s *SQLiteStore) DeletePost(postId int64) (err error) {
	_, err = s.db.Exec("DELETE FROM posts WHERE id = ?", postId)
	if err != nil {
		log.Println("deletePost error :", err)
		return err
	}
	return nil
}

func (s *SQLiteStore) CreateComment(comment Comment) (err error) {
	_, err = s.db.Exec("INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
		comment.PostID, comment.ComenterID, comment.Content, comment.CreatedAt)
	if err != nil {
		log.Println("CreateComment error :", err)
		return err
	}
	return nil
}


func (s *SQLiteStore) CreateLike(like like) (err error) {
	var exists bool
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = ? AND user_id = ?)", like.PostID, like.UserID).Scan(&exists)
	if err != nil {
		log.Println("CreateLike error :", err)
		return err
	}
	// if exists update the like
	if exists {
		_, err = s.db.Exec("UPDATE likes SET is_like = ? WHERE post_id = ? AND user_id = ?", like.IsLike, like.PostID, like.UserID)
		if err != nil {
			log.Println("CreateLike error :", err)
			return err
		}
		return nil
	} else if !exists {
		query := `INSERT INTO likes (is_like, post_id, user_id) VALUES (?, ?, ?)`
		_, err = s.db.Exec(query, like.IsLike, like.PostID, like.UserID)
		if err != nil {
			log.Println("CreateLike error :", err)
			return err
		}
	}
	return nil
}

func (s *SQLiteStore) RemoveLikes(like like) (err error) {
	_, err = s.db.Exec("DELETE FROM likes WHERE post_id = ? AND user_id = ?", like.PostID, like.UserID)
	if err != nil {
		log.Println("RemoveLikes error :", err)
		return err
	}
	return nil
}
