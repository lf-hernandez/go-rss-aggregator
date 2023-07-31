package main

import (
	"fmt"
	"net/http"

	"github.com/lf-hernandez/go-rss-aggregator/internal/auth"
	"github.com/lf-hernandez/go-rss-aggregator/internal/database"
)

type authenticatedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authenticatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		apiKey, apiKeyError := auth.GetAPIKey(request.Header)

		if apiKeyError != nil {
			respondWithError(w, 403, fmt.Sprintln("Auth error: ", apiKeyError))

			cfg.DB.GetUserByAPIKey(request.Context(), apiKey)
		}

		user, getUserError := cfg.DB.GetUserByAPIKey(request.Context(), apiKey)

		if getUserError != nil {
			respondWithError(w, 400, fmt.Sprintln("Error retrieving user: ", getUserError))
			return
		}

		handler(w, request, user)
	}
}
