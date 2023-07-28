package main

import "net/http"

func handlerError(responseWriter http.ResponseWriter, t *http.Request) {
	respondWithError(responseWriter, 400, "Something went wrong")
}
