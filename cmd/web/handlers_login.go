package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/npras/snippetbox/internal/models"
	"github.com/npras/snippetbox/internal/validator"
)

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
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
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	id, err := app.user.Insert(signupForm.Name, signupForm.Email, signupForm.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			signupForm.AddFieldError("email", err.Error())
			data := app.newTemplateData(r)
			data.Form = signupForm
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	msg := fmt.Sprintf("User with id %d created successfully! Login!!!", id)
	app.sessionManager.Put(r.Context(), "flash", msg)
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
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

func validateSignupFields(f *userSignupForm) {
	f.CheckAndAddFieldError(validator.NotBlank(f.Name), "name", "name can't be blank")
	f.CheckAndAddFieldError(validator.NotBlank(f.Email), "email", "email can't be blank")
	f.CheckAndAddFieldError(validator.Matches(f.Email, validator.EmailRX), "email", "email isn't valid")
	f.CheckAndAddFieldError(validator.NotBlank(f.Password), "password", "password can't be blank")
	f.CheckAndAddFieldError(validator.AtLeastMaxChars(f.Password, 8), "password", "password must be at least 8 chars")
}
