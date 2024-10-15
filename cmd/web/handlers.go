package main

import (
	"errors"
	"fmt"
	"github.com/VsProger/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

type postCreateForm struct {
	Title       string
	Text        string
	Categories  []int
	FieldErrors map[string]string
}

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

	post, categories, err := app.posts.Get(id)
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
	data.Categories = categories

	app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *application) showPostCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = postCreateForm{}

	app.render(w, r, http.StatusOK, "create.html", data)
}

func (app *application) doPostCreate(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	app.clientError(w, http.StatusBadRequest)
	//	return
	//}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	categories := r.Form["categories[]"]
	categoryIDs := make([]int, 0, len(categories))

	for _, cid := range categories {
		id, err := strconv.Atoi(cid)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		categoryIDs = append(categoryIDs, id)
	}
	form := postCreateForm{
		Title:       r.PostForm.Get("title"),
		Text:        r.PostForm.Get("content"),
		Categories:  categoryIDs,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Text) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	if len(form.Categories) == 0 {
		form.FieldErrors["categories"] = "This field must contain at least 1 category"
	}
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.posts.Insert(1, form.Title, form.Text, form.Categories)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}
