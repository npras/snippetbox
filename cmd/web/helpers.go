package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error(),
		slog.String("method", r.Method),
		slog.String("uri", r.URL.RequestURI()),
		slog.String("trace", string(debug.Stack())))
	//http.Error(w, "Internal Server Error. Oopsies", http.StatusInternalServerError)
	http.Error(w, fmt.Sprintf("Internal Server Error. Oopsies: %T", err), http.StatusInternalServerError)
}

//

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

//

func (app *application) newTemplateData(r *http.Request) templateData {
	loggedInUserEmail := ""
	if app.isAuthenticated(r) {
		loggedInUserEmail = app.sessionManager.GetString(r.Context(), "loggedInUserEmail")
	}
	return templateData{
		CurrentYear:       time.Now().Year(),
		Flash:             app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated:   app.isAuthenticated(r),
		CSRFToken:         nosurf.Token(r),
		LoggedInUserEmail: loggedInUserEmail,
	}
}

//

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s doesn't exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	// return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
