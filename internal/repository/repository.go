package repository

import (
	"database/sql"
	"forum/internal/repository/auth"
	"forum/internal/repository/filter"
	"forum/internal/repository/posts"
)

type Repository struct {
	auth.Authorization
	posts.Posts
	filter.Filter
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		Authorization: auth.NewAuthRepo(db),
		Posts:         posts.NewPostRepo(db),
		Filter:        filter.NewFilterRepo(db),
	}
}
