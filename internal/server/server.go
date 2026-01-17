package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/Khovanskiy5/yopass/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	cfg      *config.Config
	logger   *zap.Logger
	registry *prometheus.Registry
}

func NewServer(cfg *config.Config, logger *zap.Logger, registry *prometheus.Registry) *Server {
	return &Server{
		cfg:      cfg,
		logger:   logger,
		registry: registry,
	}
}

func (s *Server) Start(handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:      fmt.Sprintf("%s:%d", s.cfg.Address, s.cfg.Port),
		Handler:   handler,
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}

	go func() {
		s.logger.Info("Starting yopass server", zap.String("address", srv.Addr))
		var err error
		if s.cfg.TLSCert != "" && s.cfg.TLSKey != "" {
			err = srv.ListenAndServeTLS(s.cfg.TLSCert, s.cfg.TLSKey)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("yopass stopped unexpectedly", zap.Error(err))
		}
	}()

	return srv
}

func (s *Server) StartMetrics() *http.Server {
	if s.cfg.MetricsPort <= 0 {
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.cfg.Address, s.cfg.MetricsPort),
		Handler: mux,
	}

	go func() {
		s.logger.Info("Starting metrics server", zap.String("address", srv.Addr))
		var err error
		if s.cfg.TLSCert != "" && s.cfg.TLSKey != "" {
			err = srv.ListenAndServeTLS(s.cfg.TLSCert, s.cfg.TLSKey)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("metrics server stopped unexpectedly", zap.Error(err))
		}
	}()

	return srv
}

func (s *Server) Shutdown(ctx context.Context, servers ...*http.Server) {
	for _, srv := range servers {
		if srv != nil {
			s.logger.Info("Shutting down server", zap.String("address", srv.Addr))
			if err := srv.Shutdown(ctx); err != nil {
				s.logger.Error("Error shutting down server", zap.Error(err))
			}
		}
	}
}
