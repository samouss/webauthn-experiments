package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

const authTokenHeader = "Authorization"
const authTokenPrefix = "Bearer "

type auth struct {
	handler   http.Handler
	userStore *UserInMemoryStore
}

// Auth returns a HTTP middleware that validates a Bearer authentication header.
func Auth(userStore *UserInMemoryStore) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &auth{
			handler:   handler,
			userStore: userStore,
		}
	}
}

func (a *auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get(authTokenHeader)
	if !strings.HasPrefix(header, authTokenPrefix) {
		WriteJSONErrorf(w, http.StatusUnauthorized, "The %q header is required with a %q format.", authTokenHeader, strings.TrimSpace(authTokenPrefix))
		return
	}

	token := &Token{}
	if err := token.Decode(strings.TrimPrefix(header, authTokenPrefix)); err != nil {
		WriteJSONErrorf(w, http.StatusUnauthorized, "The %q header is expected to be base64 encoded.", authTokenHeader)
		return
	}

	// We don't have a real authentication mechanism for the experiements. We check
	// that the given user exists in the store. In real life we should have used
	// a JWT or similar to validate the request.
	if _, err := a.userStore.Get(token.UserID); err != nil {
		if errors.Is(err, ErrNotFound) {
			WriteJSONError(w, http.StatusUnauthorized, "Invalid credentials.")
			return
		}

		log.Println(err)
		Write500JSONError(w)
		return
	}

	a.handler.ServeHTTP(w, r)
}
