package service

import (
	"time"

	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/internal/repository/auth"
	"github.com/VsProger/snippetbox/pkg"
)

type Auth interface {
	CreateUser(user models.User) error
	GetUserByToken(token string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	CheckUser(user *models.User) error
	GetUserByUsername(username string) (models.User, error)
	CheckPassword(user models.User) error
	SetSession(user *models.User) (string, error)
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

func (a *AuthService) GetUserByUsername(username string) (models.User, error) {
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

func (a *AuthService) CheckPassword(user models.User) error {
	checkedUser, err := a.repo.GetUserByEmail(user.Email)
	if err != nil {
		return models.ErrUserNotFound
	}
	if !pkg.CheckPasswordHash(user.Password, checkedUser.Password) {
		return models.ErrInvalidPassword
	}
	return nil
}

func (a *AuthService) SetSession(user *models.User) (string, error) {
	a.repo.DeleteSessionByUserID(user.ID)
	token := pkg.GenerateToken()

	session := models.Session{
		UserID:  user.ID,
		Token:   token,
		ExpTime: time.Now().Add(3 * time.Hour),
	}
	if err := a.repo.CreateSession(session); err != nil {
		return "hui", err
	}
	return session.Token, nil
}
