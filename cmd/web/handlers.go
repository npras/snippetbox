package main

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/npras/snippetbox/internal/models"
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
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

//

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("CREATING SNIPPET")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expiresAt, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fieldErrors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "title can't be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "title should be less than 100 chars"
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "content can't be blank"
	}

	validExpiresAt := []int{1, 7, 365}
	if !slices.Contains(validExpiresAt, expiresAt) {
		fieldErrors["expiresAt"] = "expiresAt should be one of 1, 7, 365"
	}

	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

	id, err := app.snippet.Insert(title, content, expiresAt)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	redirectToURL := fmt.Sprintf("/snippet/view/%d", id)
	http.Redirect(w, r, redirectToURL, http.StatusSeeOther)
}
