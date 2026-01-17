package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Khovanskiy5/yopass/internal/secret/domain"
	"github.com/Khovanskiy5/yopass/internal/secret/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type SecretHandler struct {
	service service.SecretService
	logger  *zap.Logger
}

func NewSecretHandler(service service.SecretService, logger *zap.Logger) *SecretHandler {
	return &SecretHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SecretHandler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	var secret domain.Secret
	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		h.sendError(w, "Unable to parse json", http.StatusBadRequest)
		return
	}

	key, err := h.service.CreateSecret(secret)
	if err != nil {
		code := http.StatusBadRequest
		if err.Error() == "failed to store secret in database" {
			code = http.StatusInternalServerError
		}
		h.sendError(w, err.Error(), code)
		return
	}

	h.sendJSON(w, map[string]string{"message": key}, http.StatusOK)
}

func (h *SecretHandler) GetSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, no-cache")
	key := mux.Vars(r)["key"]

	secret, err := h.service.GetSecret(key)
	if err != nil {
		h.sendError(w, "Secret not found", http.StatusNotFound)
		return
	}

	data, err := secret.ToJSON()
	if err != nil {
		h.logger.Error("Failed to encode secret", zap.Error(err))
		h.sendError(w, "Failed to encode secret", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		h.logger.Error("Failed to write response", zap.Error(err))
	}
}

func (h *SecretHandler) GetSecretStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, no-cache")
	key := mux.Vars(r)["key"]

	oneTime, err := h.service.GetSecretStatus(key)
	if err != nil {
		h.sendError(w, "Secret not found", http.StatusNotFound)
		return
	}

	h.sendJSON(w, map[string]bool{"oneTime": oneTime}, http.StatusOK)
}

func (h *SecretHandler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	deleted, err := h.service.DeleteSecret(key)
	if err != nil {
		h.sendError(w, "Failed to delete secret", http.StatusInternalServerError)
		return
	}

	if !deleted {
		h.sendError(w, "Secret not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SecretHandler) OptionsSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.WriteHeader(http.StatusOK)
}

func (h *SecretHandler) sendError(w http.ResponseWriter, msg string, code int) {
	h.logger.Debug("Sending error response", zap.String("message", msg), zap.Int("code", code))
	h.sendJSON(w, map[string]string{"message": msg}, code)
}

func (h *SecretHandler) sendJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
	}
}
