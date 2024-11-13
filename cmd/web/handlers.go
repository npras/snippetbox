package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/npras/snippetbox/internal/models"
	"github.com/npras/snippetbox/internal/validator"
)

//

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//panic("OOPS. Things blew up!")
	snippets, err := app.snippet.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

//

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	snippet, err := app.snippet.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

//

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("RENDERING CREATE SNIPPET")
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		ExpiresAt: 365,
	}
	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

//

type snippetCreateForm struct {
	Title     string
	Content   string
	ExpiresAt int
	validator.Validator
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("CREATING SNIPPET")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expiresAt, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:     r.PostForm.Get("title"),
		Content:   r.PostForm.Get("content"),
		ExpiresAt: expiresAt,
	}

	validateSnippetCreateFields(&form)
	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "create.tmpl.html", data)
		return
	}

	id, err := app.snippet.Insert(form.Title, form.Content, form.ExpiresAt)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	redirectToURL := fmt.Sprintf("/snippet/view/%d", id)
	http.Redirect(w, r, redirectToURL, http.StatusSeeOther)
}

func validateSnippetCreateFields(f *snippetCreateForm) {
	f.CheckAndAddFieldError(validator.NotBlank(f.Title), "title", "title can't be blank")
	f.CheckAndAddFieldError(validator.LessThanMaxChars(f.Title, 100), "title", "title should be less than 100 chars")
	f.CheckAndAddFieldError(validator.NotBlank(f.Content), "content", "content can't be blank")
	f.CheckAndAddFieldError(validator.PermittedValue(f.ExpiresAt, 1, 7, 365), "expiresAt", "expiresAt should be one of 1, 7, 365")
}
