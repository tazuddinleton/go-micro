package main

import (
	"logger-service/data"
	"net/http"
)

func (app *AppConfig) Log(w http.ResponseWriter, r *http.Request) {
	var reqPayload struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	err := app.readJSON(w, r, &reqPayload)

	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = app.Models.LogEntry.Insert(data.LogEntry{
		Title: reqPayload.Name,
		Data:  reqPayload.Data,
	})
	if err != nil {
		_ = app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resp := jsonResponse{Message: "Logged.", Error: false}
	err = app.writeJSON(w, http.StatusOK, resp)
	if err != nil {
		return
	}
}
