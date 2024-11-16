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

	mux.Handle("GET /user/signup", dynamicMiddlewares.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamicMiddlewares.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamicMiddlewares.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamicMiddlewares.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", dynamicMiddlewares.ThenFunc(app.userLogoutPost))

	standardMiddlewares := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standardMiddlewares.Then(mux)
}
