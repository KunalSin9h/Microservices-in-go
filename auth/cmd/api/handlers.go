package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &requestPayload)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	// validate user against the database

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.errorJson(w, errors.New("invalid credentials"))
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		app.errorJson(w, errors.New("invalid credentials"))
		return
	}

	err = app.logAuth("[AUTH]", "Authenticated as "+user.FirstName)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in as %s", user.Email),
		Data:    user,
	}

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) logAuth(name, data string) error {

	var logEntry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	logEntry.Name = name
	logEntry.Data = data

	jsonByte, _ := json.MarshalIndent(logEntry, "", "\t")

	_, err := http.Post("http://logger:5003/log", "application/json", bytes.NewBuffer(jsonByte))

	if err != nil {
		return err
	}

	return nil

}
