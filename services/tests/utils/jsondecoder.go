// Package utils contains utility functions
package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

type JsonMarshalError struct {
	Msg        string
	StatusCode int
}

func (mr *JsonMarshalError) Error() string { return mr.Msg }

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) *JsonMarshalError {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			return &JsonMarshalError{
				Msg:        "content-type header is not application/json",
				StatusCode: http.StatusUnsupportedMediaType,
			}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return &JsonMarshalError{
				Msg:        fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset),
				StatusCode: http.StatusBadRequest,
			}
		case errors.Is(err, io.ErrUnexpectedEOF):
			return &JsonMarshalError{
				Msg:        "request body contains badly-formed JSON",
				StatusCode: http.StatusBadRequest,
			}
		case errors.As(err, &unmarshalTypeError):
			return &JsonMarshalError{
				Msg:        fmt.Sprintf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset),
				StatusCode: http.StatusBadRequest,
			}
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return &JsonMarshalError{
				Msg:        fmt.Sprintf("request body contains unknown field %s", fieldName),
				StatusCode: http.StatusBadRequest,
			}
		case errors.Is(err, io.EOF):
			return &JsonMarshalError{
				Msg:        "request body must not be empty",
				StatusCode: http.StatusBadRequest,
			}
		case err.Error() == "http: request body too large":
			return &JsonMarshalError{
				Msg:        "request body must not be larger than 1MB",
				StatusCode: http.StatusRequestEntityTooLarge,
			}
		default:
			return &JsonMarshalError{err.Error(), http.StatusInternalServerError}
		}
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return &JsonMarshalError{
			Msg:        "request body must only contain a single JSON object",
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil
}
