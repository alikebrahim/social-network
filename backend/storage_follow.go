package main

import (
	"errors"
	"log"
)

func (s *SQLiteStore) CreateFollowRequest(follow follow) error {
	var profileType string
	err := s.db.QueryRow("SELECT profile_type FROM users WHERE id = ?", follow.FollowedID).Scan(&profileType)
	if err != nil {
		log.Print("s.db.QueryRow: ", err)
		return err
	}
	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ?", follow.FollowerID, follow.FollowedID).Scan(&count)
	if err != nil {
		log.Print("s.db.QueryRow: ", err)
		return err
	}
	if count > 0 {
		log.Print("already following")
		err := errors.New("already following")
		return err
	}

	status := "pending"
	if profileType == "public" {
		status = "accepted"
	}
	_, err = s.db.Exec("INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)", follow.FollowerID, follow.FollowedID, status)
	if err != nil {
		log.Print("s.db.Exec: ", err)
		return err
	}
	return nil
}

func (s *SQLiteStore) AcceptFollowRequest(followerId int64, followedId int64) error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'pending'", followerId, followedId).Scan(&count)
	if err != nil {
		log.Print("s.db.QueryRow: ", err)
		return err
	}
	if count == 0 {
		log.Print("no follow request")
		err := errors.New("no follow request")
		return err
	}

	_, err = s.db.Exec("UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND followed_id = ?", followerId, followedId)
	if err != nil {
		log.Print("s.db.Exec: ", err)
		return err
	}

	return nil
}

func (s *SQLiteStore) DeleteFollowRequest(followerId int64, followedId int64) error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'pending'", followerId, followedId).Scan(&count)
	if err != nil {
		log.Print("s.db.QueryRow: ", err)
		return err
	}
	if count == 0 {
		log.Print("no follow request")
		err := errors.New("no follow request")
		return err
	}
	_, err = s.db.Exec("DELETE FROM followers WHERE follower_id = ? AND followed_id = ?", followerId, followedId)
	if err != nil {
		log.Print("s.db.QueryRow: ", err)
		return err
	}
	return nil
}

func (s *SQLiteStore) GetFollowRequests(userId int64) ([]UserAccount, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, u.nickname, u.avatar
		FROM followers f
		JOIN users u ON f.followed_id = u.id
		WHERE f.follower_id = ? AND f.status = 'pending'
	`
	rows, err := s.db.Query(query, userId)
	if err != nil {
		log.Print("s.db.Query error: ", err)
		return nil, err
	}
	defer rows.Close()

	var followRequests []UserAccount
	for rows.Next() {
		var follower UserAccount
		err := rows.Scan(&follower.ID, &follower.First_name, &follower.Last_name, &follower.Nickname, &follower.Avatar)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		followRequests = append(followRequests, follower)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error during rows iteration:", err)
		return nil, err
	}

	log.Println("Follow requests:", followRequests)
	return followRequests, nil
}

func (s *SQLiteStore) isFollowing(requesterId int64,userId int64) (bool, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'accepted'", 
		requesterId, userId,
	).Scan(&count)
		if err != nil {
		log.Print("s.db.QueryRow: ", err)
		return false, err
	}
	log.Print("count: ", count)
	if count == 0 {
		return false, nil
	}
	return true, nil
}
