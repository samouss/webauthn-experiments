package main

import "fmt"

// UserInMemoryStore represents an in-memory data store for users.
type UserInMemoryStore struct {
	store *InMemoryStore
}

// Get retrieves the user for the given ID from the store.
func (s *UserInMemoryStore) Get(ID string) (User, error) {
	u, err := s.store.Get(ID)
	if err != nil {
		return User{}, err
	}

	user, uok := u.(User)
	if !uok {
		return User{}, fmt.Errorf("unexpected type")
	}

	return user, nil
}

// Set stores the given ID/user pair in the store.
func (s *UserInMemoryStore) Set(ID string, user User) {
	s.store.Set(ID, user)
}
