package main

import (
	"fmt"
	"net/http"

	"github.com/lf-hernandez/go-rss-aggregator/internal/auth"
	"github.com/lf-hernandez/go-rss-aggregator/internal/database"
)

type authenticatedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authenticatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, apiKeyError := auth.GetAPIKey(r.Header)

		if apiKeyError != nil {
			respondWithError(w, 403, fmt.Sprint("Auth error: ", apiKeyError))

			cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		}

		user, getUserError := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)

		if getUserError != nil {
			respondWithError(w, 400, fmt.Sprint("Error retrieving user: ", getUserError))
			return
		}

		handler(w, r, user)
	}
}
