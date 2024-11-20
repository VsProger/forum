package handlers

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/VsProger/snippetbox/internal/models"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	nameFunction := "indexHandler"
	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusNotFound, nameFunction)
		return
	}

	var username string

	allPosts, err := h.service.GetPosts()
	if err != nil {
		log.Fatal(err)
		ErrorHandler(w, http.StatusInternalServerError, nameFunction)
		return
	}
	result := map[string]interface{}{
		"Posts":    allPosts,
		"Username": username,
	}
	tmpl, err := template.ParseFiles("/home/student/Desktop/forum/ui/html/pages/home.html")
	if err != nil {
		log.Fatal(err)
		ErrorHandler(w, http.StatusInternalServerError, nameFunction)
		return
	}
	if err = tmpl.Execute(w, result); err != nil {
		log.Fatal(err)
		ErrorHandler(w, http.StatusInternalServerError, nameFunction)
		return
	}

	ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
	return
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
		}
		checkUser, _ := h.service.GetUserByEmail(user.Email)
		if checkUser.Email == user.Email {
			ErrorHandlerWithTemplate(tmpl, w, errors.New("Email already used"), http.StatusBadRequest)
			return
		}
		checkUser, _ = h.service.GetUserByUSername(user.Username)
		if checkUser.Username == user.Username {
			ErrorHandlerWithTemplate(tmpl, w, errors.New("Username already used"), http.StatusBadRequest)
			return
		}
		if err := h.service.CheckUser(user); err != nil {
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
