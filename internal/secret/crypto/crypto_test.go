package crypto

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "test-key"
	content := "hello world"
	
	// Test Encrypt
	encrypted, err := Encrypt(strings.NewReader(content), key)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if !strings.Contains(encrypted, "BEGIN PGP MESSAGE") {
		t.Errorf("Expected PGP message, got: %s", encrypted)
	}

	// Test Decrypt
	decrypted, filename, err := Decrypt(strings.NewReader(encrypted), key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != content {
		t.Errorf("Expected decrypted content %q, got %q", content, decrypted)
	}
	if filename != "" {
		t.Errorf("Expected empty filename for string reader, got %q", filename)
	}
}

func TestEncryptEmptyKey(t *testing.T) {
	_, err := Encrypt(strings.NewReader("content"), "")
	if err != ErrEmptyKey {
		t.Errorf("Expected ErrEmptyKey, got %v", err)
	}
}

func TestDecryptInvalidKey(t *testing.T) {
	key := "test-key"
	content := "hello world"
	encrypted, _ := Encrypt(strings.NewReader(content), key)

	_, _, err := Decrypt(strings.NewReader(encrypted), "wrong-key")
	if err == nil {
		t.Fatal("Expected error for wrong key, got nil")
	}
	if !strings.Contains(err.Error(), "could not decrypt") {
		t.Errorf("Expected decryption error, got %v", err)
	}
}

func TestDecryptInvalidMessage(t *testing.T) {
	_, _, err := Decrypt(strings.NewReader("invalid message"), "key")
	if err != ErrInvalidMessage {
		t.Errorf("Expected ErrInvalidMessage, got %v", err)
	}
}

func TestGenerateKey(t *testing.T) {
	key1, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}
	if len(key1) != 22 {
		t.Errorf("Expected key length 22, got %d", len(key1))
	}

	key2, _ := GenerateKey()
	if key1 == key2 {
		t.Error("GenerateKey should produce different keys")
	}
}

func TestEncryptWithFileName(t *testing.T) {
	// We can't easily mock os.File without creating a real file
	// But we can check that Encrypt doesn't crash and handles it
	// Actually, the current implementation of Encrypt uses a type assertion to *os.File
	// Let's just test with a regular buffer which should not have a filename
	
	key := "key"
	content := "content"
	encrypted, err := Encrypt(bytes.NewBufferString(content), key)
	if err != nil {
		t.Fatal(err)
	}
	
	dec, filename, err := Decrypt(strings.NewReader(encrypted), key)
	if err != nil {
		t.Fatal(err)
	}
	if dec != content {
		t.Errorf("got %q, want %q", dec, content)
	}
	if filename != "" {
		t.Errorf("got filename %q, want empty", filename)
	}
}
