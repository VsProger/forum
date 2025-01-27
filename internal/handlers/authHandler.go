package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/pkg/oauth"
	"golang.org/x/oauth2"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	nameFunction := "indexHandler"
	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusNotFound, nameFunction)
		return
	}
	if r.Method == http.MethodGet {
		var username string
		var role string
		var user models.User
		session, err := r.Cookie("session")
		if err == nil {
			user, err = h.service.GetUserByToken(session.Value)
			if err == nil {
				username = user.Username
				role = user.Role
			}

		}
		allPosts, err := h.service.GetPosts()
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}

		isRequestSent, err := h.service.CheckRequest(user.ID)
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}

		result := map[string]interface{}{
			"Posts":       allPosts,
			"CurrentUser": user,
			"Username":    username,
			"Role":        role,
			"RequestSent": isRequestSent,
		}
		tmpl, err := template.ParseFiles("ui/html/pages/home.html")
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		if err = tmpl.Execute(w, result); err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
	} else {
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
		return
	}
}

func (h *Handler) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	config := oauth.GetGoogleOAuth2Config()

	url := config.AuthCodeURL(oauth.GetGoogleOAuth2State(), oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	log.Printf("Получен код авторизации: %s", code)

	config := oauth.GetGoogleOAuth2Config()

	// Обмен кода на токен
	token, err := config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to exchange the token: %s", err), http.StatusInternalServerError)
		return
	}

	// Используем токен для получения информации о пользователе
	client := config.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to get user info: %s", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, fmt.Sprintf("Unable to parse user info: %s", err), http.StatusInternalServerError)
		return
	}

	// Проверяем, существует ли пользователь с этим email в базе
	user, err := h.service.Auth.GetUserByEmail(userInfo.Email)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, fmt.Sprintf("Error fetching user from database: %s", err), http.StatusInternalServerError)
		return
	}

	if user.ID == 0 {
		// Если пользователь не существует, создаем нового
		newUser := models.User{
			Username: userInfo.Name,
			Email:    userInfo.Email,
			GoogleID: &userInfo.ID,
		}
		fmt.Print(newUser)
		err := h.service.Auth.CreateUserGoogle(newUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to create user: %s", err), http.StatusInternalServerError)
			return
		}
		user = newUser
	}
	// } else {
	// 	// Если пользователь существует, обновляем его данные, если необходимо
	// 	if *user.GoogleID == "" {
	// 		*user.GoogleID = userInfo.ID
	// 		err := h.service.Auth.UpdateUserWithGoogleData(userInfo.ID)
	// 		if err != nil {
	// 			http.Error(w, fmt.Sprintf("Unable to update user with Google data: %s", err), http.StatusInternalServerError)
	// 			return
	// 		}
	// 	}
	// }

	// Создаем сессию для пользователя
	sessionToken, err := h.service.Auth.SetSession(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to create session: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Print(sessionToken)
	// Сохраняем сессионный токен в cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Expires:  time.Now().Add(3 * time.Hour), // Время жизни сессии
		HttpOnly: true,
		Path:     "/",
	})
	// a
	// Вход прошел успешно, перенаправляем пользователя на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) githubLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth.GitHubAuthURL()
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) GitHubLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Генерация URL для авторизации GitHub
	url := oauth.GitHubAuthURL()
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Получение токена через колбэк
	token, err := oauth.GitHubCallback(r)
	if err != nil {
		log.Printf("GitHub Callback Error: %v", err)
		http.Error(w, "Failed to authenticate with GitHub", http.StatusInternalServerError)
		return
	}

	// Получение данных пользователя из GitHub API
	userInfo, err := oauth.GetGitHubUserInfo(token.AccessToken)
	if err != nil {

		log.Printf("GitHub User Info Error: %v", err)
		http.Error(w, "Failed to get user info from GitHub", http.StatusInternalServerError)
		return
	}

	fmt.Print(userInfo.Email)

	// Проверка существующего пользователя или создание нового
	user, err := h.service.Auth.GetUserByEmailGithub(userInfo.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Database Error: %v", err)
		http.Error(w, "Failed to fetch user from database", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Eto nash ID %d", user.ID)

	if user.ID == 0 {
		// Создание нового пользователя
		newUser := models.User{
			Username: userInfo.Username,
			Email:    userInfo.Email,
			GitHubID: userInfo.GitHubID,
		}
		if err := h.service.Auth.CreateUserGitHub(newUser); err != nil {
			log.Printf("User Creation Error: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		user = newUser

	} else if *user.GitHubID == 0 {
		// Update the existing user with GitHub data
		user.GitHubID = userInfo.GitHubID
		user.Username = userInfo.Username
		user.Email = userInfo.Email

		fmt.Print(token.AccessToken)

		if err := h.service.Auth.UpdateUserWithGitHubData(token.AccessToken); err != nil {
			log.Printf("User Update Error: %v", err)
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}

	// Создание сессии
	sessionToken, err := h.service.Auth.SetSession(&user)
	if err != nil {
		log.Printf("Session Error: %v", err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Установка cookie с токеном
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		HttpOnly: true,
		Expires:  time.Now().Add(3 * time.Hour),
		Path:     "/",
	})

	// Перенаправление на главную страницу
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	nameFunction := "LoginPage"
	if r.URL.Path != "/login" {
		ErrorHandler(w, http.StatusNotFound, nameFunction)
		return
	}
	tmpl, err := template.ParseFiles("ui/html/pages/login.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, nameFunction)
		return
	}
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, nil); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
	} else if r.Method == http.MethodPost {
		user := models.User{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		if err := h.service.Auth.CheckPassword(user); err != nil {
			if err == models.ErrInvalidPassword || err == models.ErrUserNotFound {
				ErrorHandlerWithTemplate(tmpl, w, err, http.StatusUnauthorized)
				return
			} else {
				ErrorHandler(w, http.StatusInternalServerError, nameFunction)
				return
			}
		}

		realUser, err := h.service.Auth.GetUserByEmail(user.Email)
		if err != nil {

			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}

		token, err := h.service.Auth.SetSession(&realUser)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    token,
			Expires:  time.Now().Add(3 * time.Hour),
			HttpOnly: true,
		})
		w.Header().Set("Content-Type", "application/json")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
	}
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	nameFunction := "Register"
	if r.URL.Path != "/register" {
		ErrorHandler(w, http.StatusNotFound, nameFunction)
		return
	}
	tmpl, err := template.ParseFiles("ui/html/pages/signup.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, nameFunction)
		return
	}
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, nil); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
	case "POST":
		user := &models.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Role:     models.UserRole,
		}
		checkUser, err := h.service.GetUserByEmail(user.Email)
		if checkUser.Email == user.Email {
			log.Println(err)
			ErrorHandlerWithTemplate(tmpl, w, errors.New("Email already used"), http.StatusBadRequest)
			return
		}

		if len(user.Email) < 6 || len(user.Email) > 254 {
			ErrorHandlerWithTemplate(tmpl, w, errors.New("Email must be between 6 and 254 characters"), http.StatusBadRequest)
			return
		}

		if err := h.service.CheckUser(user); err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}

		if err := h.service.CreateUser(*user); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	nameFunction := "Logout"
	if r.URL.Path != "/logout" {
		ErrorHandler(w, http.StatusNotFound, nameFunction)
		return
	}
	switch r.Method {
	case "GET":
		// Получаем куку сессии
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			// Если кука не найдена, ошибку не генерируем (пользователь может быть не в системе)
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}

		// Удаляем сессию на сервере
		if err := h.service.Auth.DeleteSession(sessionCookie.Value); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}

		// Удаляем cookie сессии
		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,  // Это удаляет куку
			Path:   "/", // Обязательно указываем путь, чтобы она была удалена на всех страницах
		})

		// Перенаправляем пользователя на главную страницу
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
		return
	}
}
