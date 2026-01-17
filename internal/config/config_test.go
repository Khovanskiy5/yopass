package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func TestLoad(t *testing.T) {
	// Reset pflag and viper to avoid interference with other tests or global state
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	viper.Reset()

	// Set environment variables to test viper integration
	os.Setenv("YOPASS_PORT", "9999")
	defer os.Unsetenv("YOPASS_PORT")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Check default values and env override
	if cfg.Port != 9999 {
		t.Errorf("Expected Port 9999, got %d", cfg.Port)
	}
	if cfg.Database != "memcached" {
		t.Errorf("Expected default Database memcached, got %s", cfg.Database)
	}
	if cfg.MaxLength != 5242880 {
		t.Errorf("Expected default MaxLength 5242880, got %d", cfg.MaxLength)
	}
}

func TestConfigDefaults(t *testing.T) {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	viper.Reset()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Address != "" {
		t.Errorf("Expected empty default address, got %s", cfg.Address)
	}
	if cfg.Redis != "redis://localhost:6379/0" {
		t.Errorf("Expected default redis URL, got %s", cfg.Redis)
	}
	if cfg.CORSAllowOrigin != "*" {
		t.Errorf("Expected default CORS allow origin *, got %s", cfg.CORSAllowOrigin)
	}
	if !cfg.PrefetchSecret {
		t.Error("Expected PrefetchSecret to be true by default")
	}
}
