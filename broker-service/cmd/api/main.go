package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = "80"

type AppConfig struct{}

func main() {

	app := AppConfig{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	log.Printf("Starting broker service on %s\n", PORT)
	srv.ListenAndServe()
}
