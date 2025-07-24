package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/config"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/models"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Статическая проверка что SubscriptionRepository реализует repository.SubscriptionRepository
var _ repository.SubscriptionRepository = (*SubscriptionRepository)(nil)

// SubscriptionRepository postgres реализация repository.SubscriptionRepository
type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(ctx context.Context, cfg config.DBConfig) *SubscriptionRepository {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name,
	)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		slog.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}

	slog.Info("Connected to postgres database")
	return &SubscriptionRepository{pool: pool}
}

func (r *SubscriptionRepository) Close() {
	r.pool.Close()
	slog.Info("Disconnected from postgres database")
}

func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, sub *models.Subscription) error {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
                  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.pool.QueryRow(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&sub.ID)
	if err != nil {
		slog.Error("Failed to create subscription", "error", err, "subscription", sub)
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	slog.Debug("Subscription created", "id", sub.ID.String(), "user_id", sub.UserID)
	return nil
}

func (r *SubscriptionRepository) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	query := `SELECT * FROM subscriptions WHERE id = $1`
	sub := &models.Subscription{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&sub)
	if err != nil {
		slog.Error("Failed to get subscription", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	slog.Debug("Subscription found", "subscription", sub)
	return sub, nil
}

func (r *SubscriptionRepository) GetAllSubscriptions(ctx context.Context) ([]models.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (r *SubscriptionRepository) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *SubscriptionRepository) ReplaceSubscription(ctx context.Context, sub *models.ReplaceSubscriptionRequest) error {
	//TODO implement me
	panic("implement me")
}

func (r *SubscriptionRepository) UpdateSubscription(ctx context.Context, sub *models.UpdateSubscriptionRequest) error {
	//TODO implement me
	panic("implement me")
}

func (r *SubscriptionRepository) GetTotalCostWithFilters(ctx context.Context, request *models.TotalCostRequest) (int, error) {
	//TODO implement me
	panic("implement me")
}
