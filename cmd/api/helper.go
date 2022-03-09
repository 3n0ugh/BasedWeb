package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type envelope map[string]interface{}

// take the enveloped map data and marshall it and return it.
func (app *application) prettyJSON(data envelope) ([]byte, error) {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil, err
	}

	js = append(js, '\n')

	return js, err
}

// take the json data and write it into response
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, header http.Header) error {
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

// read json data and send custom error message if there is an error
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// max limit size of byte reading
	maxBytes := 1_048_576

	// limit the size of body reader
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// show unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// parse json to struct
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)",
				syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)",
				unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unkown key %s", fieldName)
		case errors.Is(err, errors.New("http: request body too large")):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	// if we can decode again, it's mean there are more than one json value
	// which we don't want, so we return custom error message
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
