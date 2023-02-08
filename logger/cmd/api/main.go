package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"logger/data"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	PORT          = "5003"
	RPC_PORT      = "5031"
	GRPC_PORT     = "5032"
	MONGO_DB_CONN = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	conn, err := connectToMongo()

	if err != nil {
		log.Printf("[MONGO CONN] Error connecting to mongodb: %v", err)
		os.Exit(1)
	}

	client = conn

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("[MONGO DIS] Error disconnecting to mongodb: %v", err)
			os.Exit(1)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	app.serve()
}

func (app *Config) serve() {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", PORT),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("Starting server at port %s\n", PORT)
	log.Fatal(server.ListenAndServe())
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(MONGO_DB_CONN)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	})

	c, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		return nil, err
	}

	return c, nil
}