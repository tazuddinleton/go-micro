package main

import (
	"auth/cmd/data"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"time"
)

var Port = flag.String("port", "80", "auth service port")

const DbConnRetry = 10

type AppConfig struct {
	db     *sql.DB
	Models data.Models
}

func main() {
	flag.Parse()
	log.Println("Starting authentication service")

	db := connectToDB()
	if db == nil {
		log.Panic("Could not connect to db")
	}

	app := &AppConfig{
		db:     db,
		Models: data.New(db),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", *Port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	retry := 0
	for retry < DbConnRetry {
		conn, err := openDB(dsn)
		if err != nil {
			retry++
			log.Println("DB is not ready yet ...")
			time.Sleep(2 * time.Second)
			continue
		}
		log.Println("Connected to DB!")
		return conn
	}
	log.Println("Could not connect to db, retry max out.")
	return nil
}
