package yopass

import "errors"

// ErrNotFound is returned when a secret is not found in the repository
var ErrNotFound = errors.New("secret not found")

// Repository interface for secret storage
type Repository interface {
	// Get returns the secret for the given key
	Get(key string) (Secret, error)
	// Put stores the secret for the given key
	Put(key string, secret Secret) error
	// Delete removes the secret for the given key
	Delete(key string) (bool, error)
	// Status returns whether the secret exists and if it is one-time
	Status(key string) (bool, error)
}
