package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lf-hernandez/go-rss-aggregator/internal/database"
)

func (apiConfiguration *apiConfig) handlerCreateFeed(responseWriter http.ResponseWriter, request *http.Request, user database.User) {
	type parameteres struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(request.Body)

	params := parameteres{}
	decoderError := decoder.Decode(&params)

	if decoderError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintln("Error parsing JSON:", decoderError))
		return
	}

	feed, createFeedError := apiConfiguration.Database.CreateFeed(request.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if createFeedError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintln("Couldn't create users:, ", createFeedError))
		return
	}

	respondWithJSON(responseWriter, 201, databaseFeedToFeed(feed))
}

func (apiConfiguration *apiConfig) handlerGetFeeds(responseWriter http.ResponseWriter, request *http.Request) {
	feeds, getFeedsError := apiConfiguration.Database.GetFeeds(request.Context())
	if getFeedsError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintln("Couldn't get feeds:, ", getFeedsError))
		return
	}

	respondWithJSON(responseWriter, 200, databaseFeedsToFeeds(feeds))
}
