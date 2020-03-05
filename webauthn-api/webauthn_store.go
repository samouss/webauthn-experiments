package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/sessions"
)

const sessionKey = "auth"

// WebauthnStore represents the session store for the Webauthn.
type WebauthnStore struct {
	store  sessions.Store
	secure bool
}

// WebauthnData represents the session data required by the Webauthn.
type WebauthnData struct {
	ID      string                `json:"id"`
	Session *webauthn.SessionData `json:"session"`
}

// Get retrieves the current session from the request, if none an error is returned.
func (ws *WebauthnStore) Get(w http.ResponseWriter, r *http.Request) (*WebauthnData, error) {
	session, err := ws.store.Get(r, sessionKey)
	if err != nil {
		return nil, fmt.Errorf("get: get session: %w", err)
	}

	// Cleanup the session
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return nil, fmt.Errorf("get: delete session: %w", err)
	}

	v, ok := session.Values["data"]
	if !ok {
		return nil, fmt.Errorf("get: empty session")
	}

	raw, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("get: corrupt session")
	}

	var data WebauthnData
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &data, nil
}

// Set stores the session data within the store.
func (ws *WebauthnStore) Set(w http.ResponseWriter, r *http.Request, data *WebauthnData) error {
	session, err := ws.store.Get(r, sessionKey)
	if err != nil {
		return fmt.Errorf("set: get session: %w", err)
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("set: marshal: %w", err)
	}

	session.Options.Secure = ws.secure
	session.Options.HttpOnly = true
	session.Options.MaxAge = 0
	session.Values["data"] = raw

	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("set: save session: %w", err)
	}

	return nil
}
