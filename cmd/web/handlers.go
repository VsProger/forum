package main

import (
	"errors"
	"fmt"
	"github.com/VsProger/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//if r.URL.Path != "/" {
	//	app.notFound(w)
	//	return
	//}

	snippets, err := app.posts.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Posts = snippets

	app.render(w, r, http.StatusOK, "home.html", data)

}

func (app *application) postView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	post, err := app.posts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Post = post

	app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *application) showPostCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display form for creating a new snippet..."))
}

//func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
//	//if r.Method != http.MethodPost {
//	//	w.Header().Set("Allow", http.MethodPost)
//	//	app.clientError(w, http.StatusBadRequest)
//	//	return
//	//}
//
//	title := "0 snail"
//	content := "snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
//	expires := 7
//
//	id, err := app.posts.Insert(title, content, expires)
//	if err != nil {
//		app.serverError(w, r, err)
//		return
//	}
//
//	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
//}

func (app *application) doCreatePost(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	app.clientError(w, http.StatusBadRequest)
	//	return
	//}

	authorID := 1
	categoryID := 1

	title := "HEEllllllllll YEEEaaaaaa snail"
	content := "snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"

	id, err := app.posts.Insert(authorID, title, content, categoryID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}
