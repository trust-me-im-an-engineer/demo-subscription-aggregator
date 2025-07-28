package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/config"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/repository/postgres"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/service/subscription"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/web"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/pkg/logger"

	// Import generated docs
	_ "github.com/trust-me-im-an-engineer/demo-subscription-agregator/docs"
)

// @title Subscription Aggregator API
// @version 1.0
// @description A service for managing user subscriptions
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

func main() {
	logger.Init()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := postgres.New(context.Background(), cfg.DB)
	if err != nil {
		slog.Error("failed to initialize repository", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	service := subscription.NewService(&db)
	router := web.NewRouter(service)

	server := &http.Server{
		Addr:    cfg.App.Port,
		Handler: router,
	}

	// Create a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		slog.Info("starting server", "port", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	sig := <-sigChan
	slog.Info("received signal, shutting down gracefully with "+cfg.App.ShutdownTimeout.String()+" seconds deadline", "signal", sig)

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.App.ShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server shutdown complete")
}
