package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VsProger/snippetbox/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
	GetUserByID(id int) (models.User, error)
	DeleteSession(token string) error
	GetUserByGoogleID(googleID string) (models.User, error)
	UpdateUserWithGoogleData(id string) error
	GetUserFromGoogleToken(token string) (models.User, error)
	CreateGoogleUser(user models.User) error
	CreateGithubUser(user models.User) error
	GetUserFromGitHubToken(token string) (models.User, error)
	UpdateUserWithGitHubData(user models.User) error
	GetUserByEmailGithub(email string) (models.User, error)
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		DB: db,
	}
}

func (auth *AuthRepo) CreateUser(user models.User) error {
	query := `INSERT INTO User (Username, Email, Password, Role) VALUES ($1, $2, $3, 'user')`
	_, err := auth.DB.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}
	return nil
}

func (auth *AuthRepo) CreateGoogleUser(user models.User) error {
	query := `INSERT INTO User (Username, Email, Password, GoogleID, Role) VALUES ($1, $2, $3, $4, 'user')`
	_, err := auth.DB.Exec(query, user.Username, user.Email, user.Password, user.GoogleID)
	if err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}
	return nil
}

func (auth *AuthRepo) CreateGithubUser(user models.User) error {
	query := `INSERT INTO User (Username, Email, Password, GitHubID, Role) VALUES ($1, $2, $3, $4, 'user')`
	_, err := auth.DB.Exec(query, user.Username, user.Email, user.Password, user.GitHubID)
	if err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}
	return nil
}

func (auth *AuthRepo) GetUserByToken(token string) (models.User, error) {
	query := `SELECT u.ID, u.Email, u.Username, u.Password, u.Role
	        FROM Session INNER JOIN User u
			ON u.ID = Session.UserID
			WHERE Session.Token = ?`
	var user models.User

	if err := auth.DB.QueryRow(query, token).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Role); err != nil {
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

func (r *AuthRepo) GetUserByEmailGithub(email string) (models.User, error) {
	query := `SELECT ID, Username, Email, Password, GitHubID FROM User WHERE Email = ?`

	row := r.DB.QueryRow(query, email)
	user := models.User{}

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.GitHubID)
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

func (auth *AuthRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	query := `SELECT * FROM User WHERE ID = ?`

	if err := auth.DB.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (auth *AuthRepo) DeleteSession(token string) error {
	query := `DELETE FROM Session WHERE Token = ?`

	_, err := auth.DB.Exec(query, token)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepo) GetUserByGoogleID(googleID string) (models.User, error) {
	var user models.User
	query := `SELECT ID, Username, Email, GoogleID FROM User WHERE GoogleID = ?`

	row := r.DB.QueryRow(query, googleID)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.GoogleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user with GoogleID %s not found", googleID)
		}
		return user, fmt.Errorf("unable to fetch user by GoogleID: %w", err)
	}
	return user, nil
}

func (r *AuthRepo) GetUserByGithubID(githubID string) (models.User, error) {
	var user models.User
	query := `SELECT ID, Username, Email, GithubID FROM User WHERE GithubID = ?`

	row := r.DB.QueryRow(query, githubID)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.GitHubID)
	if err != nil {
		if err == sql.ErrNoRows {

			return user, fmt.Errorf("user with GithubID %s not found", githubID)
		}
		return user, fmt.Errorf("unable to fetch user by GithubID: %w", err)
	}
	return user, nil
}

func (auth *AuthRepo) UpdateUserWithGoogleData(id string) error {
	query := `UPDATE User SET GoogleID = $1 WHERE GoogleID = ''`
	_, err := auth.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("unable to update user with Google data: %w", err)
	}
	return nil
}

func (repo *AuthRepo) GetUserFromGoogleToken(token string) (models.User, error) {
	// Initialize Google OAuth2 config
	googleOauth2Config := oauth2.Config{
		ClientID:     "474394525572-pbrh9edm251u9d04e0l9l7qtqiq217bg.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8081/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Use Google OAuth2 config to create a client and fetch user info
	client := googleOauth2Config.Client(context.Background(), &oauth2.Token{AccessToken: token})
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return models.User{}, err
	}
	defer resp.Body.Close()

	// Decode the response into a user struct
	var googleUser struct {
		ID       string `json:"id"`
		Username string `json:"name"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return models.User{}, err
	}

	// Find the user by their Google ID (or email, depending on your implementation)
	user, err := repo.GetUserByGoogleID(googleUser.ID)
	if err != nil {
		return models.User{}, err
	}

	// Return the user information
	return user, nil
}

func (repo *AuthRepo) GetUserFromGitHubToken(token string) (models.User, error) {

	const githubUserInfoURL = "https://api.github.com/user"

	client := &http.Client{}
	req, err := http.NewRequest("GET", githubUserInfoURL, nil)
	if err != nil {
		return models.User{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return models.User{}, err
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID       string `json:"id"`
		Username string `json:"login"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return models.User{}, err
	}

	user, err := repo.GetUserByGithubID(githubUser.ID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (auth *AuthRepo) UpdateUserWithGitHubData(user models.User) error {
	query := `UPDATE User SET GitHubID = $1 WHERE GitHubID IS NULL`
	_, err := auth.DB.Exec(query, user.GitHubID)
	if err != nil {
		return fmt.Errorf("unable to update user with GitHub data: %w", err)
	}
	return nil
}

func (r *AuthRepo) GetUserRole(userid string) (models.User, error) {
	var user models.User
	query := `SELECT Role FROM User WHERE ID = ?`

	row := r.DB.QueryRow(query, userid)
	err := row.Scan(&user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		return user, fmt.Errorf("unable to fetch user by userid: %w", err)
	}
	return user, nil
}
