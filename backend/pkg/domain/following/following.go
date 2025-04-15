package following

// FollowRequest represents a request to follow another user
type FollowRequest struct {
	FollowerID int64 `json:"follower_id"`
	FollowedID int64 `json:"followed_id"`
}

// FollowStatus represents the status of a follow relationship
type FollowStatus struct {
	IsFollowing bool `json:"is_following"`
}