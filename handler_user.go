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
		respondWithError(responseWriter, 400, fmt.Sprintln("Error parsing JSON:", decoderError))
		return
	}

	user, createUserError := apiConfiguration.Database.CreateUser(request.Context(), database.CreateUserParams{
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

func (apiConfiguration *apiConfig) handlerGetPostsByUser(responseWriter http.ResponseWriter, request *http.Request, user database.User) {
	posts, getPostsDBError := apiConfiguration.Database.GetPostsByUser(request.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  10,
	})

	if getPostsDBError != nil {
		respondWithError(responseWriter, 400, fmt.Sprintf("Couln't get posts: %v", getPostsDBError))
		return
	}

	respondWithJSON(responseWriter, 200, databasePostsToPosts(posts))

}
