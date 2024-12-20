package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	// Конфигурация OAuth 2.0 для Google
	OAuth2Config = oauth2.Config{
		ClientID:     "474394525572-vj65k8l3fnv0p0pp1i0c2ve31bnu137f.apps.googleusercontent.com", // Получите эти данные из Google Developer Console
		ClientSecret: "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",                                      // Получите эти данные из Google Developer Console
		RedirectURL:  "http://localhost:8081/auth/google/callback",                               // URL перенаправления
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	// Для защиты от CSRF атак генерируем случайную строку
	OAuth2StateString = "random"
)

func GetGoogleOAuth2Config() oauth2.Config {
	return OAuth2Config
}

func GetGoogleOAuth2State() string {
	return OAuth2StateString
}

// package oauth

// import (
// 	"net/http"

// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/google"
// )

// var GoogleOauth2Config = oauth2.Config{
// 	ClientID:     "474394525572-vj65k8l3fnv0p0pp1i0c2ve31bnu137f.apps.googleusercontent.com",
// 	ClientSecret: "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",
// 	RedirectURL:  "http://localhost:8081/auth/google/callback",
// 	Scopes:       []string{"email", "profile"},
// 	Endpoint:     google.Endpoint,
// }

// func GoogleAuthURL() string {
// 	return GoogleOauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline)
// }

// func GoogleCallback(r *http.Request) (*oauth2.Token, error) {
// 	code := r.URL.Query().Get("code")
// 	return GoogleOauth2Config.Exchange(r.Context(), code)
// }
