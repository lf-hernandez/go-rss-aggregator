package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lf-hernandez/go-rss-aggregator/internal/database"
)

func (apiConfiguration *apiConfig) handlerCreateUser(responseWriter http.ResponseWriter, r *http.Request) {
	type parameteres struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameteres{}
	decoderError := decoder.Decode(&params)

	if decoderError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintf("Error parsing JSON: %v", decoderError))
		return
	}

	user, createUserError := apiConfiguration.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if createUserError != nil {
		respondWithError(responseWriter, 400, fmt.Sprint("Couldn't create users:, ", createUserError))
		return
	}

	respondWithJSON(responseWriter, 200, databaseUserToUser(user))
}
