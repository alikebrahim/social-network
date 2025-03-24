package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// NOTE: migration handling
// goose -dir migrationsDir sqlite3 up
// goose -dir migrationsDir sqlite3 down

type Storage interface {
	// test data
	AddTestAccount() error
	// USER ACCOUNTS
	CreateUserAccount(*UserAccount) (id int64, err error)
	GetSeesionToken(int64) (string, error)
	AuthenticateUser(string, string) (string, error)
	DeleteSession(string) error
	GetUserIdBySession(string) (int64, error)

	// follow requests
	CreateFollowRequest(follow) error
	AcceptFollowRequest(int64, int64) error
	DeleteFollowRequest(int64, int64) error
	GetFollowRequests(int64) ([]UserAccount, error)
	
	

	// DeleteUserAccount(int) error
	// EditUserAccount(*UserAccount) error
	// GetUserAccountByID(int) (*UserAccount, error)
	// GetUserAccounts() ([]*UserAccount, error)

	// POSTS
	CreatePost(*Post) (id int64, err error)
	GetPostByID(int64, int64) ([]Post, error)
	CreatePost(Post) (id int64, err error)
	IsPostOwner(int64, int64) (bool, error)
	EditPost(Post) error
	DeletePost(int64) error
	CreateComment(Comment) (err error)
	CreateLike(like) error
	RemoveLikes(like) error

	// GROUPS
	CreateGroup(userID int64, group *CreateGroupRequest) (int64, error)
	GetGroup(groupID int64) (*Group, error)
	GetGroups() ([]*Group, error)
	SearchGroups(query string) ([]*Group, error)
	
	// GROUP MEMBERSHIP
	InviteToGroup(groupID, inviterID, inviteeID int64) error
	GetGroupInvites(userID int64) ([]*Group, error)
	AcceptGroupInvite(groupID, userID int64) error
	RejectGroupInvite(groupID, userID int64) error
	
	RequestJoinGroup(groupID, userID int64) error
	GetGroupJoinRequests(groupID int64) ([]*GroupMember, error)
	AcceptGroupJoinRequest(groupID, userID int64) error
	RejectGroupJoinRequest(groupID, userID int64) error
	
	IsGroupMember(groupID, userID int64) (bool, error)
	IsGroupCreator(groupID, userID int64) (bool, error)
	GetGroupMembers(groupID int64) ([]*GroupMember, error)
	
	// GROUP EVENTS
	CreateGroupEvent(event *GroupEvent) (int64, error)
	GetGroupEvents(groupID int64) ([]*GroupEvent, error)
	RespondToEvent(eventID, userID int64, response string) error
	GetEventResponses(eventID int64) ([]*EventResponse, error)
	
	// GROUP POSTS
	CreateGroupPost(groupID, userID int64, post *Post) (int64, error)
	GetGroupPosts(groupID int64) ([]*Post, error)
	
	// CHAT
	SaveChat(chat *Chat) (int64, error)
	GetUserChats(userID int64) ([]*Chat, error)
	GetChatHistory(userID1, userID2 int64) ([]*Chat, error)
	
	// GROUP CHAT
	SaveGroupChat(chat *GroupChat) (int64, error)
	GetGroupChatHistory(groupID int64) ([]*GroupChat, error)
	
	// profile
	GetProfileData(*Profile) (*Profile, error)
	SetProfilePrivacy(Profile, string) error
	isFollowing(int64, int64) (bool, error)
}

type SQLiteStore struct {
	db *sql.DB
}

func (s *SQLiteStore) Init() error {
	return s.db.Ping()
}

func NewSQLiteStore() (*SQLiteStore, error) {
	connStr := "./pkg/db/sqlite/database.sql"

	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{
		db: db,
	}, nil
}

// CreateGroup creates a new group
func (s *SQLiteStore) CreateGroup(userID int64, req *CreateGroupRequest) (int64, error) {
	query := `INSERT INTO groups (creator_id, title, description, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?)`
	
	now := time.Now()
	result, err := s.db.Exec(query, userID, req.Title, req.Description, now, now)
	if err != nil {
		return 0, err
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Add creator as a member with accepted status
	memberQuery := `INSERT INTO group_members (group_id, user_id, status, created_at) 
                   VALUES (?, ?, ?, ?)`
	_, err = s.db.Exec(memberQuery, groupID, userID, "accepted", now)
	if err != nil {
		return 0, err
	}

	return groupID, nil
}

// GetGroup retrieves a group by ID
func (s *SQLiteStore) GetGroup(groupID int64) (*Group, error) {
	query := `SELECT id, creator_id, title, description, created_at, updated_at 
              FROM groups WHERE id = ?`
	
	var group Group
	err := s.db.QueryRow(query, groupID).Scan(
		&group.ID,
		&group.CreatorID,
		&group.Title,
		&group.Description,
		&group.CreatedAt,
		&group.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return &group, nil
}

// GetGroups retrieves all groups
func (s *SQLiteStore) GetGroups() ([]*Group, error) {
	// Check if the database has groups table
	var count int
	err := s.db.QueryRow(`SELECT count(name) FROM sqlite_master 
                          WHERE type='table' AND name='groups'`).Scan(&count)
	if err != nil {
		return nil, err
	}
	
	// If groups table doesn't exist, return empty slice
	if count == 0 {
		return []*Group{}, nil
	}
	
	// Check if there are any groups
	err = s.db.QueryRow(`SELECT COUNT(*) FROM groups`).Scan(&count)
	if err != nil {
		return nil, err
	}
	
	// If no groups exist, return empty slice
	if count == 0 {
		return []*Group{}, nil
	}
	
	// Query all groups - simple query with WHERE clause to filter NULLs
	query := `SELECT id, creator_id, title, description, created_at, updated_at
              FROM groups 
              WHERE id IS NOT NULL
              ORDER BY created_at DESC`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var groups []*Group
	for rows.Next() {
		var group Group
		err := rows.Scan(
			&group.ID,
			&group.CreatorID,
			&group.Title,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return groups, nil
}

// SearchGroups searches groups by title or description
func (s *SQLiteStore) SearchGroups(query string) ([]*Group, error) {
	// Check if the database has groups table
	var count int
	err := s.db.QueryRow(`SELECT count(name) FROM sqlite_master 
                          WHERE type='table' AND name='groups'`).Scan(&count)
	if err != nil {
		return nil, err
	}
	
	// If groups table doesn't exist, return empty slice
	if count == 0 {
		return []*Group{}, nil
	}
	
	// Check if there are any groups
	err = s.db.QueryRow(`SELECT COUNT(*) FROM groups`).Scan(&count)
	if err != nil {
		return nil, err
	}
	
	// If no groups exist, return empty slice
	if count == 0 {
		return []*Group{}, nil
	}
	
	// Search groups - simple query with WHERE clause to filter NULLs
	sqlQuery := `SELECT id, creator_id, title, description, created_at, updated_at
                FROM groups 
                WHERE (title LIKE ? OR description LIKE ?) AND id IS NOT NULL
                ORDER BY created_at DESC`
	
	searchParam := "%" + query + "%"
	rows, err := s.db.Query(sqlQuery, searchParam, searchParam)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var groups []*Group
	for rows.Next() {
		var group Group
		err := rows.Scan(
			&group.ID,
			&group.CreatorID,
			&group.Title,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return groups, nil
}

// InviteToGroup invites a user to a group
func (s *SQLiteStore) InviteToGroup(groupID, inviterID, inviteeID int64) error {
	// Check if inviter is a member of the group
	isMember, err := s.IsGroupMember(groupID, inviterID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrUnauthorized
	}
	
	// Check if invitee is already a member or has a pending invitation
	query := `SELECT COUNT(*) FROM group_members 
              WHERE group_id = ? AND user_id = ?`
	
	var count int
	err = s.db.QueryRow(query, groupID, inviteeID).Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return ErrAlreadyExists
	}
	
	// Create invitation (pending membership)
	insertQuery := `INSERT INTO group_members (group_id, user_id, status, created_at) 
                   VALUES (?, ?, ?, ?)`
	_, err = s.db.Exec(insertQuery, groupID, inviteeID, "pending", time.Now())
	if err != nil {
		return err
	}
	
	return nil
}

// IsGroupMember checks if a user is a member of a group
func (s *SQLiteStore) IsGroupMember(groupID, userID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM group_members 
              WHERE group_id = ? AND user_id = ? AND status = 'accepted'`
	
	var count int
	err := s.db.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// IsGroupCreator checks if a user is the creator of a group
func (s *SQLiteStore) IsGroupCreator(groupID, userID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM groups WHERE id = ? AND creator_id = ?`
	
	var count int
	err := s.db.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}
