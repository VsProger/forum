package oauth

import (
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubOauth2Config = oauth2.Config{
	ClientID:     "your-github-client-id",
	ClientSecret: "your-github-client-secret",
	RedirectURL:  "http://localhost:8080/auth/github/callback",
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
