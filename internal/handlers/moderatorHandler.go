package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) reportPost(w http.ResponseWriter, r *http.Request) {
	nameFunction := "reportPostHandler"
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println("Error parsing form:", err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		postIDstr := r.FormValue("postId")
		reason := "Breaks forum rules"
		fmt.Println("post", postIDstr)
		if postIDstr == "" {
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		var postID int
		postID, err = strconv.Atoi(postIDstr)

		if err != nil {
			log.Println("Invalid post ID format:", err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		session, err := r.Cookie("session")
		if err != nil {
			log.Println("Error getting session cookie:", err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		user, err := h.service.GetUserByToken(session.Value)
		if err != nil {
			log.Println("Error getting user by token:", err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		err = h.service.ReportPost(postID, user.ID, reason)
		if err != nil {
			log.Println("Error reporting post:", err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		http.Redirect(w, r, "/posts/"+postIDstr, http.StatusSeeOther)
	} else {
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
		return
	}
}
