package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

// StartRegister returns a HTTP function that starts a Webauthn registration session.
func StartRegister(authn *webauthn.WebAuthn, webauthnStore *WebauthnStore, userStore *UserInMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			WriteJSONError(w, http.StatusUnprocessableEntity, "The body contains an unexpected value.")
			return
		}

		user := &User{
			ID:    body.ID,
			Email: body.Email,
		}

		if err := user.Validate(); err != nil {
			WriteJSONErrorf(w, http.StatusUnprocessableEntity, "The user %s.", err)
			return
		}

		u, err := userStore.Get(user.Email)
		if err != nil && !errors.Is(err, ErrNotFound) {
			log.Println(err)
			Write500JSONError(w)
			return
		}
		// User exists with credentials, return an error. We use a two step process
		// to create the user. Users might stop after the first step but we don't want
		// them to not be able to pursue the process later.
		if !errors.Is(err, ErrNotFound) && len(u.Credentials) > 0 {
			WriteJSONError(w, http.StatusUnprocessableEntity, "The user email already exists.")
			return
		}

		userStore.Set(user.Email, *user)

		options, data, err := authn.BeginRegistration(
			user,
			webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
				UserVerification: protocol.VerificationRequired,
			}),
		)
		if err != nil {
			log.Println(err)
			Write500JSONError(w)
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

// CompleteRegister returns a HTTP function that completes a Webauthn registration session.
func CompleteRegister(authn *webauthn.WebAuthn, webauthnStore *WebauthnStore, userStore *UserInMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Token string `json:"token"`
		}

		data, err := webauthnStore.Get(w, r)
		if err != nil {
			WriteJSONError(w, http.StatusBadRequest, "The registration didn't start for this user.")
			return
		}

		user, err := userStore.Get(data.ID)
		if errors.Is(err, ErrNotFound) {
			WriteJSONError(w, http.StatusBadRequest, "The registration didn't start for this user.")
			return
		}
		if err != nil {
			log.Println(err)
			Write500JSONError(w)
			return
		}

		credential, err := authn.FinishRegistration(&user, *data.Session, r)
		if err != nil {
			log.Println(err)
			Write500JSONError(w)
			return
		}

		user.Credentials = append(user.Credentials, *credential)

		userStore.Set(user.Email, user)

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
