package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Khovanskiy5/yopass/internal/config"
	"github.com/Khovanskiy5/yopass/internal/repository"
	"github.com/Khovanskiy5/yopass/internal/secret/handler"
	"github.com/Khovanskiy5/yopass/internal/secret/service"
	"github.com/Khovanskiy5/yopass/internal/server"
	"github.com/Khovanskiy5/yopass/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// 2. Initialize infrastructure
	logger := utils.NewLogger()
	registry := utils.NewRegistry()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg, logger, registry); err != nil {
		logger.Fatal("server error", zap.Error(err))
	}
}

func run(ctx context.Context, cfg *config.Config, logger *zap.Logger, registry *prometheus.Registry) error {
	// 3. Setup repository
	repo, err := repository.NewRepository(cfg, logger)
	if err != nil {
		return err
	}

	// 4. Setup business logic
	var allowedExpirationsI32 []int32
	for _, e := range cfg.AllowedExpirations {
		allowedExpirationsI32 = append(allowedExpirationsI32, int32(e))
	}
	secretService := service.NewSecretService(
		repo,
		cfg.MaxLength,
		cfg.ForceOneTimeSecrets,
		allowedExpirationsI32,
	)

	// 5. Setup handlers
	secretHandler := handler.NewSecretHandler(secretService, logger)
	configHandler := handler.NewConfigHandler(cfg, logger)

	// 6. Setup router
	router := server.NewRouter(cfg, secretHandler, configHandler, registry)

	// 7. Start servers
	srvManager := server.NewServer(cfg, logger, registry)
	apiSrv := srvManager.Start(router)
	metricsSrv := srvManager.StartMetrics()

	// 8. Wait for termination signal or context cancellation
	<-ctx.Done()
	logger.Info("Shutting down servers")

	// 9. Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	srvManager.Shutdown(shutdownCtx, apiSrv, metricsSrv)
	logger.Info("Server gracefully stopped")
	return nil
}
