package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/config"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Статическая проверка что SubscriptionRepository реализует repository.SubscriptionRepository
var _ repository.SubscriptionRepository = (*SubscriptionRepository)(nil)

// SubscriptionRepository postgres реализация repository.SubscriptionRepository
type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.DBConfig) (SubscriptionRepository, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name,
	)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return SubscriptionRepository{}, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return SubscriptionRepository{}, fmt.Errorf("failed to ping database: %w", err)
	}
	slog.Info("connected to postgres database")
	return SubscriptionRepository{pool: pool}, nil
}

func (r *SubscriptionRepository) Close() {
	r.pool.Close()
	slog.Info("disconnected from postgres database")
}

func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, sub repository.Subscription) (uuid.UUID, error) {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
                  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("failed to create subscription: %w", err)
	}

	slog.Debug("subscription created", "id", id.String(), "user_id", sub.UserID)
	return id, nil
}

func (r *SubscriptionRepository) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (repository.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`
	sub := repository.Subscription{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		return repository.Subscription{}, fmt.Errorf("failed to get subscription: %w", err)
	}

	slog.Debug("subscription found", "subscription", sub)
	return sub, nil
}

func (r *SubscriptionRepository) GetAllSubscriptions(ctx context.Context) ([]repository.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all subscriptions: %w", err)
	}

	subs := make([]repository.Subscription, 0)
	for rows.Next() {
		var sub repository.Subscription
		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get all subscriptions: %w", err)
		}
		subs = append(subs, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get all subscriptions: %w", err)
	}

	slog.Debug("all subscriptions fetched")
	return subs, nil
}

func (r *SubscriptionRepository) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	slog.Debug("subscription deleted", "id", id)
	return nil
}

func (r *SubscriptionRepository) UpdateSubscription(ctx context.Context, sub repository.Subscription) error {
	query := `UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6`
	_, err := r.pool.Exec(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	slog.Debug("subscription updated", "id", sub.ID)
	return nil
}

func (r *SubscriptionRepository) GetTotalCostWithFilters(ctx context.Context, filter repository.SubscriptionFilter) (int, error) {
	query := `SELECT SUM(price) FROM subscriptions WHERE 1 = 1`
	args := make([]any, 0, 1)
	argID := 1

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argID)
		args = append(args, *filter.UserID)
		argID++
	}
	if filter.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name ILIKE $%d", argID)
		args = append(args, "%"+*filter.ServiceName+"%")
		argID++
	}
	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND start_date >= $%d", argID)
		args = append(args, *filter.StartDate)
		argID++
	}
	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND end_date <= $%d", argID)
		args = append(args, *filter.EndDate)
	}

	var totalCost int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("failed to get total cost with filters: %w", err)
	}

	slog.Debug("total cost with filters calculated", "total_cost", totalCost, "filter", filter)
	return totalCost, nil
}
