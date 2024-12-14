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
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	validateSignupFields(&form)
	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	id, err := app.user.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", err.Error())
			data := app.newTemplateData(r)
			data.Form = form
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

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("GET: /user/login")
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("POST /user/login")
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	validateLoginFields(&form)
	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err := app.user.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect!!!!")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.sessionManager.Put(r.Context(), "loggedInUserEmail", form.Email)
	app.sessionManager.Put(r.Context(), "flash", "Login Succeeded LOLL! Now create a snippet@!!")
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("POST /user/logout")
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Remove(r.Context(), "loggedInUserEmail")
	app.sessionManager.Put(r.Context(), "flash", "Logout successful! GET LOST!!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func validateSignupFields(f *userSignupForm) {
	f.CheckAndAddFieldError(validator.NotBlank(f.Name), "name", "name can't be blank")
	f.CheckAndAddFieldError(validator.NotBlank(f.Email), "email", "email can't be blank")
	f.CheckAndAddFieldError(validator.Matches(f.Email, validator.EmailRX), "email", "email isn't valid")
	f.CheckAndAddFieldError(validator.NotBlank(f.Password), "password", "password can't be blank")
	f.CheckAndAddFieldError(validator.AtLeastMaxChars(f.Password, 8), "password", "password must be at least 8 chars")
}

func validateLoginFields(f *userLoginForm) {
	f.CheckAndAddFieldError(validator.NotBlank(f.Email), "email", "email can't be blank")
	f.CheckAndAddFieldError(validator.Matches(f.Email, validator.EmailRX), "email", "email isn't valid")
	f.CheckAndAddFieldError(validator.NotBlank(f.Password), "password", "password can't be blank")
}
