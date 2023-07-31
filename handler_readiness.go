package main

import "net/http"

func handlerReadiness(responseWriter http.ResponseWriter, request *http.Request) {
	respondWithJSON(responseWriter, 200, struct{}{})
}
