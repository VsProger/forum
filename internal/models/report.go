package models

type Report struct {
	UserID       int    `json:"user_id"`
	UserName     string `json:"user_name"`
	UserEmail    string `json:"user_email"`
	PostID       int    `json:"post_id"`
	PostTitle    string `json:"post_title"`
	ReportReason string `json:"report_reason"`
}
