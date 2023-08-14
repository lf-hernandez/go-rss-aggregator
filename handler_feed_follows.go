package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/lf-hernandez/go-rss-aggregator/internal/database"
	"github.com/lf-hernandez/go-rss-aggregator/models"
)

func (apiConfiguration *apiConfig) handlerCreateFeedFollow(responseWriter http.ResponseWriter, request *http.Request, authenticatedUser database.User) {
	type parameteres struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(request.Body)

	params := parameteres{}
	decoderError := decoder.Decode(&params)

	if decoderError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintln("Error parsing JSON: ", decoderError))
		return
	}

	feedFollow, createFeedError := apiConfiguration.Database.CreateFeedFollow(request.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    authenticatedUser.ID,
		FeedID:    params.FeedID,
	})

	if createFeedError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintln("Couldn't create feed follow:, ", createFeedError))
		return
	}

	respondWithJSON(responseWriter, 201, models.DatabaseFeedFollowToFeedFollow(feedFollow))
}

func (apiConfiguration *apiConfig) handlerGetFeedFollows(responseWriter http.ResponseWriter, request *http.Request, authenticatedUser database.User) {
	feedFollows, getFeedsError := apiConfiguration.Database.GetFeedFollows(request.Context(), authenticatedUser.ID)
	if getFeedsError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintln("Couldn't get feeds follows, ", getFeedsError))
		return
	}

	respondWithJSON(responseWriter, 200, models.DatabaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(responseWriter http.ResponseWriter, request *http.Request, authenticatedUser database.User) {
	feedFollowIDString := chi.URLParam(request, "feedFollowID")
	feedFollowID, uuidParsingError := uuid.Parse(feedFollowIDString)
	if uuidParsingError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintln("Couldn't parse feed follow id: ", uuidParsingError))
	}

	databaseDeleteError := apiCfg.Database.DeleteFeedFollow(request.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: authenticatedUser.ID,
	})

	if databaseDeleteError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintln("Couldn't delete feed follow: ", databaseDeleteError))
	}

	respondOk(responseWriter)
}
