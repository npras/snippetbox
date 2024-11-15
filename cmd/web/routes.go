package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(app.config.staticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamicMiddlewares := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamicMiddlewares.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamicMiddlewares.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamicMiddlewares.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamicMiddlewares.ThenFunc(app.snippetCreatePost))

	standardMiddlewares := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standardMiddlewares.Then(mux)
}
