package handlers

import (
	"forum/internal/models"
	"html/template"
	"net/http"
)

func (h *Handler) filterByCategory(w http.ResponseWriter, r *http.Request) {
	nameFunction := "filterByCategory"
	if r.URL.Path != "/filter" {
		ErrorHandler(w, http.StatusNotFound, nameFunction)
		return
	}
	switch r.Method {
	case "GET":
		tmpl, err := template.ParseFiles("ui/html/pages/home.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		var user models.User
		username := ""
		session, err := r.Cookie("session")
		if err == nil {
			user, err = h.service.GetUserByToken(session.Value)
			if err == nil {
				username = user.Username
			}
		}
		if err := r.ParseForm(); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		categories := r.Form["categories"]
		if len(categories) == 0 || len(categories) == 4 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		category, err := h.service.GetCategoryByName(categories)
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}

		posts, err := h.service.FilterByCategories(category)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		result := map[string]interface{}{
			"Posts":    posts,
			"Username": username,
		}

		if err := tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
		return
	}
}
