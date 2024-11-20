package service

import (
	repo "github.com/VsProger/snippetbox/internal/repository"
	authService "github.com/VsProger/snippetbox/internal/service/auth"

	postService "github.com/VsProger/snippetbox/internal/service/posts"
)

type Service struct {
	authService.Auth
	postService.PostService
	// 	filter.Filter
}

func NewService(repo *repo.Repository) *Service {
	return &Service{
		Auth:        authService.NewAuthService(repo.Authorization),
		PostService: postService.NewPostService(repo.Posts),
		// Filter:      filter.NewFilterService(repo.Filter),
	}
}
