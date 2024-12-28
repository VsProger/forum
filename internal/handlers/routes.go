package handlers

import "net/http"

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/ui/static/", http.StripPrefix("/ui/static/", http.FileServer(http.Dir("./ui/static"))))

	mux.Handle("/myposts", h.AuthMiddleware(http.HandlerFunc(h.userPosts)))
	mux.Handle("/filter", http.HandlerFunc(h.filterByCategory))
	mux.Handle("/mylikedposts", h.AuthMiddleware(http.HandlerFunc(h.likePostsByUser)))
	mux.HandleFunc("/posts/reactions", h.addReaction)
	mux.HandleFunc("/posts/", h.getPost)
	mux.HandleFunc("/posts/create", h.createPost)

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

// func (h *Handler) routes() http.Handler {
// 	router := httprouter.New()

// 	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		h.notFound(w)
// 	})

// 	fileServer := http.FileServer(http.Dir("./ui/static/"))
// 	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

// 	router.HandlerFunc(http.MethodGet, "/", app.home)
// 	router.HandlerFunc(http.MethodGet, "/post/view/:id", app.postView)
// 	router.HandlerFunc(http.MethodGet, "/post/create", app.showPostCreate)
// 	router.HandlerFunc(http.MethodPost, "/post/create", app.doPostCreate)

// 	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

// 	return standard.Then(router)
// }
