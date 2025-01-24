package main

import (
	"net/http"
	"strings"

	"github.com/maksimfisenko/goform-server-app/internal/data"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RoleTitle string `json:"role_title"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.RoleTitle == "" {
		input.RoleTitle = "RESPONDER"
	}

	role, err := app.storage.Roles.GetByTitle(strings.ToUpper(input.RoleTitle))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := &data.User{
		RoleID:      role.ID,
		Name:        input.Name,
		Email:       input.Email,
		IsActivated: true,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.storage.Users.Insert(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
