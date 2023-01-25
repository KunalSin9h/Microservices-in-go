package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

/*
General Json Response
*/
type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(data); err != nil {
		return err
	}

	/*
		struct{} -> empty struct
		struct{}{} -> instance of empty struct
		means nothing
	*/
	if err := dec.Decode(&struct{}{}); err != nil {
		/*
			After extracting `data` there should nothing left
		*/
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func (app *Config) writeJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	// Json -> bytes array
	outBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, val := range headers[0] {
			w.Header()[key] = val
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(outBytes)

	if err != nil {
		return err
	} else {
		return nil
	}

}

func (app *Config) errorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var response jsonResponse
	response.Error = true
	response.Message = err.Error()

	return app.writeJson(w, statusCode, response)
}
