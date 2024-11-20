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
	GetUserByToken(token string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	DeleteSessionByUserID(userID int) error
	CreateSession(sessions models.Session) error
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

func (auth *AuthRepo) GetUserByToken(token string) (models.User, error) {
	query := `SELECT u.ID, u.Email, u.Username, u.Password
	        FROM Session INNER JOIN User u
			ON u.ID = Session.UserID
			WHERE Session.Token = ?`
	var user models.User

	if err := auth.DB.QueryRow(query, token).Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		return user, err
	}
	return user, nil
}

func (r *AuthRepo) GetUserByEmail(email string) (models.User, error) {
	query := `SELECT ID, Username, Email, Password FROM User WHERE Email = ?`

	row := r.DB.QueryRow(query, email)
	user := models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *AuthRepo) GetUserByUsername(username string) (models.User, error) {
	query := `SELECT ID, Username, Email, Password FROM User WHERE Username = ?`

	row := r.DB.QueryRow(query, username)
	user := models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (auth *AuthRepo) DeleteSessionByUserID(userID int) error {
	query := `DELETE FROM Session WHERE UserID = ?`
	_, err := auth.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (auth *AuthRepo) CreateSession(session models.Session) error {
	query := `INSERT INTO Session (UserID, Token, ExpTime) VALUES ($1, $2, $3)`

	_, err := auth.DB.Exec(query, session.UserID, session.Token, session.ExpTime)
	if err != nil {
		return err
	}
	return nil
}
