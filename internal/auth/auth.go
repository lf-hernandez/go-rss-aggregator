package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("no authentication data found")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", errors.New("bad auth header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("bad auth header")
	}

	return vals[1], nil
}
