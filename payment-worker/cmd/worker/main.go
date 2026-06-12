package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/company/payment-worker/internal/client"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "payment-worker")
	coreClient := client.NewCoreClient(env("CORE_SERVICE_URL", "http://localhost:9001"))
	_ = coreClient

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info("worker started", "core_service_url", env("CORE_SERVICE_URL", "http://localhost:9001"))

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("worker stopped")
			return
		case <-ticker.C:
			logger.Info("scheduled task placeholder", "task", "query-order-compensation")
		}
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.LUTC)
}

