package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type envelope map[string]any

// helper function for writing JSON responses
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Marshal JSON from the passed in data with indentation nesting for readability
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// append a newline to the JSON object for easier readability
	js = append(js, '\n')

	// set headers from a passed in set of key-value pairs
	for key, value := range headers {
		w.Header()[key] = value
	}

	// set the content header for JSON, write the status header, and the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// helper function for reading in JSON responses
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// sets maximum allowed bytes of read in JSON
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// blocks unknown JSON fields from being passed into decoder
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// decode the JSON data using the decoder object
	if err := dec.Decode(dst); err != nil {
		return err
	}
	// ensure there is only one JSON object in the request body
	err := dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}
