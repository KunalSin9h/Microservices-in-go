package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = "5001"

type Config struct{}

func main() {

	app := Config{}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.Routes(),
	}

	log.Printf("Started Broker Server at port: %s", PORT)
	log.Fatal(server.ListenAndServe())
}
