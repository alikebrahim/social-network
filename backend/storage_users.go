package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (s *SQLiteStore) CreateUserAccount(user *UserAccount) (id int64, err error) {
	log.Print("1")
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print("bcrypt.GenerateFromPassword: ", err)
		return 0, err
	}
	log.Print("2")

	user.Password = string(hashedPass)
	qurey := `INSERT INTO users (email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, profile_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if user.Profile_type == "" {
		user.Profile_type = "public"
	}

	log.Print("user: ", user.Email)
	res, err := s.db.Exec(qurey, user.Email, user.Password, user.First_name, user.Last_name, user.Date_of_birth, user.Avatar, user.Nickname, user.About_me, user.Profile_type)
	if err != nil {
		log.Print("s.db.Exec: ", err)
		return 0, err
	}

	userID, err := res.LastInsertId()
	if err != nil {
		log.Print("res.LastInsertId: ", err)
		return 0, err
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		log.Print("generateSessionToken: ", err)
		return 0, err
	}

	qurey = `INSERT INTO sessions (user_id, session_token) VALUES (?, ?)`
	_, err = s.db.Exec(qurey, userID, sessionToken)
	if err != nil {
		log.Print("s.db.Exec: ", err)
		return 0, err
	}

	return userID, nil
}

func (s *SQLiteStore) EditUserAccount(usrAcc *UserAccount) error {
	// NOTE: should profile update be restricted (which attrs?)?
	query := `UPDATE users SET email = ?, password = ?, first_name = ?, last_name = ?, date_of_birth = ?, avatar = ?, nickname = ?, about_me = ?, profile_type = ? WHERE id = ?`
	_, err := s.db.Exec(query,
		usrAcc.Email,
		usrAcc.Password,
		usrAcc.First_name,
		usrAcc.Last_name,
		usrAcc.Date_of_birth,
		usrAcc.Avatar,
		usrAcc.Nickname,
		usrAcc.About_me,
		usrAcc.Profile_type,
		usrAcc.ID,
	)
	return err
}

func (s *SQLiteStore) DeleteUserAccount(usrID int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := s.db.Exec(query, usrID)
	return err
}

func (s *SQLiteStore) GetUserAccountByID(usrID int) (*UserAccount, error) {
	query := `SELECT id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, profile_type FROM users WHERE id = ?`
	row := s.db.QueryRow(query, usrID)

	userAccount := new(UserAccount)
	err := row.Scan(
		&userAccount.ID,
		&userAccount.Email,
		&userAccount.Password,
		&userAccount.First_name,
		&userAccount.Last_name,
		&userAccount.Date_of_birth,
		&userAccount.Avatar,
		&userAccount.Nickname,
		&userAccount.About_me,
		&userAccount.Profile_type,
	)
	if err != nil {
		return nil, err
	}
	return userAccount, nil
}

func (s *SQLiteStore) GetUserAccounts() ([]*UserAccount, error) {
	rows, err := s.db.Query("SELECT id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, profile_type FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*UserAccount{}
	for rows.Next() {
		userAccount := new(UserAccount)
		err := rows.Scan(
			&userAccount.ID,
			&userAccount.Email,
			&userAccount.Password,
			&userAccount.First_name,
			&userAccount.Last_name,
			&userAccount.Date_of_birth,
			&userAccount.Avatar,
			&userAccount.Nickname,
			&userAccount.About_me,
			&userAccount.Profile_type,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, userAccount)
	}
	return accounts, nil
}

func (s *SQLiteStore) AuthenticateUser(email string, password string) (string, error) {
	// Check if the user exists
	query := `SELECT id, password FROM users WHERE email = ?`
	var hashedPassword string
	var userID int64
	err := s.db.QueryRow(query, email).Scan(&userID, &hashedPassword)
	if err != nil {
		log.Print("s.db.QueryRow: ", err)
		return "", err
	}
	// Compare the stored hashed password with the one provided
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		log.Print("bcrypt.CompareHashAndPassword: ", err)
		return "", err
	}

	// Generate a new session token
	var sessionToken string
	query = `SELECT session_token FROM sessions WHERE user_id = ?`
	err = s.db.QueryRow(query, userID).Scan(&sessionToken)

	if err != nil { // No existing session, generate a new one
		sessionToken, err = generateSessionToken()
		if err != nil {
			log.Print("generateSessionToken: ", err)
			return "", err
		}

		// Insert the new session
		query = `INSERT INTO sessions (user_id, session_token) VALUES (?, ?)`
		_, err = s.db.Exec(query, userID, sessionToken)
		if err != nil {
			log.Print("s.db.Exec: ", err)
			return "", err
		}
	}

	return sessionToken, nil
}

func (s *SQLiteStore) GetSeesionToken(usrID int64) (string, error) {
	query := `SELECT session_token FROM sessions WHERE user_id = ?`
	row := s.db.QueryRow(query, usrID)

	var sessionToken string
	err := row.Scan(&sessionToken)
	if err != nil {
		return "", err
	}
	return sessionToken, nil
}

func (s *SQLiteStore) DeleteSession(sessionToken string) error {
	qurrey := `DELETE FROM sessions WHERE session_token = ?`
	_, err := s.db.Exec(qurrey, sessionToken)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) AddTestAccount() error {
	users := []struct {
		Email         string
		Password      string
		FirstName     string
		LastName      string
		DateOfBirth   string
		Avatar        string
		Nickname      string
		AboutMe       string
		ProfileType   string
	}{
		{"test1@example.com", "password123", "John", "Doe", "1990-01-01", "avatar1.png", "johndoe", "About John", "public"},
		{"test2@example.com", "password123", "Jane", "Smith", "1992-05-10", "avatar2.png", "janesmith", "About Jane", "private"},
		{"test3@example.com", "password123", "Mike", "Johnson", "1988-07-22", "avatar3.png", "mikej", "About Mike", "public"},
		{"test4@example.com", "password123", "Emily", "Davis", "1995-09-30", "avatar4.png", "emilyd", "About Emily", "public"},
		{"test5@example.com", "password123", "Robert", "Brown", "1985-11-15", "avatar5.png", "robb", "About Robert", "private"},
		{"test6@example.com", "password123", "Laura", "Wilson", "1998-04-25", "avatar6.png", "lauraw", "About Laura", "public"},
		{"test7@example.com", "password123", "Chris", "Miller", "1993-08-12", "avatar7.png", "chrism", "About Chris", "private"},
		{"test8@example.com", "password123", "Sophia", "Anderson", "1997-06-19", "avatar8.png", "sophiaa", "About Sophia", "public"},
		{"test9@example.com", "password123", "David", "Thomas", "1989-03-14", "avatar9.png", "davidth", "About David", "public"},
		{"test10@example.com", "password123", "Olivia", "Martinez", "1994-02-27", "avatar10.png", "oliviam", "About Olivia", "private"},
		{"test11@example.com", "password123", "Daniel", "Garcia", "1991-11-05", "avatar11.png", "danielg", "About Daniel", "public"},
		{"test12@example.com", "password123", "Ella", "Rodriguez", "1996-07-08", "avatar12.png", "ellar", "About Ella", "public"},
		{"test13@example.com", "password123", "James", "Hernandez", "1990-09-21", "avatar13.png", "jamesh", "About James", "private"},
		{"test14@example.com", "password123", "Ava", "Lopez", "1999-01-17", "avatar14.png", "avalo", "About Ava", "public"},
		{"test15@example.com", "password123", "William", "Clark", "1987-12-09", "avatar15.png", "williamc", "About William", "public"},
	}

	var userIDs []struct {
		ID          int64
		ProfileType string
	}

	// Insert Users into DB
	for _, user := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		query := `INSERT INTO users (email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, profile_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		res, err := s.db.Exec(query, user.Email, string(hashedPassword), user.FirstName, user.LastName, user.DateOfBirth, user.Avatar, user.Nickname, user.AboutMe, user.ProfileType)
		if err != nil {
			log.Print("s.db.Exec (users): ", err)
			return err
		}
		userID, err := res.LastInsertId()
		if err != nil {
			return err
		}
		userIDs = append(userIDs, struct {
			ID          int64
			ProfileType string
		}{userID, user.ProfileType})
	}

	// Insert Followers (Privacy Aware)
	for _, follower := range userIDs {
		for _, followee := range userIDs {
			if follower.ID != followee.ID { // Prevent self-following
				status := "accepted"
				if followee.ProfileType == "private" {
					status = "pending"
				}
				query := `INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)`
				_, err := s.db.Exec(query, follower.ID, followee.ID, status)
				if err != nil {
					log.Print("s.db.Exec (followers): ", err)
					return err
				}
			}
		}
	}

	// Insert Posts & Handle Privacy Settings
	for _, user := range userIDs {
		posts := []struct {
			Content  string
			Image    string
			Privacy  string
		}{
			{"Public post content", "public_image.jpg", "public"},
			{"Private post content", "private_image.jpg", "private"},
		}

		for _, post := range posts {
			query := `INSERT INTO posts (user_id, content, image, privacy, created_at) VALUES (?, ?, ?, ?, datetime('now'))`
			res, err := s.db.Exec(query, user.ID, post.Content, post.Image, post.Privacy)
			if err != nil {
				log.Print("s.db.Exec (posts): ", err)
				return err
			}
			postID, err := res.LastInsertId()
			if err != nil {
				return err
			}

			// If post is private, assign specific viewers
			if post.Privacy == "private" {
				for _, viewer := range userIDs {
					if viewer.ID != user.ID { // Exclude creator from restrictions
						query := `INSERT INTO post_visibility (post_id, viewer_id) VALUES (?, ?)`
						_, err := s.db.Exec(query, postID, viewer.ID)
						if err != nil {
							log.Print("s.db.Exec (post_visibility): ", err)
							return err
						}
					}
				}
			}
		}
	}

	log.Print("Test accounts, posts, and visibility settings added successfully!")
	return nil
}



func (s *SQLiteStore) GetUserIdBySession(session string) (int64, error) {
	log.Print("session: ", session)
	query := `SELECT user_id FROM sessions WHERE session_token = ?`
	var userID int64
	err := s.db.QueryRow(query, session).Scan(&userID)
	if err != nil {
		return 0, err
	}
	log.Print("userID: ", userID)
	return userID, nil
}
