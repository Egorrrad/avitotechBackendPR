package domain

// PostUsersSetIsActiveJSONBody defines parameters for PostUsersSetIsActive.
type PostUsersSetIsActiveJSONBody struct {
	IsActive bool   `json:"is_active"`
	UserID   string `json:"user_id"`
}

// User defines model for User.
type User struct {
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}
