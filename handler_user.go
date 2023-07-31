package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lf-hernandez/go-rss-aggregator/internal/database"
)

func (apiConfiguration *apiConfig) handlerCreateUser(responseWriter http.ResponseWriter, request *http.Request) {
	type parameteres struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(request.Body)

	params := parameteres{}
	decoderError := decoder.Decode(&params)

	if decoderError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintln("Error parsing JSON: %v", decoderError))
		return
	}

	user, createUserError := apiConfiguration.DB.CreateUser(request.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if createUserError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintln("Couldn't create users:, ", createUserError))
		return
	}

	respondWithJSON(responseWriter, 201, databaseUserToUser(user))
}

func (apiConfiguration *apiConfig) handlerGetUser(responseWriter http.ResponseWriter, request *http.Request, user database.User) {
	respondWithJSON(responseWriter, 200, databaseUserToUser(user))
}
