package main

import "net/http"

/*
Broker is http handler function
*/
func (app *Config) broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.writeJson(w, http.StatusOK, payload)
}
