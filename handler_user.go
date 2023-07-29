package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lf-hernandez/go-rss-aggregator/internal/auth"
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

	respondWithJSON(responseWriter, 201, databaseUserToUser(user))
}

func (apiConfiguration *apiConfig) handlerGetUser(responseWriter http.ResponseWriter, r *http.Request) {
	apiKey, apiKeyError := auth.GetAPIKey(r.Header)

	if apiKeyError != nil {
		respondWithError(responseWriter, 403, fmt.Sprint("Auth error: ", apiKeyError))

		apiConfiguration.DB.GetUserByAPIKey(r.Context(), apiKey)
	}

	user, getUserError := apiConfiguration.DB.GetUserByAPIKey(r.Context(), apiKey)

	if getUserError != nil {
		respondWithError(responseWriter, 400, fmt.Sprint("Error retrieving user: ", getUserError))
		return
	}

	respondWithJSON(responseWriter, 200, databaseUserToUser(user))
}
