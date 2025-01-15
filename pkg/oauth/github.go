package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VsProger/snippetbox/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubOauth2Config = oauth2.Config{
	ClientID:     "Ov23liop6ipn43yQRXfw",
	ClientSecret: "52faf14f32e6efe5b76741d5bd91c485e848c392",
	RedirectURL:  "http://localhost:8081/auth/github/callback",
	Scopes:       []string{"user:email"},
	Endpoint:     github.Endpoint,
}

func GitHubAuthURL() string {
	return githubOauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline)
}

func GitHubCallback(r *http.Request) (*oauth2.Token, error) {
	code := r.URL.Query().Get("code")
	return githubOauth2Config.Exchange(r.Context(), code)
}

func GetGitHubUserInfo(accessToken string) (*models.User, error) {
	const githubUserInfoURL = "https://api.github.com/user"

	client := &http.Client{}
	req, err := http.NewRequest("GET", githubUserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	// Устанавливаем авторизационный заголовок
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch GitHub user info: %s", resp.Status)
	}

	// Декодируем ответ в структуру
	var githubUser struct {
		GitHubID int64  `json:"id"`
		Username string `json:"login"`
		Email    string `json:"email,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, err
	}

	// Если email отсутствует, пытаемся получить его через дополнительный запрос
	if githubUser.Email == "" {
		email, err := GetGitHubUserEmail(accessToken)
		if err != nil {
			return nil, err
		}
		githubUser.Email = email
	}

	// Возвращаем информацию о пользователе
	return &models.User{
		GitHubID: &githubUser.GitHubID,
		Username: githubUser.Username,
		Email:    githubUser.Email,
	}, nil
}

func GetGitHubUserEmail(accessToken string) (string, error) {
	const githubUserEmailURL = "https://api.github.com/user/emails"

	client := &http.Client{}
	req, err := http.NewRequest("GET", githubUserEmailURL, nil)
	if err != nil {
		return "", err
	}

	// Устанавливаем авторизационный заголовок
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch GitHub user emails: %s", resp.Status)
	}

	// Декодируем ответ в структуру
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	// Проверяем, если есть primary email
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	return "", nil
}
