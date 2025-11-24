package domain

// PostUsersSetIsActiveJSONBody defines parameters for PostUsersSetIsActive.
type PostUsersSetIsActiveJSONBody struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
}

// User defines model for User.
type User struct {
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}
