package models

type User struct {
	ID         int     `json:"id"`
	Email      string  `json:"email"`
	Username   string  `json:"username"`
	Password   string  `json:"password,omitempty"`
	GoogleID   *string `json:"google_id,omitempty"`
	GitHubID   *int64  `json:"github_id,omitempty"`
	OAuthToken string  `json:"oauth_token,omitempty"`
	Role       string  `json:"Role"`
}
