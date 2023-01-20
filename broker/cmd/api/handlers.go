package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

/*
Broker is http handler function
*/
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	log.Println("Hit Broker Handler")
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	outBytes, _ := json.MarshalIndent(payload, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(outBytes)
}
