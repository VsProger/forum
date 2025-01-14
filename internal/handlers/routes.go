package handlers

import "net/http"

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/ui/static/", http.StripPrefix("/ui/static/", http.FileServer(http.Dir("./ui/static"))))

	mux.Handle("/myposts", h.AuthMiddleware(http.HandlerFunc(h.userPosts)))
	mux.Handle("/filter", http.HandlerFunc(h.filterByCategory))
	mux.Handle("/mylikedposts", h.AuthMiddleware(http.HandlerFunc(h.likePostsByUser)))
	mux.Handle("/mydislikedposts", h.AuthMiddleware(http.HandlerFunc(h.dislikePostsByUser)))
	mux.Handle("/posts/create", h.AuthMiddleware(http.HandlerFunc(h.createPost)))
	mux.Handle("/posts/reactions", h.AuthMiddleware(http.HandlerFunc(h.addReaction)))
	mux.Handle("/postsdelete/", h.RoleMiddleware([]string{"admin", "moderator"}, http.HandlerFunc(h.DeletePost)))

	///todo
	// mux.Handle("/posts/report", h.RoleMiddleware([]string{"moderator"}, http.HandlerFunc(h.)))
	// mux.Handle("/user/upgrade", h.RoleMiddleware([]string{"admin"}, http.HandlerFunc(h.)))
	// mux.Handle("/user/downgrade", h.RoleMiddleware([]string{"admin"}, http.HandlerFunc(h.)))
	mux.Handle("/adminpage", h.RoleMiddleware([]string{"admin"}, http.HandlerFunc(h.adminpage)))

	mux.Handle("/postsedit/", h.AuthMiddleware(http.HandlerFunc(h.editPost)))

	mux.HandleFunc("/posts/", h.getPost)
	mux.HandleFunc("/userComments/", h.userComments)

	mux.HandleFunc("/auth/google", h.GoogleLoginHandler)
	mux.HandleFunc("/notifications", h.GetNotificationsHandler)

	mux.HandleFunc("/auth/google/callback", h.GoogleCallbackHandler)

	mux.HandleFunc("/auth/github", h.githubLogin)
	mux.HandleFunc("/auth/github/callback", h.GitHubCallbackHandler)

	mux.HandleFunc("/", h.home)
	mux.HandleFunc("/login", h.login)
	mux.HandleFunc("/register", h.register)
	mux.HandleFunc("/logout", h.logout)

	return h.AllHandler(mux)
}
