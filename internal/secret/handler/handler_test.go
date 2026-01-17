package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Khovanskiy5/yopass/internal/config"
	"github.com/Khovanskiy5/yopass/internal/secret/domain"
	"github.com/gorilla/mux"
	"go.uber.org/zap/zaptest"
)

type mockService struct {
	createKey string
	createErr error
	getSecret domain.Secret
	getErr    error
	status    bool
	statusErr error
	deleteRes bool
	deleteErr error
}

func (m *mockService) CreateSecret(secret domain.Secret) (string, error) {
	return m.createKey, m.createErr
}
func (m *mockService) GetSecret(key string) (domain.Secret, error) {
	return m.getSecret, m.getErr
}
func (m *mockService) GetSecretStatus(key string) (bool, error) {
	return m.status, m.statusErr
}
func (m *mockService) DeleteSecret(key string) (bool, error) {
	return m.deleteRes, m.deleteErr
}

func TestSecretHandler_CreateSecret(t *testing.T) {
	svc := &mockService{createKey: "test-key"}
	h := NewSecretHandler(svc, zaptest.NewLogger(t))

	secret := domain.Secret{Message: "-----BEGIN PGP MESSAGE-----\n...\n-----END PGP MESSAGE-----", Expiration: 3600}
	body, _ := json.Marshal(secret)
	req := httptest.NewRequest(http.MethodPost, "/secret", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateSecret(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["message"] != "test-key" {
		t.Errorf("expected message test-key, got %s", resp["message"])
	}
}

func TestSecretHandler_GetSecret(t *testing.T) {
	svc := &mockService{getSecret: domain.Secret{Message: "encrypted"}}
	h := NewSecretHandler(svc, zaptest.NewLogger(t))

	req := httptest.NewRequest(http.MethodGet, "/secret/test-key", nil)
	req = mux.SetURLVars(req, map[string]string{"key": "test-key"})
	w := httptest.NewRecorder()

	h.GetSecret(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp domain.Secret
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Message != "encrypted" {
		t.Errorf("expected message encrypted, got %s", resp.Message)
	}
}

func TestConfigHandler_GetConfig(t *testing.T) {
	cfg := &config.Config{DisableUpload: true}
	h := NewConfigHandler(cfg, zaptest.NewLogger(t))

	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	w := httptest.NewRecorder()

	h.GetConfig(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]bool
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["DISABLE_UPLOAD"] != true {
		t.Errorf("expected DISABLE_UPLOAD true, got %v", resp["DISABLE_UPLOAD"])
	}
}
