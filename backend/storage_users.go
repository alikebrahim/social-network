package main

import "fmt"

func (s *SQLiteStore) CreateUserAccount(usrAcc *UserAccount) (id int64, err error) {
	query := `INSERT INTO users (email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, profile_type)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := s.db.Exec(query,
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
	if err != nil {
		fmt.Println("db exec: ", err)
		return 0, err
	}

	// ID for inserted usrAcc
	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	fmt.Printf("Inserted with ID: %d\n", id)
	return id, nil
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
