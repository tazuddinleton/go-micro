package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"logger-service/data"
)

var webPort = flag.String("port", "80", "Logger service port")

// var rpcPort = flag.String("rpcport", "5001", "RPC port")
var mongoUrl = flag.String("mongourl", "mongodb://mongo:27017", "Mongo db connection string")

// var gRpcPort = flag.String("grpcport", "50001", "gRpc port")
type AppConfig struct {
	Models *data.Models
}

var mongoClient *mongo.Client

func main() {
	flag.Parse()

	conn, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	mongoClient = conn

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	app := &AppConfig{Models: data.NewModels(mongoClient)}

	srv := http.Server{
		Handler: app.routes(),
		Addr:    fmt.Sprintf(":%s", *webPort),
	}
	log.Println(fmt.Sprintf("logger service started at http://localhost:%s", *webPort))
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	opt := options.Client().ApplyURI(*mongoUrl)
	username := os.Getenv("MongoUsername")
	password := os.Getenv("MongoPassword")
	log.Println(fmt.Sprintf("credentials: username: %s, pass: %s", username, password))
	opt.SetAuth(
		options.Credential{Username: username, Password: password, AuthMechanism: "SCRAM-SHA-256"},
	)

	conn, err := mongo.Connect(context.TODO(), opt)
	if err != nil {
		log.Println("error connecting to mongodb", err)
		return nil, err
	}

	err = conn.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("error connecting to mongodb", err)
		return nil, err
	}

	log.Println("connected to mongodb!")
	return conn, nil
}
