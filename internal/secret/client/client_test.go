package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Khovanskiy5/yopass/internal/secret/domain"
)

func TestFetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/secret/test-id" {
			t.Errorf("Expected path /secret/test-id, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(serverResponse{Message: "decrypted-content"})
	}))
	defer ts.Close()

	got, err := Fetch(ts.URL, "test-id")
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}
	if got != "decrypted-content" {
		t.Errorf("Expected decrypted-content, got %s", got)
	}
}

func TestStore(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		var s domain.Secret
		json.NewDecoder(r.Body).Decode(&s)
		if s.Message != "encrypted-content" {
			t.Errorf("Expected encrypted-content, got %s", s.Message)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serverResponse{Message: "stored-id"})
	}))
	defer ts.Close()

	s := domain.Secret{Message: "encrypted-content"}
	got, err := Store(ts.URL, s)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}
	if got != "stored-id" {
		t.Errorf("Expected stored-id, got %s", got)
	}
}

func TestFetchError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serverResponse{Message: "error message"})
	}))
	defer ts.Close()

	_, err := Fetch(ts.URL, "any-id")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !contains(err.Error(), "error message") {
		t.Errorf("Expected error to contain 'error message', got %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))
}

func TestServerErrorUnwrap(t *testing.T) {
	inner := http.ErrHandlerTimeout
	err := &ServerError{err: inner}
	if err.Unwrap() != inner {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), inner)
	}
}
