package storage

import (
	"socialNetwork/pkg/domain/auth"
	"socialNetwork/pkg/domain/chat"
	"socialNetwork/pkg/domain/following"
	"socialNetwork/pkg/domain/groups"
	"socialNetwork/pkg/domain/notifications"
	"socialNetwork/pkg/domain/posts"
	"socialNetwork/pkg/domain/profile"
)

// Storage defines the interface for all storage operations
type Storage interface {
	// Initialize storage
	Init() error


	// Auth operations
	CreateUserAccount(*auth.UserAccount) (id int64, err error)
	AuthenticateUser(string, string) (string, error)
	DeleteSession(string) error
	GetUserIdBySession(string) (int64, error)

	// Following operations
	CreateFollowRequest(following.FollowRequest) error
	AcceptFollowRequest(int64, int64) error
	DeleteFollowRequest(int64, int64) error
	GetFollowRequests(int64) ([]auth.UserAccount, error)
	IsFollowing(int64, int64) (bool, error)

	// Posts operations
	GetPostByID(int64, int64) ([]posts.Post, error)
	CreatePost(posts.Post) (id int64, err error)
	IsPostOwner(int64, int64) (bool, error)
	EditPost(posts.Post) error
	DeletePost(int64) error
	CreateComment(posts.Comment) (err error)
	CreateLike(posts.Like) error
	RemoveLikes(posts.Like) error
	CanUserSeePost(int64, int64) (bool, error)

	// Groups operations
	CreateGroup(userID int64, group *groups.CreateGroupRequest) (int64, error)
	GetGroup(groupID int64) (*groups.Group, error)
	GetGroups() ([]*groups.Group, error)
	SearchGroups(query string) ([]*groups.Group, error)

	// Group membership
	InviteToGroup(groupID, inviterID, inviteeID int64) error
	GetGroupInvites(userID int64) ([]*groups.Group, error)
	AcceptGroupInvite(groupID, userID int64) error
	RejectGroupInvite(groupID, userID int64) error
	RequestJoinGroup(groupID, userID int64) error
	GetGroupJoinRequests(groupID int64) ([]*groups.GroupMember, error)
	AcceptGroupJoinRequest(groupID, userID int64) error
	RejectGroupJoinRequest(groupID, userID int64) error
	IsGroupMember(groupID, userID int64) (bool, error)
	IsGroupCreator(groupID, userID int64) (bool, error)
	GetGroupMembers(groupID int64) ([]*groups.GroupMember, error)

	// Group events
	CreateGroupEvent(event *groups.GroupEvent) (int64, error)
	GetGroupEvents(groupID int64) ([]*groups.GroupEvent, error)
	RespondToEvent(eventID, userID int64, response string) error
	GetEventResponses(eventID int64) ([]*groups.EventResponse, error)

	// Group posts
	CreateGroupPost(groupID, userID int64, post *posts.Post) (int64, error)
	GetGroupPosts(groupID int64) ([]*posts.Post, error)

	// Chat operations
	SaveChat(chat *chat.Chat) (int64, error)
	GetUserChats(userID int64) ([]*chat.Chat, error)
	GetChatHistory(userID1, userID2 int64) ([]*chat.Chat, error)

	// Group chat
	SaveGroupChat(chat *chat.GroupChat) (int64, error)
	GetGroupChatHistory(groupID int64) ([]*chat.GroupChat, error)

	// Profile operations
	GetProfileData(*profile.Profile) (*profile.Profile, error)
	SetProfilePrivacy(profile.Profile, string) error
	
	// Notification operations
	CreateNotification(*notifications.Notification) error
	GetNotifications(userID int64, limit, offset int) (notifications.NotificationsList, error)
	MarkNotificationAsRead(notificationID, userID int64) error
	MarkAllNotificationsAsRead(userID int64) error
	GetUnreadNotificationCount(userID int64) (int, error)
}