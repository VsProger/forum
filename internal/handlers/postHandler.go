package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/pkg"
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
		ErrorHandler(w, http.StatusInternalServerError, nameFunction)
		return
	}
	result := map[string]interface{}{
		"Posts":    allPosts,
		"Username": username,
	}
	tmpl, err := template.ParseFiles("ui/html/pages/index.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, nameFunction)
		return
	}
	if err = tmpl.Execute(w, result); err != nil {
		ErrorHandler(w, http.StatusInternalServerError, nameFunction)
		return
	}

	ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
	return
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	nameFunction := "CreatePost"
	tmpl, err := template.ParseFiles("ui/html/pages/createPost.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, nameFunction)
		return
	}
	if r.Method == http.MethodGet {
		if err := tmpl.Execute(w, nil); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
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
		categories := r.Form["categories"]

		post := models.Post{
			Title: r.FormValue("title"),
			Text:  r.FormValue("text"),
		}

		for _, name := range categories {
			post.Categories = append(post.Categories, models.Category{Name: name})
		}
		if err := pkg.VallidatePost(post); err != nil {
			ErrorHandlerWithTemplate(tmpl, w, err, http.StatusBadRequest)
			return
		}
		post.AuthorID = user.ID
		if err := h.service.PostService.CreatePost(post); err != nil {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	nameFunction := "getPost"
	tmpl, err := template.ParseFiles("ui/html/pages/post.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "getPost")
		return
	}
	if r.Method == http.MethodGet {
		idStr := r.URL.Path[len("/posts/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		post, err := h.service.GetPostByID(id)
		if err != nil || idStr == "" || len(idStr) > 2 || id > 50 || id <= 0 {
			ErrorHandler(w, http.StatusNotFound, nameFunction)
			return
		}
		var username string
		session, err := r.Cookie("session")
		if err == nil {
			user, err := h.service.GetUserByToken(session.Value)
			if err == nil {
				username = user.Username
			}
		}
		result := map[string]interface{}{
			"Post":          post,
			"Authenticated": username,
		}

		if err = tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "getPost")
			return
		}
	} else if r.Method == http.MethodPost {
		idStr := r.URL.Path[len("/posts/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		session, err := r.Cookie("session")
		if err != nil {
			ErrorHandler(w, http.StatusUnauthorized, nameFunction)
			return
		}
		user, err := h.service.Auth.GetUserByToken(session.Value)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		post, err := h.service.GetPostByID(id)
		if err != nil {
			if idStr == "" || len(idStr) > 2 || id > 50 || id <= 0 {
				ErrorHandler(w, http.StatusNotFound, nameFunction)
				return
			}
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		result := map[string]interface{}{
			"Post":          post,
			"Authenticated": user.Username,
		}
		comment := models.Comment{
			Text:     r.FormValue("text"),
			PostID:   id,
			AuthorID: user.ID,
		}

		if err := h.service.CreateComment(comment); err != nil {
			if err == models.ErrEmptyComment || err == models.ErrInvalidComment || err == models.ErrNotAscii {
				ErrorHandler(w, http.StatusBadRequest, nameFunction)
				return
			}
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		if err := tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
