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
	query := `UPDATE account SET email = ?, password = ?, first_name = ?, last_name = ?, date_of_birth = ?, avatar = ?, nickname = ?, about_me = ?, profile_type = ?, bio = ? WHERE id = ?`
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
	)
	return err
}

func (s *SQLiteStore) DeleteUserAccount(usrID int) error {
	query := `DELETE FROM account WHERE id = ?`
	_, err := s.db.Exec(query, usrID)
	return err
}

func (s *SQLiteStore) GetUserAccountByID(usrID int) (*UserAccount, error) {
	//NOTE: edit query
	query := `SELECT email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, profile_type FROM account WHERE id = ?`
	row := s.db.QueryRow(query, usrID)

	userAccount := new(UserAccount)
	err := row.Scan(
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
	//NOTE: edit query
	rows, err := s.db.Query("SELECT email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, profile_type FROM account")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*UserAccount{}
	for rows.Next() {
		userAccount := new(UserAccount)
		err := rows.Scan(
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
	