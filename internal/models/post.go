package models

import (
	"time"
)

type Post struct {
	ID           int
	AuthorID     int
	Title        string
	Text         string
	ImageURL     string
	LikeCount    int
	DislikeCount int
	Username     string
	CreationTime time.Time
	CategoryId   []int
	Comment      []Comment
	Categories   []Category
	Category     string
}

type Category struct {
	ID   int
	Name string
}

type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id,omitempty"` // Optional for comments
	CommentID int       `json:"comment_id,omitempty"`
	Type      string    `json:"type"` // e.g., "like", "dislike", "comment"
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	IsRead    bool      `json:"is_read"`
	Username  string    `json:"Username"`
}
