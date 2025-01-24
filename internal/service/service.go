package service

import (
	repo "github.com/VsProger/snippetbox/internal/repository"
	"github.com/VsProger/snippetbox/internal/service/admin"
	authService "github.com/VsProger/snippetbox/internal/service/auth"
	filter "github.com/VsProger/snippetbox/internal/service/filter"
	postService "github.com/VsProger/snippetbox/internal/service/posts"
)

type Service struct {
	authService.Auth
	postService.PostService
	filter.Filter
	admin.Admin
}

func NewService(repo *repo.Repository) *Service {
	return &Service{
		Auth:        authService.NewAuthService(repo.Authorization),
		PostService: postService.NewPostService(repo.Posts),
		Filter:      filter.NewFilterService(repo.Filter),
		Admin:       admin.NewAdminService(repo.Admin),
	}
}
