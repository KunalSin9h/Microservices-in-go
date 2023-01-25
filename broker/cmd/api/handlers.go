package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var request RequestPayload
	err := app.readJson(w, r, &request)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	switch request.Action {
	case "auth":
		app.authenticate(w, request.Auth)
	default:
		app.errorJson(w, errors.New("unknown action"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, load AuthPayload) {
	jsonData, err := json.MarshalIndent(load, "", "\t")
	if err != nil {
		app.errorJson(w, errors.New("unable to marshal data"), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", "http://auth/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("invalid credentials"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling auth service"))
	}

	var response jsonResponse

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	if response.Error {
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticate"
	payload.Data = response.Data

	app.writeJson(w, http.StatusAccepted, payload)
}
