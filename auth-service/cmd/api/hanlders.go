package main

import (
	"errors"
	"log"
	"net/http"
)

func (app *AppConfig) Authenticate(w http.ResponseWriter, r *http.Request) {
	var reqPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &reqPayload)

	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log.Println("request payload: ", reqPayload)
	user, err := app.Models.User.GetByEmail(reqPayload.Email)
	if err != nil || user == nil {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	log.Println("valid user found: ", user)
	valid, err := user.PasswordMatches(reqPayload.Password)

	if !valid {
		log.Println("password didn't  match!   ", reqPayload.Password)
	}
	if err != nil || !valid {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	resPayload := jsonResponse{
		Message: "Authenticated",
		Error:   false,
		Data:    user,
	}

	err = app.writeJSON(w, http.StatusOK, resPayload)
	if err != nil {
		return
	}
}
