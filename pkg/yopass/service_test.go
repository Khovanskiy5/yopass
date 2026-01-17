package yopass

import (
	"testing"
)

type mockRepo struct {
	putErr error
	getErr error
	secret Secret
}

func (m *mockRepo) Get(key string) (Secret, error) {
	return m.secret, m.getErr
}
func (m *mockRepo) Put(key string, secret Secret) error {
	return m.putErr
}
func (m *mockRepo) Delete(key string) (bool, error) {
	return true, nil
}
func (m *mockRepo) Status(key string) (bool, error) {
	return m.secret.OneTime, nil
}

func TestCreateSecret(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo, 100, false)

	tests := []struct {
		name    string
		secret  Secret
		wantErr bool
	}{
		{
			name: "Valid secret",
			secret: Secret{
				Message:    "-----BEGIN PGP MESSAGE-----\n...\n-----END PGP MESSAGE-----",
				Expiration: 3600,
			},
			wantErr: false,
		},
		{
			name: "Invalid PGP",
			secret: Secret{
				Message:    "plain text",
				Expiration: 3600,
			},
			wantErr: true,
		},
		{
			name: "Invalid expiration",
			secret: Secret{
				Message:    "-----BEGIN PGP MESSAGE-----\n...\n-----END PGP MESSAGE-----",
				Expiration: 123,
			},
			wantErr: true,
		},
		{
			name: "Too long message",
			secret: Secret{
				Message:    "-----BEGIN PGP MESSAGE-----\n" + string(make([]byte, 200)) + "\n-----END PGP MESSAGE-----",
				Expiration: 3600,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateSecret(tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCustomExpirations(t *testing.T) {
	repo := &mockRepo{}
	svc := NewServiceWithExpirations(repo, 100, false, []int32{60})

	s := Secret{
		Message:    "-----BEGIN PGP MESSAGE-----\n...\n-----END PGP MESSAGE-----",
		Expiration: 60,
	}
	if _, err := svc.CreateSecret(s); err != nil {
		t.Errorf("Expected success for 60s expiration, got %v", err)
	}

	s.Expiration = 3600
	if _, err := svc.CreateSecret(s); err == nil {
		t.Error("Expected error for 3600s expiration when only 60s is allowed")
	}
}

func TestForceOneTime(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo, 100, true)

	s := Secret{
		Message:    "-----BEGIN PGP MESSAGE-----\n...\n-----END PGP MESSAGE-----",
		Expiration: 3600,
		OneTime:    false,
	}
	if _, err := svc.CreateSecret(s); err == nil {
		t.Error("Expected error when forceOneTimeSecrets is true but secret is not one-time")
	}

	s.OneTime = true
	if _, err := svc.CreateSecret(s); err != nil {
		t.Errorf("Expected success when secret is one-time, got %v", err)
	}
}
