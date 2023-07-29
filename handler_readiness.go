package main

import "net/http"

func handlerReadiness(responseWriter http.ResponseWriter, r *http.Request) {
	respondWithJSON(responseWriter, 200, struct{}{})
}
