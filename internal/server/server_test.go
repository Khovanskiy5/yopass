package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Khovanskiy5/yopass/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func TestNewServer(t *testing.T) {
	cfg := &config.Config{Address: "127.0.0.1", Port: 1337}
	logger := zap.NewNop()
	reg := prometheus.NewRegistry()
	s := NewServer(cfg, logger, reg)

	if s.cfg != cfg {
		t.Error("Config not set correctly")
	}
	if s.logger != logger {
		t.Error("Logger not set correctly")
	}
	if s.registry != reg {
		t.Error("Registry not set correctly")
	}
}

func TestServerStartShutdown(t *testing.T) {
	cfg := &config.Config{Address: "127.0.0.1", Port: 0} // Port 0 for random available port
	logger := zap.NewNop()
	reg := prometheus.NewRegistry()
	s := NewServer(cfg, logger, reg)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := s.Start(handler)
	
	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Shutdown(ctx, srv)
}

func TestMetricsServerStartShutdown(t *testing.T) {
	cfg := &config.Config{Address: "127.0.0.1", MetricsPort: 11337}
	logger := zap.NewNop()
	reg := prometheus.NewRegistry()
	s := NewServer(cfg, logger, reg)

	srv := s.StartMetrics()
	if srv == nil {
		t.Fatal("Expected metrics server to be created")
	}

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Shutdown(ctx, srv)
}

func TestMetricsServerDisabled(t *testing.T) {
	cfg := &config.Config{MetricsPort: -1}
	s := NewServer(cfg, zap.NewNop(), prometheus.NewRegistry())
	srv := s.StartMetrics()
	if srv != nil {
		t.Error("Metrics server should be disabled for MetricsPort <= 0")
	}
}
