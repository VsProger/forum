package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

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
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		allRepots, err := h.service.GetReports()
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}

		requests, err := h.service.GetRequests()

		result := map[string]interface{}{
			"Users":       allUsers,
			"CurrentUser": user,
			"Username":    username,
			"Role":        role,
			"Reports":     allRepots,
			"Requests":    requests,
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

func (h *Handler) upgradeOrDowngradeUser(w http.ResponseWriter, r *http.Request) {
	nameFunction := "upgradeOrDowngradeUser"
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid route", http.StatusBadRequest)
		return
	}
	action := parts[2]

	if r.Method == http.MethodPost {
		// Parse form values to get the user ID
		err := r.ParseForm()
		if err != nil {
			log.Println("Failed to parse form:", err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}

		userIDStr := r.FormValue("id")
		if userIDStr == "" {
			log.Println("User ID not provided")
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}

		var userID int
		_, err = fmt.Sscanf(userIDStr, "%d", &userID)
		if err != nil {
			log.Println("Invalid user ID format:", err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		switch action {
		case "upgrade":
			err = h.service.UpgradeUser(userID)
		case "downgrade":
			err = h.service.Downgrade(userID)
		default:
			ErrorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}
		// Call the service method to upgrade the user
		if err != nil {
			log.Println("Failed to upgrade user:", err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}

		http.Redirect(w, r, "/adminpage", http.StatusSeeOther)
	} else {
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
	}
}

func (h *Handler) requestRole(w http.ResponseWriter, r *http.Request) {
	nameFunction := "requestRole"
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println("Error parsing form:", err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		userIDStr := r.FormValue("id")
		if userIDStr == "" {
			log.Println("User ID not provided")
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Println("Invalid user ID format:", err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		if ok, err := h.service.CheckRequest(userID); ok {
			log.Println("User already requested role", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		err = h.service.RequestRole(userID)
		if err != nil {
			log.Println("Error requesting role:", err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
	}
}

func (h *Handler) approveUser(w http.ResponseWriter, r *http.Request) {
	nameFunction := "approveUser"
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println("Error parsing form:", err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		userIDStr := r.FormValue("id")
		if userIDStr == "" {
			log.Println("User ID not provided")
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Println("Invalid user ID format:", err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		err = h.service.ApproveRequest(userID)
		if err != nil {
			log.Println("Error approving request:", err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		http.Redirect(w, r, "/adminpage", http.StatusSeeOther)
	} else {
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
	}
}

func (h *Handler) declineUser(w http.ResponseWriter, r *http.Request) {
	nameFunction := "rejectUser"
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println("Error parsing form:", err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		userIDStr := r.FormValue("id")
		if userIDStr == "" {
			log.Println("User ID not provided")
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Println("Invalid user ID format:", err)
			ErrorHandler(w, http.StatusBadRequest, nameFunction)
			return
		}
		err = h.service.RejectRequest(userID)
		if err != nil {
			log.Println("Error rejecting request:", err)
			ErrorHandler(w, http.StatusInternalServerError, nameFunction)
			return
		}
		http.Redirect(w, r, "/adminpage", http.StatusSeeOther)
	} else {
		ErrorHandler(w, http.StatusMethodNotAllowed, nameFunction)
	}
}
