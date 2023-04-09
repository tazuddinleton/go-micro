package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *AppConfig) Broker(w http.ResponseWriter, r *http.Request) {

	res := jsonResponse{
		Error:   false,
		Message: "Hit the broker!",
	}

	_ = app.writeJSON(w, http.StatusAccepted, res)
}

func (app *AppConfig) HandleSubmission() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestPayload RequestPayload

		err := json.NewDecoder(r.Body).Decode(&requestPayload)
		if err != nil {
			_ = app.errorJSON(w, errors.New("error occurred"), http.StatusBadRequest)
		}

		switch requestPayload.Action {
		case "authenticate":
			log.Println("handling authentication request...")
			app.HandleAuthentication(w, requestPayload.Auth)
		default:
			_ = app.errorJSON(w, errors.New("unknown action"), http.StatusBadRequest)
		}
	}
}

func (app *AppConfig) HandleAuthentication(w http.ResponseWriter, auth AuthPayload) {

	jsonData, err := json.MarshalIndent(auth, "", "\t")
	if err != nil {
		_ = app.errorJSON(w, errors.New("error calling auth service"), http.StatusBadRequest)
		return
	}

	request, err := http.NewRequest(http.MethodPost, "http://auth-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	client := &http.Client{}
	//res, err := client.Post("http://localhost:8092/authenticate", "Content-Type: application/json", bytes.NewBuffer(jsonData))

	res, err := client.Do(request)
	if err != nil {
		log.Println("error calling auth service", err)
		_ = app.errorJSON(w, errors.New("error calling auth service"), http.StatusBadRequest)
		return
	}
	defer res.Body.Close()

	log.Println("response from auth service: status ", res.StatusCode)
	if res.StatusCode == http.StatusUnauthorized {
		log.Println("unauthorized", err)
		_ = app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Println("error calling auth service", err)
		_ = app.errorJSON(w, errors.New("error calling auth service"), http.StatusUnauthorized)
		return
	}

	var jsonAuth jsonResponse
	err = json.NewDecoder(res.Body).Decode(&jsonAuth)

	if err != nil {
		log.Println("error calling auth service", err)
		_ = app.errorJSON(w, errors.New("error calling auth service"), http.StatusBadRequest)
		return
	}

	// Unnecessary
	if jsonAuth.Error {
		log.Println("unauthorized", err)
		_ = app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonAuth.Data

	log.Println(fmt.Sprintf("sending back payload: %v", payload))
	_ = app.writeJSON(w, http.StatusOK, payload)
}
