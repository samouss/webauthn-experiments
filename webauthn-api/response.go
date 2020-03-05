package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Write500JSONError sends a 500 with a generic message.
func Write500JSONError(w http.ResponseWriter) {
	WriteJSONError(w, http.StatusUnauthorized, http.StatusText(http.StatusInternalServerError)+".")
}

// WriteJSONErrorf sends a response with the given code and format specifier.
func WriteJSONErrorf(w http.ResponseWriter, status int, format string, values ...interface{}) {
	WriteJSONError(w, status, fmt.Sprintf(format, values...))
}

// WriteJSONError sends a response with the given code and message.
func WriteJSONError(w http.ResponseWriter, status int, message string) {
	type response struct {
		Message string `json:"message"`
	}

	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&response{Message: message}); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
