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
