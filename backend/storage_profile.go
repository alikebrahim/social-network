package main

import (
	"log"
)

func (s *SQLiteStore) GetProfileData(Profile *Profile) (*Profile, error) {

	err := s.db.QueryRow(`SELECT id, email, first_name, last_name, about_me, avatar, nickname, profile_type
		FROM users
		WHERE id = ?`, Profile.UserData.ID).Scan(
		&Profile.UserData.ID,
		&Profile.UserData.Email,
		&Profile.UserData.First_name,
		&Profile.UserData.Last_name,
		&Profile.UserData.About_me,
		&Profile.UserData.Avatar,
		&Profile.UserData.Nickname,
		&Profile.UserData.Profile_type,
	)
	if err != nil {
		log.Println("Error getting profile data:", err)
		return nil, err
	}

	postRows, err := s.db.Query("SELECT id, user_id, content, image, privacy, created_at FROM posts WHERE user_id = ?", Profile.UserData.ID)
	if err != nil {
		log.Println("Error getting posts:", err)
		return nil, err
	}
	defer postRows.Close()
	for postRows.Next() {
		var post Post
		if err := postRows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&post.Image,
			&post.Privacy,
			&post.CreatedAt,
		); err != nil {
			log.Println("Error scanning post:", err)
			return nil, err
		}
		if post.Privacy == "public" {
			Profile.Posts = append(Profile.Posts, post)
		}

		canSee, err := s.CanUserSeePost(post.ID, Profile.UserData.ID)
		if err != nil {
			log.Println("Error checking post visibility:", err)
			return nil, err
		}
		if canSee {
			Profile.Posts = append(Profile.Posts, post)
		}
	}

	return Profile, nil
}

func (s *SQLiteStore) SetProfilePrivacy(Profile Profile, privacy string) error {
	query := "UPDATE users SET privacy = ? WHERE id = ?"
	_, err := s.db.Exec(query, privacy, Profile.UserData.ID)
	if err != nil {
		log.Println("Error setting profile privacy:", err)
		return err
	}

	return nil
}
