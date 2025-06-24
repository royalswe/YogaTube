package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type HTTPError struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_message"`
}

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	// w.Header().Add("Content-Type", "application/json")
	// w.WriteHeader(status)
	// return json.NewEncoder(w).Encode(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return err
	}
	if _, err := w.Write(response); err != nil {
		fmt.Printf("Failed to write response: %v\n", err)
		return err
	}
	return nil
}

// Error handling for HTTTP requests

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, HTTPError{StatusCode: status, ErrorMessage: err.Error()})
}

// WrapError creates a new HTTPError with a status code
func WrapError(err error, statusCode int) *HTTPError {
	return &HTTPError{StatusCode: statusCode, ErrorMessage: err.Error()}
}

// Error implements the error interface for HTTPError
func (e *HTTPError) Error() string {
	return e.ErrorMessage
}

func MakeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			status := http.StatusBadRequest // Default status code
			if httpErr, ok := err.(*HTTPError); ok {
				status = httpErr.StatusCode
			}
			WriteError(w, status, err)
		}
	}
}
