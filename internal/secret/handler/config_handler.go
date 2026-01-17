package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Khovanskiy5/yopass/internal/config"
	"go.uber.org/zap"
)

type ConfigHandler struct {
	cfg    *config.Config
	logger *zap.Logger
}

func NewConfigHandler(cfg *config.Config, logger *zap.Logger) *ConfigHandler {
	return &ConfigHandler{
		cfg:    cfg,
		logger: logger,
	}
}

func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Content-Type", "application/json")

	cfgMap := map[string]interface{}{
		"DISABLE_UPLOAD":        h.cfg.DisableUpload,
		"PREFETCH_SECRET":       h.cfg.PrefetchSecret,
		"DISABLE_FEATURES":      h.cfg.DisableFeatures,
		"NO_LANGUAGE_SWITCHER":  h.cfg.NoLanguageSwitcher,
		"FORCE_ONETIME_SECRETS": h.cfg.ForceOneTimeSecrets,
	}

	if h.cfg.PrivacyNoticeURL != "" {
		cfgMap["PRIVACY_NOTICE_URL"] = h.cfg.PrivacyNoticeURL
	}
	if h.cfg.ImprintURL != "" {
		cfgMap["IMPRINT_URL"] = h.cfg.ImprintURL
	}

	if err := json.NewEncoder(w).Encode(cfgMap); err != nil {
		h.logger.Error("Failed to encode config response", zap.Error(err))
	}
}

func (h *ConfigHandler) OptionsConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.WriteHeader(http.StatusOK)
}
