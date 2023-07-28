package main

import "net/http"

func handlerReadiness(responseWriter http.ResponseWriter, t *http.Request) {
	respondWithJSON(responseWriter, 200, struct{}{})
}
