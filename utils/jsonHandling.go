package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must only have a single json value")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	if len(headers) > 0 {
		for key, val := range headers[0] {
			w.Header()[key] = val
		}
	}
	
	w.Header().Set("Content-type", "application-json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	return err
}