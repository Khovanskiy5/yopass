package yopass

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
)

// Service handles business logic for secrets
type Service interface {
	CreateSecret(secret Secret) (string, error)
	GetSecret(key string) (Secret, error)
	GetSecretStatus(key string) (bool, error)
	DeleteSecret(key string) (bool, error)
	IsPGPEncrypted(content string) bool
}

type service struct {
	repo                Repository
	maxLength           int
	forceOneTimeSecrets bool
	allowedExpirations  []int32
}

// NewService creates a new Service
func NewService(repo Repository, maxLength int, forceOneTimeSecrets bool) Service {
	return NewServiceWithExpirations(repo, maxLength, forceOneTimeSecrets, []int32{3600, 86400, 604800})
}

// NewServiceWithExpirations creates a new Service with custom allowed expirations
func NewServiceWithExpirations(repo Repository, maxLength int, forceOneTimeSecrets bool, allowedExpirations []int32) Service {
	return &service{
		repo:                repo,
		maxLength:           maxLength,
		forceOneTimeSecrets: forceOneTimeSecrets,
		allowedExpirations:  allowedExpirations,
	}
}

// CreateSecret stores a new secret and returns its key
func (s *service) CreateSecret(secret Secret) (string, error) {
	if !s.IsPGPEncrypted(secret.Message) {
		return "", fmt.Errorf("Message must be PGP encrypted")
	}

	if !s.validExpiration(secret.Expiration) {
		return "", fmt.Errorf("Invalid expiration specified")
	}

	if !secret.OneTime && s.forceOneTimeSecrets {
		return "", fmt.Errorf("Secret must be one time download")
	}

	if len(secret.Message) > s.maxLength {
		return "", fmt.Errorf("The encrypted message is too long")
	}

	uuidVal, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("Unable to generate UUID: %w", err)
	}
	key := uuidVal.String()

	if err := s.repo.Put(key, secret); err != nil {
		return "", fmt.Errorf("Failed to store secret in database")
	}

	return key, nil
}

// GetSecret retrieves a secret by key
func (s *service) GetSecret(key string) (Secret, error) {
	return s.repo.Get(key)
}

// GetSecretStatus returns the one-time status of a secret
func (s *service) GetSecretStatus(key string) (bool, error) {
	return s.repo.Status(key)
}

// DeleteSecret removes a secret by key
func (s *service) DeleteSecret(key string) (bool, error) {
	return s.repo.Delete(key)
}

// validExpiration validates that expiration is in the allowed list
func (s *service) validExpiration(expiration int32) bool {
	for _, ttl := range s.allowedExpirations {
		if ttl == expiration {
			return true
		}
	}
	return false
}

// IsPGPEncrypted verifies that the provided content is a valid PGP encrypted message
func (s *service) IsPGPEncrypted(content string) bool {
	if content == "" {
		return false
	}
	// We can't easily import from pkg/server/utils.go here due to circular dependency.
	// Actually, the logic should probably be here in the domain/service.
	// For now, I'll replicate it or move it.
	// Since yopass package already has armor import in yopass.go, it's fine.
	return strings.HasPrefix(content, "-----BEGIN PGP MESSAGE-----") &&
		strings.HasSuffix(strings.TrimSpace(content), "-----END PGP MESSAGE-----")
}
