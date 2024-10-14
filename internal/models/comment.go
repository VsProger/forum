package models

type Comment struct {
	ID           int
	Text         string
	PostID       int
	AuthorID     int
	LikeCount    int
	DislikeCount int
	Username     string
}
