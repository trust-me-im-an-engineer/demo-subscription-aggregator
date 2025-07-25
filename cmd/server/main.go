package main

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"time"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/config"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/repository"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/repository/postgres"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/pkg/logger"
)

func main() {
	logger.Init()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config: ", err)
	}

	db, err := postgres.New(context.Background(), cfg.DB)
	if err != nil {
		slog.Error("failed to connect to database")
		os.Exit(1)
	}
	defer db.Close()

	id, err := db.CreateSubscription(context.Background(), repository.Subscription{
		ServiceName: "test service",
		Price:       123,
		UserID:      uuid.UUID{},
		StartDate:   time.Now(),
		EndDate:     sql.NullTime{},
	})
	if err != nil {
		slog.Error("failed to create subscription", "error", err)
		os.Exit(1)
	}

	slog.Info("subscription created", "id", id.String())
}
