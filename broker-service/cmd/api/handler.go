package main

import (
	"broker/cmd/model"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const (
	loggerServiceUrl = "http://logger-service/log"
	authServiceUrl   = "http://auth-service/authenticate"
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
		var requestPayload model.RequestPayload

		err := json.NewDecoder(r.Body).Decode(&requestPayload)
		if err != nil {
			_ = app.errorJSON(w, errors.New("error occurred"), http.StatusBadRequest)
		}

		switch requestPayload.Action {
		case "authenticate":
			log.Println("handling authentication request...")
			app.HandleAuthentication(w, requestPayload.Auth)
			break
		case "log":
			log.Println("handling logging request...")
			app.HandleLogging(w, requestPayload.Log)
			break
		default:
			_ = app.errorJSON(w, errors.New("unknown action"), http.StatusBadRequest)
		}
	}
}

func (app *AppConfig) HandleAuthentication(w http.ResponseWriter, auth model.AuthPayload) {
	res, err := app.doRequest(authServiceUrl, http.MethodPost, auth)
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

	_ = app.writeJSON(w, http.StatusOK, success(jsonAuth.Data, "Authenticated"))
}

func (app *AppConfig) HandleLogging(w http.ResponseWriter, l model.LogPayload) {
	res, err := app.doRequest(loggerServiceUrl, http.MethodPost, l)
	if err != nil {
		log.Println("failed to call logger-service", err)
		_ = app.errorJSON(w, errors.New("failed to call logger-service"), http.StatusInternalServerError)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Println("failed to call logger-service", err)
		_ = app.errorJSON(w, errors.New("request not valid"), http.StatusBadRequest)
	}

	_ = app.writeJSON(w, http.StatusOK, success(nil, "logged."))
}

func success(data any, message string) jsonResponse {
	var payload jsonResponse
	payload.Error = false
	payload.Message = message
	payload.Data = data
	return payload
}

func (app *AppConfig) doRequest(url, method string, payload any) (*http.Response, error) {
	jsonData, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
