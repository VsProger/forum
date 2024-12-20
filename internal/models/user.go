package models

type User struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"`
	GoogleID   string `json:"google_id,omitempty"`
	OAuthToken string `json:"oauth_token,omitempty"`
}
