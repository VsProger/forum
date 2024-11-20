package service

import (
	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/internal/repository/auth"
)

type Auth interface {
	CreateUser(user models.User) error
	GetUserByToken(token string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	CheckUser(user *models.User) error
	GetUserByUsername(username string) (models.User, error)
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

func (a *AuthService) GetUserByToken(token string) (models.User, error) {
	user, err := a.repo.GetUserByToken(token)
	if err != nil {
		return user, nil
	}
	return user, nil
}

func (a *AuthService) GetUserByEmail(email string) (models.User, error) {
	user, err := a.repo.GetUserByEmail(email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (a *AuthService) GetUserByUSername(username string) (models.User, error) {
	user, err := a.repo.GetUserByUsername(username)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (a *AuthService) CheckUser(user *models.User) error {
	if err := pkg.ValidatePassword(user.Password); err != nil {
		return err
	}

	if err := pkg.ValidateUsername(user.Username); err != nil {
		return err
	}

	if err := pkg.ValidateEmail(user.Email); err != nil {
		return err
	}
	return nil
}
