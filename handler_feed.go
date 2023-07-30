package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lf-hernandez/go-rss-aggregator/internal/database"
)

func (apiConfiguration *apiConfig) handlerCreateFeed(responseWriter http.ResponseWriter, r *http.Request, user database.User) {
	type parameteres struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameteres{}
	decoderError := decoder.Decode(&params)

	if decoderError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintf("Error parsing JSON: %v", decoderError))
		return
	}

	feed, createFeedError := apiConfiguration.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if createFeedError != nil {
		respondWithError(responseWriter, 400, fmt.Sprint("Couldn't create users:, ", createFeedError))
		return
	}

	respondWithJSON(responseWriter, 201, databaseFeedToFeed(feed))
}

func (apiConfiguration *apiConfig) handlerGetFeeds(responseWriter http.ResponseWriter, r *http.Request) {
	feeds, getFeedsError := apiConfiguration.DB.GetFeeds(r.Context())
	if getFeedsError != nil {
		respondWithError(responseWriter, 400, fmt.Sprint("Couldn't get feeds:, ", getFeedsError))
		return
	}

	respondWithJSON(responseWriter, 201, databaseFeedsToFeeds(feeds))
}
