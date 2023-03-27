// Filename: cmd/api/helpers.go

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

// func (app *application) readIDParam(w http.ResponseWriter, r *http.Request) (int64, error) {
func (app *application) readIDParam(r *http.Request) (int64, error) {
	// Use the "ParamsFromContext()" function to get the request context as a slice
	params := httprouter.ParamsFromContext(r.Context())
	// Get the value of the "id" parameter
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) WriteJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Convert our map into a JSON object
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	// Add a newline to make viewing on the terminal easier
	js = append(js, '\n')
	// Add the headers
	for key, value := range headers {
		w.Header()[key] = value
	}
	// Specifiy that we will serve our responses using JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Write the []byte slice containing the JSON response body
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, desitination interface{}) error {
	// specify the max size of our JSON request body
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	// we will try to decode the JSON request
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	// start the decoding
	err := dec.Decode(desitination)
	if err != nil {
		//something went wrong
		// let's find out the type of error
		var syntaxError *json.SyntaxError
		var unMarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		// check for max bytes
		var maxBytesError *http.MaxBytesError
		// let's check for the type of decode error
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON(at character  %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unMarshalTypeError):
			// whixk field has the wrong type
			if unMarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unMarshalTypeError.Field)
			}
			return fmt.Errorf("body contains badly-formed JSON(at character %d)", unMarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
			// unmappable field/ not -ecistent field
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}

	}
	// Let's call the decoder agian to check if ther are any
	// trailing json objects
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}
