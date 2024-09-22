package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/npras/snippetbox/internal/models"
)

//

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	snippets, err := app.snippet.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := templateData{Snippets: snippets}
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
	data := templateData{Snippet: snippet}
	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

//

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

//

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "99 snail"
	content := "99 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expiresAt := 10

	id, err := app.snippet.Insert(title, content, expiresAt)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	redirectToURL := fmt.Sprintf("/snippet/view/%d", id)
	http.Redirect(w, r, redirectToURL, http.StatusSeeOther)
}
