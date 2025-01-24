package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/internal/repository/auth"
	"github.com/VsProger/snippetbox/pkg"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type Auth interface {
	CreateUser(user models.User) error
	GetUserByToken(token string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByEmailGithub(email string) (models.User, error)
	CheckUser(user *models.User) error
	GetUserByUsername(username string) (models.User, error)
	CheckPassword(user models.User) error
	SetSession(user *models.User) (string, error)
	DeleteSession(token string) error
	CreateUserFromOAuth(token *oauth2.Token) (models.User, error)
	GetUserByGoogleID(token string) (models.User, error)
	UpdateUserWithGoogleData(token string) error
	CreateUserGoogle(user models.User) error
	CreateUserGitHub(user models.User) error
	UpdateUserWithGitHubData(token string) error
}

var googleOauth2Config = oauth2.Config{
	ClientID:     "474394525572-vj65k8l3fnv0p0pp1i0c2ve31bnu137f.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:8081/auth/google/callback",
	Scopes:       []string{"email", "profile"},
	Endpoint:     google.Endpoint,
}

var githubOauth2Config = oauth2.Config{
	ClientID:     "Ov23liop6ipn43yQRXfw",                       // Replace with your GitHub Client ID
	ClientSecret: "52faf14f32e6efe5b76741d5bd91c485e848c392",   // Replace with your GitHub Client Secret
	RedirectURL:  "http://localhost:8081/auth/github/callback", // Adjust based on your environment
	Scopes:       []string{"read:user", "user:email"},
	Endpoint:     github.Endpoint,
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

func (a *AuthService) CreateUserGoogle(user models.User) error {
	var err error
	if err != nil {
		return err
	}
	return a.repo.CreateGoogleUser(user)
}

func (a *AuthService) CreateUserGitHub(user models.User) error {
	var err error
	if err != nil {
		return err
	}
	return a.repo.CreateGithubUser(user)
}

func (a *AuthService) GetUserByToken(token string) (models.User, error) {
	user, err := a.repo.GetUserByToken(token)
	if err != nil {
		return user, nil
	}
	return user, nil
}

func (a *AuthService) GetUserByGoogleID(token string) (models.User, error) {
	user, err := a.repo.GetUserByGoogleID(token)
	if err != nil {
		return user, nil
	}
	return user, nil
}

func (a *AuthService) UpdateUserWithGoogleData(token string) error {
	// Get the user by the token (likely the Google ID)
	user, err := a.repo.GetUserFromGoogleToken(token) // Assuming this method retrieves the user based on the Google token
	if err != nil {
		return err
	}

	// Update user with the new Google data
	err = a.repo.UpdateUserWithGoogleData(*user.GoogleID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) GetUserByEmail(email string) (models.User, error) {
	user, err := a.repo.GetUserByEmail(email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (a *AuthService) GetUserByEmailGithub(email string) (models.User, error) {
	user, err := a.repo.GetUserByEmailGithub(email)
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

func (a *AuthService) DeleteSession(token string) error {
	return a.repo.DeleteSession(token)
}

func (a *AuthService) CreateUserFromOAuth(token *oauth2.Token) (models.User, error) {
	var user models.User
	var userInfoURL string

	// OAuth2 configuration for Google
	googleOauth2Config := oauth2.Config{
		ClientID:     "474394525572-vj65k8l3fnv0p0pp1i0c2ve31bnu137f.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8081/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	// Determine the provider and set the appropriate user info URL
	if token.Extra("id_token") != nil { // Google
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	} else { // Assume GitHub
		userInfoURL = "https://api.github.com/user"
	}

	// Create an OAuth2 client
	var oauth2Config oauth2.Config
	if userInfoURL == "https://www.googleapis.com/oauth2/v2/userinfo" {
		oauth2Config = googleOauth2Config
	} else {
		oauth2Config = githubOauth2Config
	}
	client := oauth2Config.Client(oauth2.NoContext, token)

	// Make a request to the user info endpoint
	resp, err := client.Get(userInfoURL)
	if err != nil {
		log.Println("Failed to fetch user info:", err)
		return user, err
	}
	defer resp.Body.Close()

	// Parse the response
	if userInfoURL == "https://www.googleapis.com/oauth2/v2/userinfo" {
		var googleUser struct {
			ID       string `json:"id"`
			Username string `json:"name"`
			Email    string `json:"email"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
			log.Println("Failed to parse Google user info:", err)
			return user, err
		}

		// Populate the user model
		user.Email = googleUser.Email
		user.Username = googleUser.Username
	} else {
		var githubUser struct {
			ID       string `json:"id"`
			Username string `json:"login"`
			Email    string `json:"email"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
			log.Println("Failed to parse GitHub user info:", err)
			return user, err
		}

		// Populate the user model
		user.Email = githubUser.Email
		user.Username = githubUser.Username
	}

	// Check if the user already exists based on the email
	existingUser, err := a.repo.GetUserByEmail(user.Email)
	if err != nil && err != models.ErrUserNotFound {
		// Handle error if getting the user fails for any reason
		log.Println("Error checking existing user:", err)
		return user, err
	}

	// If the user already exists, return the existing user
	if existingUser.ID != 0 {
		return existingUser, nil
	}

	// Create a new user if one doesn't exist
	if err := a.repo.CreateUser(user); err != nil {
		log.Println("Error creating user:", err)
		return user, err
	}

	// Pass the token to the repository method with the oauth2Config
	userFromRepo, err := a.repo.GetUserFromGoogleToken(token.AccessToken)
	if err != nil {
		log.Println("Error retrieving user from Google token:", err)
		return user, err
	}

	return userFromRepo, nil
}

func (a *AuthService) UpdateUserWithGitHubData(token string) error {
	// Get the user by the token (likely the GitHub ID)
	user, err := a.repo.GetUserFromGitHubToken(token) // Assuming this method retrieves the user based on the GitHub token
	if err != nil {
		log.Println(err)
		return err
	}

	// Update user with the new GitHub data
	err = a.repo.UpdateUserWithGitHubData(user)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
