package main

import (
	"net/http"
)

func (app *AppConfig) Broker(w http.ResponseWriter, r *http.Request) {

	res := jsonResponse{
		Error:   false,
		Message: "Hit the broker!",
	}

	_ = app.writeJSON(w, http.StatusAccepted, res)
}
