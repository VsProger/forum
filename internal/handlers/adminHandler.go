package handlers

import (
	"html/template"
	"net/http"

	"github.com/VsProger/snippetbox/internal/models"
)

func (h *Handler) adminpage(w http.ResponseWriter, r *http.Request) {
	nameFunction := "adminpageHandler"
	if r.URL.Path != "/adminpage" {
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
		allUsers, err := h.service.GetUsers()
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		result := map[string]interface{}{
			"Users":       allUsers,
			"CurrentUser": user,
			"Username":    username,
			"Role":        role,
		}
		tmpl, err := template.ParseFiles("ui/html/pages/admin.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		if err = tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
	} else {
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
		return
	}
}
