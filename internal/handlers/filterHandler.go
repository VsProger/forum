package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/VsProger/snippetbox/internal/models"
)

func (h *Handler) likePostsByUser(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("ui/html/pages/home.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "likePosts")
		return
	}
	if r.Method == http.MethodGet {
		nameFunction := "likePosts"
		session, err := r.Cookie("session")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		user, err := h.service.GetUserByToken(session.Value)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		posts, err := h.service.FilterByLikes(user.ID)
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		result := map[string]interface{}{
			"Posts":    posts,
			"Username": user.Username,
		}
		if err = tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "likePosts")
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) dislikePostsByUser(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("ui/html/pages/home.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "likePosts")
		return
	}
	if r.Method == http.MethodGet {
		nameFunction := "likePosts"
		session, err := r.Cookie("session")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		user, err := h.service.GetUserByToken(session.Value)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		posts, err := h.service.FilterByDislikes(user.ID)
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		result := map[string]interface{}{
			"Posts":    posts,
			"Username": user.Username,
		}
		if err = tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "likePosts")
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

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
			log.Print(err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		categories := r.Form["categories"]
		if len(categories) == 0 || len(categories) == 4 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		category, err := h.service.GetCategoryByName(categories)
		if err != nil {
			log.Fatal(err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}

		posts, err := h.service.FilterByCategories(category)
		if err != nil {
			log.Print(err)
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
