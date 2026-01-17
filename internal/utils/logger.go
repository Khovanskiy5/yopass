package utils

import (
	"log"

	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	loggerCfg := zap.NewProductionConfig()
	logger, err := loggerCfg.Build()
	if err != nil {
		log.Fatalf("unable to build logger: %v", err)
	}
	zap.ReplaceGlobals(logger)
	return logger
}
