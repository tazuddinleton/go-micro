package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var Port = flag.String("port", "80", "port of broker service")

type AppConfig struct{}

func main() {
	flag.Parse()
	app := AppConfig{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", *Port),
		Handler: app.routes(),
	}

	log.Printf("Starting broker service on %s\n", *Port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
