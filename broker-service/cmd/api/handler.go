package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *AppConfig) Broker(w http.ResponseWriter, r *http.Request) {

	res := jsonResponse{
		Error:   false,
		Message: "Hit the broker!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	out, _ := json.MarshalIndent(res, "", "\t")

	w.Write(out)
}
