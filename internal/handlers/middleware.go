package handlers

import (
	"context"
	"net/http"

	"github.com/VsProger/snippetbox/logger"
)

var logg = logger.NewLogger()

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")
		if err != nil || sessionCookie.Value == "" {

			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, err := h.service.GetUserByToken(sessionCookie.Value)
		if err != nil {

			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// func (h *Handler) logRequest(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var (
// 			ip     = r.RemoteAddr
// 			proto  = r.Proto
// 			method = r.Method
// 			uri    = r.URL.RequestURI()
// 		)

// 		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

// 		next.ServeHTTP(w, r)
// 	})
// }

func (h *Handler) AllHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logg.Info(r.Method + " successfully working")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RoleMiddleware(requiredRoles []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the session cookie
		session, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Retrieve the user associated with the session token
		user, err := h.service.GetUserByToken(session.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		hasAccess := false
		for _, role := range requiredRoles {
			if user.Role == role {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Proceed to the next handler if the role matches
		next.ServeHTTP(w, r)
	})
}
