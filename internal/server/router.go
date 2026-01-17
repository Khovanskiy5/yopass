package server

import (
	"net/http"

	"github.com/Khovanskiy5/yopass/internal/config"
	"github.com/Khovanskiy5/yopass/internal/constants"
	"github.com/Khovanskiy5/yopass/internal/middleware"
	"github.com/Khovanskiy5/yopass/internal/secret/handler"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

func NewRouter(
	cfg *config.Config,
	secretHandler *handler.SecretHandler,
	configHandler *handler.ConfigHandler,
	registry *prometheus.Registry,
) http.Handler {
	mx := mux.NewRouter()
	mx.Use(middleware.Metrics(registry))
	mx.Use(middleware.CORS(cfg.CORSAllowOrigin))

	// Secret routes
	mx.HandleFunc("/secret", secretHandler.CreateSecret).Methods(http.MethodPost)
	mx.HandleFunc("/secret", secretHandler.OptionsSecret).Methods(http.MethodOptions)
	if cfg.PrefetchSecret {
		mx.HandleFunc("/secret/"+constants.KeyParameter+"/status", secretHandler.GetSecretStatus).Methods(http.MethodGet)
	}
	mx.HandleFunc("/secret/"+constants.KeyParameter, secretHandler.GetSecret).Methods(http.MethodGet)
	mx.HandleFunc("/secret/"+constants.KeyParameter, secretHandler.DeleteSecret).Methods(http.MethodDelete)

	// Config routes
	mx.HandleFunc("/config", configHandler.GetConfig).Methods(http.MethodGet)
	mx.HandleFunc("/config", configHandler.OptionsConfig).Methods(http.MethodOptions)

	// File routes (if enabled)
	if !cfg.DisableUpload {
		mx.HandleFunc("/file", secretHandler.CreateSecret).Methods(http.MethodPost)
		mx.HandleFunc("/file", secretHandler.OptionsSecret).Methods(http.MethodOptions)
		if cfg.PrefetchSecret {
			mx.HandleFunc("/file/"+constants.KeyParameter+"/status", secretHandler.GetSecretStatus).Methods(http.MethodGet)
		}
		mx.HandleFunc("/file/"+constants.KeyParameter, secretHandler.GetSecret).Methods(http.MethodGet)
		mx.HandleFunc("/file/"+constants.KeyParameter, secretHandler.DeleteSecret).Methods(http.MethodDelete)
	}

	// Static files
	mx.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.AssetPath)))

	// Security headers
	return middleware.SecurityHeaders(mx)
}
