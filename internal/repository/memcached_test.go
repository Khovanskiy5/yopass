package repository

import (
	"os"
	"testing"

	"github.com/Khovanskiy5/yopass/internal/secret/domain"
)

func TestMemcached(t *testing.T) {
	memcachedURL := os.Getenv("MEMCACHED")
	if memcachedURL == "" {
		t.Skip("Specify MEMCACHED env variable to test memcached database")
	}

	m := NewMemcached(memcachedURL)

	key := "f9fa5704-3ed2-4e60-b441-c426d3f9f3c1"
	secret := domain.Secret{Message: "foo", OneTime: true}

	err := m.Put(key, secret)
	if err != nil {
		t.Fatalf("error in Put(): %v", err)
	}

	storedSecret, err := m.Get(key)
	if err != nil {
		t.Fatalf("error in Get(): %v", err)
	}

	if storedSecret.Message != secret.Message {
		t.Fatalf("expected value %s, got %s", secret.Message, storedSecret.Message)
	}

	_, err = m.Get(key)
	if err == nil {
		t.Fatal("expected error from Get() after Delete()")
	}
}

func TestMemcachedUnits(t *testing.T) {
	t.Run("NewMemcached creates correct instance", func(t *testing.T) {
		db := NewMemcached("localhost:11211")
		m, ok := db.(*Memcached)
		if !ok {
			t.Fatal("NewMemcached should return *Memcached")
		}
		if m.client == nil {
			t.Fatal("Client should be initialized")
		}
	})
}

func TestMemcachedStatus(t *testing.T) {
	memcachedURL := os.Getenv("MEMCACHED")
	if memcachedURL == "" {
		t.Skip("Specify MEMCACHED env variable to test memcached database")
	}

	m := NewMemcached(memcachedURL)

	t.Run("Status returns correct OneTime value for existing secret", func(t *testing.T) {
		key := "test-status-onetime"
		secret := domain.Secret{Message: "test message", OneTime: true, Expiration: 3600}

		err := m.Put(key, secret)
		if err != nil {
			t.Fatalf("error in Put(): %v", err)
		}

		oneTime, err := m.Status(key)
		if err != nil {
			t.Fatalf("error in Status(): %v", err)
		}

		if oneTime != true {
			t.Fatalf("expected OneTime to be true, got %v", oneTime)
		}

		m.Delete(key)
	})

	t.Run("Status returns correct OneTime value for non-onetime secret", func(t *testing.T) {
		key := "test-status-multi"
		secret := domain.Secret{Message: "test message", OneTime: false, Expiration: 3600}

		err := m.Put(key, secret)
		if err != nil {
			t.Fatalf("error in Put(): %v", err)
		}

		oneTime, err := m.Status(key)
		if err != nil {
			t.Fatalf("error in Status(): %v", err)
		}

		if oneTime != false {
			t.Fatalf("expected OneTime to be false, got %v", oneTime)
		}

		m.Delete(key)
	})

	t.Run("Status returns error for non-existent key", func(t *testing.T) {
		_, err := m.Status("non-existent")
		if err != domain.ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}
