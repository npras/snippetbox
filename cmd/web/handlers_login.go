package main

import (
	"fmt"
	"net/http"

	"github.com/npras/snippetbox/internal/validator"
)

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            int    `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("GET: /user/signup")
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("POST /user/signup")
	var signupForm userSignupForm

	err := app.decodePostForm(r, &signupForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	validateSignupFields(&signupForm)
	if !signupForm.IsValid() {
		data := app.newTemplateData(r)
		data.Form = signupForm
		app.render(w, r, http.StatusOK, "signup.tmpl.html", data)
		return
	}
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "display login form")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new user session")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "delete the user session")
}

func validateSignupFields(f *signupForm) {
	f.CheckAndAddFieldError(validator.NotBlank(f.Name), "name", "name can't be blank")
	f.CheckAndAddFieldError(validator.LessThanMaxChars(f.Email, 100), "email", "email isn't valid")
	f.CheckAndAddFieldError(validator.AtLeastMaxChars(f.Password, 8), "password", "password must be at least 8 chars")
	f.CheckAndAddFieldError(validator.PermittedValue(f.ExpiresAt, 1, 7, 365), "expiresAt", "expiresAt should be one of 1, 7, 365")
}
