package auth

type UserAccount struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	First_name    string `json:"first_name"`
	Last_name     string `json:"last_name"`
	Date_of_birth string `json:"date_of_birth"`
	Avatar        string `json:"avatar"`
	Nickname      string `json:"nickname"`
	About_me      string `json:"about_me"`
	Profile_type  string `json:"profile_type"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	First_name    string `json:"first_name"`
	Last_name     string `json:"last_name"`
	Date_of_birth string `json:"date_of_birth,omitempty"`
	Profile_type  string `json:"profile_type"`
}

type Session struct {
	ID             string        `json:"id"`
	UserID         int64         `json:"user_id"`
	CreatedAt      string        `json:"created_at"`
	ExpiresAt      string        `json:"expires_at"`
	FollowCount    int64         `json:"followCount"`
	FollowingCount int64         `json:"followingCount"`
	Followers      []UserAccount `json:"followers"`
	Following      []UserAccount `json:"following"`
}