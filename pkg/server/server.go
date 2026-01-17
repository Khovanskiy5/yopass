package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/Khovanskiy5/yopass/pkg/yopass"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Server struct holding service and settings.
// This should be created with server.New
type Server struct {
	Service             yopass.Service
	MaxLength           int
	Registry            *prometheus.Registry
	ForceOneTimeSecrets bool
	AssetPath           string
	Logger              *zap.Logger
	TrustedProxies      []string
}

func (s *Server) sendError(w http.ResponseWriter, msg string, code int) {
	s.Logger.Debug("Sending error response", zap.String("message", msg), zap.Int("code", code))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

// createSecret creates secret
func (s *Server) createSecret(w http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var secret yopass.Secret
	if err := decoder.Decode(&secret); err != nil {
		s.sendError(w, "Unable to parse json", http.StatusBadRequest)
		return
	}

	key, err := s.Service.CreateSecret(secret)
	if err != nil {
		code := http.StatusBadRequest
		if err.Error() == "Failed to store secret in database" {
			code = http.StatusInternalServerError
		}
		s.sendError(w, err.Error(), code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": key}); err != nil {
		s.Logger.Error("Failed to write response", zap.Error(err))
	}
}

// getSecret from database
func (s *Server) getSecret(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Cache-Control", "private, no-cache")

	secretKey := mux.Vars(request)["key"]
	secret, err := s.Service.GetSecret(secretKey)
	if err != nil {
		s.sendError(w, "Secret not found", http.StatusNotFound)
		return
	}

	data, err := secret.ToJSON()
	if err != nil {
		s.Logger.Error("Failed to encode request", zap.Error(err))
		s.sendError(w, "Failed to encode secret", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		s.Logger.Error("Failed to write response", zap.Error(err))
	}
}

// getSecretStatus returns minimal status for a secret without returning the secret content
func (s *Server) getSecretStatus(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Cache-Control", "private, no-cache")
	w.Header().Set("Content-Type", "application/json")

	secretKey := mux.Vars(request)["key"]
	oneTime, err := s.Service.GetSecretStatus(secretKey)
	if err != nil {
		s.sendError(w, "Secret not found", http.StatusNotFound)
		return
	}

	resp := map[string]bool{"oneTime": oneTime}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		s.Logger.Error("Failed to write status response", zap.Error(err))
	}
}

// deleteSecret from database
func (s *Server) deleteSecret(w http.ResponseWriter, request *http.Request) {
	deleted, err := s.Service.DeleteSecret(mux.Vars(request)["key"])
	if err != nil {
		s.sendError(w, "Failed to delete secret", http.StatusInternalServerError)
		return
	}

	if !deleted {
		s.sendError(w, "Secret not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// optionsSecret handle the Options http method by returning the correct CORS headers
func (s *Server) optionsSecret(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
}

func (s *Server) configHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Content-Type", "application/json")

	config := map[string]interface{}{
		"DISABLE_UPLOAD":        viper.GetBool("disable-upload"),
		"PREFETCH_SECRET":       viper.GetBool("prefetch-secret"),
		"DISABLE_FEATURES":      viper.GetBool("disable-features"),
		"NO_LANGUAGE_SWITCHER":  viper.GetBool("no-language-switcher"),
		"FORCE_ONETIME_SECRETS": viper.GetBool("force-onetime-secrets"),
	}

	// Add optional string URLs only if they are provided
	if privacyURL := viper.GetString("privacy-notice-url"); privacyURL != "" {
		config["PRIVACY_NOTICE_URL"] = privacyURL
	}
	if imprintURL := viper.GetString("imprint-url"); imprintURL != "" {
		config["IMPRINT_URL"] = imprintURL
	}

	if err := json.NewEncoder(w).Encode(config); err != nil {
		s.Logger.Error("Failed to encode config response", zap.Error(err))
	}
}

// HTTPHandler containing all routes
func (s *Server) HTTPHandler() http.Handler {
	mx := mux.NewRouter()
	mx.Use(newMetricsMiddleware(s.Registry))
	mx.Use(corsMiddleware)

	mx.HandleFunc("/secret", s.createSecret).Methods(http.MethodPost)
	mx.HandleFunc("/secret", s.optionsSecret).Methods(http.MethodOptions)
	if viper.GetBool("prefetch-secret") {
		mx.HandleFunc("/secret/"+keyParameter+"/status", s.getSecretStatus).Methods(http.MethodGet)
	}
	mx.HandleFunc("/secret/"+keyParameter, s.getSecret).Methods(http.MethodGet)
	mx.HandleFunc("/secret/"+keyParameter, s.deleteSecret).Methods(http.MethodDelete)

	mx.HandleFunc("/config", s.configHandler).Methods(http.MethodGet)
	mx.HandleFunc("/config", s.optionsSecret).Methods(http.MethodOptions)

	if !viper.GetBool("disable-upload") {
		mx.HandleFunc("/file", s.createSecret).Methods(http.MethodPost)
		mx.HandleFunc("/file", s.optionsSecret).Methods(http.MethodOptions)
		if viper.GetBool("prefetch-secret") {
			mx.HandleFunc("/file/"+keyParameter+"/status", s.getSecretStatus).Methods(http.MethodGet)
		}
		mx.HandleFunc("/file/"+keyParameter, s.getSecret).Methods(http.MethodGet)
		mx.HandleFunc("/file/"+keyParameter, s.deleteSecret).Methods(http.MethodDelete)
	}

	mx.PathPrefix("/").Handler(http.FileServer(http.Dir(s.AssetPath)))
	return handlers.CustomLoggingHandler(nil, SecurityHeadersHandler(mx), s.httpLogFormatter())
}

const keyParameter = "{key:(?:[0-9a-f]{8}-(?:[0-9a-f]{4}-){3}[0-9a-f]{12})}"
