package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/pkg"
)

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
			log.Fatal(err)
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

func (h *Handler) userPosts(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("/home/student/forum/ui/html/pages/home.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "getPost")
		return
	}
	if r.Method == http.MethodGet {
		nameFunction := "userPosts"
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
		posts, err := h.service.GetPostsByUserId(user.ID)
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		result := map[string]interface{}{
			"Posts":    posts,
			"Username": user.Username,
		}
		if err = tmpl.Execute(w, result); err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "userPosts")
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) addReaction(w http.ResponseWriter, r *http.Request) {
	nameFunction := "addReaction"
	if r.Method == http.MethodPost {
		session, err := r.Cookie("session")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		user, err := h.service.Auth.GetUserByToken(session.Value)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		postId, err := pkg.Atoi(r.FormValue("postId"))
		if err != nil {
			ErrorHandler(w, http.StatusNotFound, nameFunction)
			return
		}
		var commentId int
		if r.FormValue("commentId") != "" {
			commentId, err = pkg.Atoi(r.FormValue("commentId"))
			if err != nil {
				ErrorHandler(w, http.StatusBadRequest, nameFunction)
				return
			}
		}
		vote, err := pkg.Atoi(r.FormValue("status"))
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		reaction := models.Reaction{
			UserID:    user.ID,
			PostID:    postId,
			CommentID: commentId,
			Vote:      vote,
		}
		if err := h.service.AddReaction(reaction); err != nil {
			if err == fmt.Errorf("specify either PostId or CommentId, not both") || strings.Contains(err.Error(), "Vote IN (-1, 1)") {
				ErrorHandler(w, http.StatusBadRequest, nameFunction)
				return
			} else if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				ErrorHandler(w, http.StatusNotFound, nameFunction)
				return
			}
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		path := "/posts/" + r.FormValue("postId")
		http.Redirect(w, r, path, http.StatusSeeOther)
	} else {
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
		return
	}
}
