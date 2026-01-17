package repository

import (
	"testing"

	"github.com/Khovanskiy5/yopass/internal/config"
	"go.uber.org/zap/zaptest"
)

func TestNewRepository(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name     string
		db       string
		wantErr  bool
	}{
		{
			name:    "memcached",
			db:      "memcached",
			wantErr: false,
		},
		{
			name:    "redis",
			db:      "redis",
			wantErr: false,
		},
		{
			name:    "unsupported",
			db:      "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Database:  tt.db,
				Memcached: "localhost:11211",
				Redis:     "redis://localhost:6379/0",
			}
			repo, err := NewRepository(cfg, logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && repo == nil {
				t.Error("NewRepository() returned nil repo without error")
			}
		})
	}
}
