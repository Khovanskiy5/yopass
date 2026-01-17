package domain

import (
	"encoding/json"
	"testing"
)

func TestSecretToJSON(t *testing.T) {
	s := Secret{
		Expiration: 3600,
		Message:    "test message",
		OneTime:    true,
	}
	
	got, err := s.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}
	
	var decoded Secret
	if err := json.Unmarshal(got, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	
	if decoded.Expiration != s.Expiration || decoded.Message != s.Message || decoded.OneTime != s.OneTime {
		t.Errorf("Decoded secret %v does not match original %v", decoded, s)
	}
}
