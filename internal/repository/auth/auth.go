package auth

import (
	"database/sql"

	"github.com/VsProger/snippetbox/internal/models"
)

type AuthRepo struct {
	DB *sql.DB
}

type Authorization interface {
	CreateUser(user models.User) error
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		DB: db,
	}
}

func (auth *AuthRepo) CreateUser(user models.User) error {
	query := `INSERT INTO User (Email, Username, Password) VALUES ($1, $2, $3)`

	_, err := auth.DB.Exec(query, user.Email, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}
