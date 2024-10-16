package main

import (
	"errors"
	"fmt"
	"github.com/VsProger/snippetbox/internal/models"
	"github.com/VsProger/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type postCreateForm struct {
	Title       string
	Text        string
	CategoryIDs []int
	validator.Validator
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
	if len(categoryIDs) == 0 {
		categoryIDs = append(categoryIDs, 5)
	}
	form := postCreateForm{
		Title:       r.PostForm.Get("title"),
		Text:        r.PostForm.Get("content"),
		CategoryIDs: categoryIDs,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Text), "content", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.posts.Insert(1, form.Title, form.Text, form.CategoryIDs)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}
