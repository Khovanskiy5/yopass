package service

import (
	"fmt"
	"strings"

	"github.com/Khovanskiy5/yopass/internal/secret/domain"
	"github.com/gofrs/uuid"
)

type SecretService interface {
	CreateSecret(secret domain.Secret) (string, error)
	GetSecret(key string) (domain.Secret, error)
	GetSecretStatus(key string) (bool, error)
	DeleteSecret(key string) (bool, error)
}

type secretService struct {
	repo                domain.Repository
	maxLength           int
	forceOneTimeSecrets bool
	allowedExpirations  []int32
}

func NewSecretService(
	repo domain.Repository,
	maxLength int,
	forceOneTimeSecrets bool,
	allowedExpirations []int32,
) SecretService {
	return &secretService{
		repo:                repo,
		maxLength:           maxLength,
		forceOneTimeSecrets: forceOneTimeSecrets,
		allowedExpirations:  allowedExpirations,
	}
}

func (s *secretService) CreateSecret(secret domain.Secret) (string, error) {
	if !s.isPGPEncrypted(secret.Message) {
		return "", fmt.Errorf("message must be PGP encrypted")
	}

	if !s.isValidExpiration(secret.Expiration) {
		return "", fmt.Errorf("invalid expiration specified")
	}

	if !secret.OneTime && s.forceOneTimeSecrets {
		return "", fmt.Errorf("secret must be one time download")
	}

	if len(secret.Message) > s.maxLength {
		return "", fmt.Errorf("the encrypted message is too long")
	}

	uuidVal, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("unable to generate UUID: %w", err)
	}
	key := uuidVal.String()

	if err := s.repo.Put(key, secret); err != nil {
		return "", fmt.Errorf("failed to store secret in database")
	}

	return key, nil
}

func (s *secretService) GetSecret(key string) (domain.Secret, error) {
	return s.repo.Get(key)
}

func (s *secretService) GetSecretStatus(key string) (bool, error) {
	return s.repo.Status(key)
}

func (s *secretService) DeleteSecret(key string) (bool, error) {
	return s.repo.Delete(key)
}

func (s *secretService) isValidExpiration(expiration int32) bool {
	for _, ttl := range s.allowedExpirations {
		if ttl == expiration {
			return true
		}
	}
	return false
}

func (s *secretService) isPGPEncrypted(content string) bool {
	if content == "" {
		return false
	}
	return strings.HasPrefix(content, "-----BEGIN PGP MESSAGE-----") &&
		strings.HasSuffix(strings.TrimSpace(content), "-----END PGP MESSAGE-----")
}
