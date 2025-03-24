package main

import (
	"log"
)

func (s *SQLiteStore) GetProfileData(Profile *Profile) (*Profile, error) {
	// queryProfile := "SELECT id, email , first_name , last_name , date_of_birth , avatar , nickname , about_me, profile_type FROM users WHERE id = ?"
	// err := s.db.QueryRow(queryProfile, Profile.UserData.ID).Scan(&Profile.UserData.ID, &Profile.UserData.Email, &Profile.UserData.First_name, &Profile.UserData.Last_name, &Profile.UserData.Date_of_birth, &Profile.UserData.Avatar, &Profile.UserData.Nickname, &Profile.UserData.About_me, &Profile.UserData.Profile_type)
	// if err != nil {
	// 	log.Println("Error getting profile data:", err)
	// 	return nil, err
	// }

	
	// queryPosts := "SELECT id, user_id, content, image, created_at FROM posts WHERE user_id = ?"
	// rows, err := s.db.Query(queryPosts, Profile.UserData.ID)
	// if err != nil {
	// 	log.Println("Error getting posts:", err)
	// 	return nil, err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var post Post
	// 	err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.CreatedAt)
	// 	if err != nil {
	// 		log.Println("Error scanning posts:", err)
	// 		return nil, err
	// 	}
	// 	Profile.Posts = append(Profile.Posts, post)
	// }

	// // get followers count
	// queryFollowers := "SELECT COUNT(*) FROM followers WHERE followed_id = ? AND status = 'accepted'"
	// err = s.db.QueryRow(queryFollowers, Profile.UserData.ID).Scan(&Profile.UserData.FollowCount)
	// if err != nil {
	// 	log.Println("Error getting followers count:", err)
	// 	return nil, err
	// }

	// log.Println("Profile.UserData.ID:", Profile.UserData.ID)

	// // get following count
	// queryFollowing := "SELECT COUNT(*) FROM followers WHERE follower_id = ? AND status = 'accepted'"
	// err = s.db.QueryRow(queryFollowing, Profile.UserData.ID).Scan(&Profile.UserData.FollowingCount)
	// if err != nil {
	// 	log.Println("Error getting following count:", err)
	// 	return nil, err
	// }

	// queryFollowersData := `SELECT id, follower_id, followed_id, status, created_at 
	//           FROM followers 
	//           WHERE followed_id = ? AND status = 'accepted'`
	// rows, err = s.db.Query(queryFollowersData, Profile.UserData.ID)
	// if err != nil {
	// 	log.Println("Error getting followers data:", err)
	// 	return nil, err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var follow UserAccount
	// 	err := rows.Scan(&follow.ID, &follow.First_name, &follow.Nickname, &follow.Avatar)
	// 	if err != nil {
	// 		log.Println("Error scanning followers data:", err)
	// 		return nil, err
	// 	}
	// 	Profile.UserData.Followers = append(Profile.UserData.Followers, follow)
	// }

	// queryFollowingData := `
	// 	SELECT u.id, u.email, u.password, u.first_name, u.last_name, u.date_of_birth, 
	// 	       u.avatar, u.nickname, u.about_me, u.profile_type 
	// 	FROM followers f
	// 	JOIN users u ON f.followed_id = u.id
	// 	WHERE f.follower_id = ? AND f.status = 'accepted'
	// `
	// rows, err = s.db.Query(queryFollowingData, Profile.UserData.ID)
	// if err != nil {
	// 	log.Println("Error getting following data:", err)
	// 	return nil, err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var follow UserAccount
	// 	err := rows.Scan(&follow.ID, &follow.First_name, &follow.Nickname, &follow.Avatar)
	// 	if err != nil {
	// 		log.Println("Error scanning following data:", err)
	// 		return nil, err
	// 	}
	// 	Profile.UserData.Following = append(Profile.UserData.Following, follow)
	// }
	err := s.db.QueryRow("SELECT id, email, first_name, last_name, bio, avatar, nickname, profile_type FROM users WHERE id = ?" , Profile.UserData.ID).Scan(&Profile.UserData.ID, &Profile.UserData.Email, &Profile.UserData.First_name, &Profile.UserData.Last_name, &Profile.UserData.About_me, &Profile.UserData.Avatar, &Profile.UserData.Nickname, &Profile.UserData.Profile_type)
	if err != nil {
		log.Println("Error getting profile data:", err)
		return nil, err
	}

	postsRows, err := s.db.Query("SELECT id, user_id, content, image, privacy, created_at FROM posts WHERE user_id = ?", Profile.UserData.ID)
	if err != nil {
		log.Println("Error getting posts:", err)
		return nil, err
	}
	defer postsRows.Close()
	for postsRows.Next() {
		var post Post
		err := postsRows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image , &post.Privacy, &post.CreatedAt)
		if err != nil {
			log.Println("Error scanning posts:", err)
			return nil, err
		}
		Profile.Posts = append(Profile.Posts, post)
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
