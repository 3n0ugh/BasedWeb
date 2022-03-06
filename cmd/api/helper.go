package main

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]interface{}

// convert data struct to json type and write into response
func writeJSON(w http.ResponseWriter, status int, data envelope, header http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for k, v := range header {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	if err != nil {
		return err
	}
	return nil
}
