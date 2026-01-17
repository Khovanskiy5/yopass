package repository

import (
	"os"
	"testing"

	"github.com/Khovanskiy5/yopass/internal/secret/domain"
)

func TestRedis(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Skip("Specify REDIS_URL env variable to test Redis database")
	}

	r, err := NewRedis(redisURL)
	if err != nil {
		t.Fatalf("error in NewRedis(): %v", err)
	}

	key := "f9fa5704-3ed2-4e60-b441-c426d3f9f3c1"
	secret := domain.Secret{Message: "foo", OneTime: true}

	err = r.Put(key, secret)
	if err != nil {
		t.Fatalf("error in Put(): %v", err)
	}

	storedVal, err := r.Get(key)
	if err != nil {
		t.Fatalf("error in Get(): %v", err)
	}

	if storedVal.Message != secret.Message {
		t.Fatalf("expected value %s, got %s", secret.Message, storedVal.Message)
	}

	_, err = r.Get(key)
	if err == nil {
		t.Fatal("expected error from Get() after Delete()")
	}
}

func TestRedisUnits(t *testing.T) {
	t.Run("NewRedis with invalid URL", func(t *testing.T) {
		_, err := NewRedis("invalid-url")
		if err == nil {
			t.Fatal("Expected error for invalid Redis URL")
		}
	})

	t.Run("NewRedis with valid URL", func(t *testing.T) {
		db, err := NewRedis("redis://localhost:6379/0")
		if err != nil {
			t.Fatalf("Expected no error for valid Redis URL, got: %v", err)
		}
		r, ok := db.(*Redis)
		if !ok {
			t.Fatal("NewRedis should return *Redis")
		}
		if r.client == nil {
			t.Fatal("Client should be initialized")
		}
	})
}

func TestRedisStatus(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Skip("Specify REDIS_URL env variable to test Redis database")
	}

	r, err := NewRedis(redisURL)
	if err != nil {
		t.Fatalf("error in NewRedis(): %v", err)
	}

	t.Run("Status returns correct OneTime value for existing secret", func(t *testing.T) {
		key := "test-status-onetime"
		secret := domain.Secret{Message: "test message", OneTime: true, Expiration: 3600}

		err := r.Put(key, secret)
		if err != nil {
			t.Fatalf("error in Put(): %v", err)
		}

		oneTime, err := r.Status(key)
		if err != nil {
			t.Fatalf("error in Status(): %v", err)
		}

		if oneTime != true {
			t.Fatalf("expected OneTime to be true, got %v", oneTime)
		}

		r.Delete(key)
	})

	t.Run("Status returns correct OneTime value for non-onetime secret", func(t *testing.T) {
		key := "test-status-multi"
		secret := domain.Secret{Message: "test message", OneTime: false, Expiration: 3600}

		err := r.Put(key, secret)
		if err != nil {
			t.Fatalf("error in Put(): %v", err)
		}

		oneTime, err := r.Status(key)
		if err != nil {
			t.Fatalf("error in Status(): %v", err)
		}

		if oneTime != false {
			t.Fatalf("expected OneTime to be false, got %v", oneTime)
		}

		r.Delete(key)
	})
}
