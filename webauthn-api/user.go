package main

import (
	"errors"

	"github.com/duo-labs/webauthn/webauthn"
)

// User represents a domain user.
type User struct {
	ID          string                `json:"id"`
	Email       string                `json:"email"`
	Credentials []webauthn.Credential `json:"-"`
}

// Validate validates the user.
func (u *User) Validate() error {
	if u.ID == "" {
		return errors.New("ID is required")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	return nil
}

// WebAuthnID returns the user ID.
func (u *User) WebAuthnID() []byte {
	return []byte(u.ID)
}

// WebAuthnName returns the user name.
func (u *User) WebAuthnName() string {
	return u.Email
}

// WebAuthnDisplayName returns the user name.
func (u *User) WebAuthnDisplayName() string {
	return u.Email
}

// WebAuthnIcon returns the user icon URL.
func (u *User) WebAuthnIcon() string {
	return ""
}

// WebAuthnCredentials returns the credentials for the user.
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}
