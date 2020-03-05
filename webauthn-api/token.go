package main

import (
	"encoding/base64"
	"encoding/json"
)

// Token represents the value used to authenticate a request.
type Token struct {
	UserID string `json:"userID"`
}

// Encode returns an encoded representation of the struct.
func (t *Token) Encode() (string, error) {
	raw, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(raw), nil
}

// Decode hydrates the target from the given input.
func (t *Token) Decode(input string) error {
	raw, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw, t); err != nil {
		return err
	}

	return nil
}
