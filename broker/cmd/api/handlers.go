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
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
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
	case "log":
		app.LogItem(w, request.Log)
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

	req, err := http.NewRequest("POST", "http://auth:5002/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

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
	payload.Message = response.Message
	payload.Data = response.Data

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) LogItem(w http.ResponseWriter, load LogPayload) {

	jsonByte, _ := json.MarshalIndent(load, "", "\t")
	res, err := http.Post("http://logger:5003/log", "application/json", bytes.NewBuffer(jsonByte))

	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusBadRequest:
		app.errorJson(w, errors.New("bad request"), http.StatusBadRequest)
	case http.StatusInternalServerError:
		app.errorJson(w, errors.New("internal server error"), http.StatusInternalServerError)
	default:

		var resPayload jsonResponse
		resPayload.Error = false
		resPayload.Message = "logged"
		app.writeJson(w, http.StatusOK, resPayload)
	}
}
