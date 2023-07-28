package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(responseWriter http.ResponseWriter, code int, payload interface{}) {
	data, marshalError := json.Marshal(payload)

	if marshalError != nil {
		log.Printf("Failed to marshal JSON repsonse: %v", payload)
		responseWriter.WriteHeader(500)
		return
	}

	responseWriter.Header().Add("Content-Type", "application/json")
	responseWriter.WriteHeader(code)
	responseWriter.Write(data)
}

func respondWithError(responseWriter http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(responseWriter, code, errResponse{Error: msg})
}
