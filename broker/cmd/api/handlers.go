package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"

	"broker/event"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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
		// app.logItem(w, request.Log)
		// app.logItemAMQP(w, request.Log)
		app.LogItemRPC(w, request.Log)
	case "mail":
		app.sendMail(w, request.Mail)
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

func (app *Config) logItem(w http.ResponseWriter, load LogPayload) {

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

func (app *Config) sendMail(w http.ResponseWriter, load MailPayload) {
	jsonData, _ := json.MarshalIndent(load, "", "\t")

	mailServiceUrl := "http://mail:5004/send"

	res, err := http.Post(mailServiceUrl, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJson(w, err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("internal server error"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Email send to " + load.To

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) logItemAMQP(w http.ResponseWriter, load LogPayload) {

	err := app.pushToQueue(load.Name, load.Data)

	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: "Logged via RabbitMQ",
	}

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, data string) error {
	p, err := event.NewProducer(app.RabbitMQ)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: data,
	}

	// Don't use Indent in production
	j, _ := json.MarshalIndent(payload, "", "\t") // j -> []byte

	err = p.Push(string(j), "log.INFO")

	return err
}

type RPCPayload struct {
	Name, Data string
}

func (app *Config) LogItemRPC(w http.ResponseWriter, load LogPayload) {
	client, err := rpc.Dial("tcp", "logger:5031")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	payload := RPCPayload{
		Name: load.Name,
		Data: load.Data,
	}

	result := new(string)

	err = client.Call("RPCServer.LogInfo", payload, result)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	res := jsonResponse{
		Error:   false,
		Message: *result,
	}

	app.writeJson(w, http.StatusAccepted, res)
}
