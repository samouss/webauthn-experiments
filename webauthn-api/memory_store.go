package main

import (
	"errors"
	"sync"
)

// ErrNotFound -
var ErrNotFound = errors.New("key not found")

// InMemoryStore represents an in-memory data store.
type InMemoryStore struct {
	data  map[string]interface{}
	mutex sync.Mutex
}

// NewInMemoryStore create an instance of an in-memory data store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: map[string]interface{}{}}
}

// Get retrieves the given key from the in-memory data store.
// The function Get is thread safe.
func (m *InMemoryStore) Get(key string) (interface{}, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	value, ok := m.data[key]
	if !ok {
		return nil, ErrNotFound
	}

	return value, nil
}

// Set stores the given key/value pair in the in-memory data store.
// The function Set is thread safe.
func (m *InMemoryStore) Set(key string, value interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = value
}
