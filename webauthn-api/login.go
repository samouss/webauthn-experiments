package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

// StartLogin returns a HTTP function that starts a Webauthn login session.
func StartLogin(authn *webauthn.WebAuthn, webauthnStore *WebauthnStore, userStore *UserInMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Email string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			WriteJSONError(w, http.StatusUnprocessableEntity, "The body contains an unexpected value.")
			return
		}

		user, err := userStore.Get(body.Email)
		if errors.Is(err, ErrNotFound) {
			WriteJSONError(w, http.StatusUnauthorized, "Invalid credentials.")
			return
		}
		if err != nil {
			log.Println(err)
			Write500JSONError(w)
			return
		}

		options, data, err := authn.BeginLogin(
			&user,
			webauthn.WithUserVerification(protocol.VerificationRequired),
		)
		if err != nil {
			WriteJSONError(w, http.StatusUnauthorized, "Invalid credentials.")
			return
		}

		err = webauthnStore.Set(w, r, &WebauthnData{
			ID:      user.Email,
			Session: data,
		})
		if err != nil {
			log.Println(err)
			Write500JSONError(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(options); err != nil {
			log.Println(err)
			Write500JSONError(w)
			return
		}
	}
}

// CompleteLogin returns a HTTP function that completes a Webauthn login session.
func CompleteLogin(authn *webauthn.WebAuthn, webauthnStore *WebauthnStore, userStore *UserInMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Token string `json:"token"`
		}

		data, err := webauthnStore.Get(w, r)
		if err != nil {
			WriteJSONError(w, http.StatusBadRequest, "The login didn't start for this user.")
			return
		}

		user, err := userStore.Get(data.ID)
		if errors.Is(err, ErrNotFound) {
			WriteJSONError(w, http.StatusUnauthorized, "Invalid credentials.")
			return
		}
		if err != nil {
			log.Println(err)
			Write500JSONError(w)
			return
		}

		if _, err := authn.FinishLogin(&user, *data.Session, r); err != nil {
			WriteJSONError(w, http.StatusUnauthorized, "Invalid credentials.")
			return
		}

		token := Token{UserID: user.Email}
		encoded, err := token.Encode()
		if err != nil {
			log.Println(err)
			Write500JSONError(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response{Token: encoded}); err != nil {
			log.Println(err)
			Write500JSONError(w)
			return
		}
	}
}
