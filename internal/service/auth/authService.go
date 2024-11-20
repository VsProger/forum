package authService

import (
	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/internal/repository/auth"
)

type Auth interface {
	CreateUser(user models.User) error
}

type AuthService struct {
	repo auth.Authorization
}

func NewAuthService(repo auth.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) CreateUser(user models.User) error {
	var err error
	if err != nil {
		return err
	}
	return a.repo.CreateUser(user)
}
