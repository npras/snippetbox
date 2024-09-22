package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error(),
		slog.String("method", r.Method),
		slog.String("uri", r.URL.RequestURI()),
		slog.String("trace", string(debug.Stack())))
	http.Error(w, "Internal Server Error. Oopsies", http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s doesn't exist", page)
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
