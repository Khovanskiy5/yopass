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

	// 3. Setup repository
	repo, err := repository.NewRepository(cfg, logger)
	if err != nil {
		logger.Fatal("failed to setup repository", zap.Error(err))
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

	// 8. Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("Shutting down servers", zap.String("signal", sig.String()))

	// 9. Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	srvManager.Shutdown(ctx, apiSrv, metricsSrv)
	logger.Info("Server gracefully stopped")
}
