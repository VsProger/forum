package repository

import (
	"database/sql"
	"github.com/VsProger/snippetbox/internal/repository/auth"
	// "github.com/VsProger/snippetbox/internal/repository/filter"
	// "github.com/VsProger/snippetbox/internal/repository/posts"
)

type Repository struct {
	auth.Authorization
	// posts.Posts
	// filter.Filter
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		Authorization: auth.NewAuthRepo(db),
	// 	Posts:         posts.NewPostRepo(db),
	// 	Filter:        filter.NewFilterRepo(db),
	// }
}
