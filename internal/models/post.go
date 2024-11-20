package models

import (
	"time"
)

type Post struct {
	ID           int
	AuthorID     int
	Title        string
	Text         string
	LikeCount    int
	DislikeCount int
	CreationTime time.Time
	CategoryID   int
}
