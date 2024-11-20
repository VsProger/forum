package service

import (
	repo "github.com/VsProger/snippetbox/internal/repository"
	authService "github.com/VsProger/snippetbox/internal/service/auth"
	"github.com/VsProger/snippetbox/internal/service/filter"
	"github.com/VsProger/snippetbox/internal/service/posts"
)

type Service struct {
	authService.Auth
	posts.PostService
	filter.Filter
}

func NewService(repo *repo.Repository) *Service {
	return &Service{
		Auth:        authService.NewAuthService(repo.Authorization),
		PostService: posts.NewPostService(repo.Posts),
		Filter:      filter.NewFilterService(repo.Filter),
	}
}
