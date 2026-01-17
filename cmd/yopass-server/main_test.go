package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Khovanskiy5/yopass/internal/config"
	"github.com/Khovanskiy5/yopass/internal/utils"
)

func TestRun(t *testing.T) {
	cfg := &config.Config{
		Address:            "127.0.0.1",
		Port:               13371, // Use a specific port for testing
		Database:           "memcached",
		Memcached:          "localhost:11211",
		MaxLength:          1000,
		AllowedExpirations: []int{3600},
		AssetPath:          "../../public", // dummy path
	}

	logger := utils.NewLogger()
	registry := utils.NewRegistry()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Run in a goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- run(ctx, cfg, logger, registry)
	}()

	// Wait a bit for server to start
	time.Sleep(500 * time.Millisecond)

	// Check if we can ping the server
	resp, err := http.Get("http://127.0.0.1:13371/config")
	if err != nil {
		t.Fatalf("Failed to ping server: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// Cancel context to stop server
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("Server did not shut down in time")
	}
}
