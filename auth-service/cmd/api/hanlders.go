package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const loggerServiceUrl = "http://logger-service/log"

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
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Error:   false,
		Data:    user,
	}
	log.Println("calling logger-service")
	err = app.logAuthentication("authentication", fmt.Sprintf("Logged in %s", user.Email))
	if err != nil {
		log.Println("error calling logger-service")
	}

	err = app.writeJSON(w, http.StatusOK, resPayload)
	if err != nil {
		return
	}
}

func (app *AppConfig) logAuthentication(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, err := json.MarshalIndent(entry, "", "/t")
	if err != nil {
		return err
	}

	client := http.Client{}
	log.Println(string(jsonData))
	req, err := http.NewRequest(http.MethodPost, loggerServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
