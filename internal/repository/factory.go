package repository

import (
	"fmt"

	"github.com/Khovanskiy5/yopass/internal/config"
	"github.com/Khovanskiy5/yopass/internal/secret/domain"
	"go.uber.org/zap"
)

func NewRepository(cfg *config.Config, logger *zap.Logger) (domain.Repository, error) {
	switch cfg.Database {
	case "memcached":
		logger.Debug("Configuring Memcached", zap.String("address", cfg.Memcached))
		return NewMemcached(cfg.Memcached), nil
	case "redis":
		logger.Debug("Configuring Redis", zap.String("url", cfg.Redis))
		return NewRedis(cfg.Redis)
	default:
		return nil, fmt.Errorf("unsupported database: %s", cfg.Database)
	}
}
