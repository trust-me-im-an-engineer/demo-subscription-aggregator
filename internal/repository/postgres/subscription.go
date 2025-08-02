package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/config"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/repository"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // PostgreSQL driver for golang-migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // File source for golang-migrate
)

// Статическая проверка что SubscriptionRepository реализует repository.SubscriptionRepository
var _ repository.SubscriptionRepository = (*SubscriptionRepository)(nil)

// SubscriptionRepository postgres реализация repository.SubscriptionRepository
type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.DBConfig) (SubscriptionRepository, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
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

	migrationPath := "file://./migrations" // Path to your migration files (relative to where the app runs)
	m, err := migrate.New(migrationPath, dsn)
	if err != nil {
		pool.Close() // Close the pool if migration instance creation fails
		return SubscriptionRepository{}, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	slog.Info("Running database migrations...")
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		pool.Close() // Close the pool on migration error
		return SubscriptionRepository{}, fmt.Errorf("failed to run migrations: %w", err)
	}
	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("No new database migrations to apply.")
	} else {
		slog.Info("Database migrations applied successfully.")
	}

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
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			// Unique constrain failed
			if pgxError.Code == "23505" {
				return id, repository.ErrSubscriptionAlreadyExists
			}
		}
		return id, err
	}

	slog.Debug("subscription created", "id", id.String(), "user_id", sub.UserID)
	return id, nil
}

func (r *SubscriptionRepository) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (repository.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`
	sub := repository.Subscription{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.Subscription{}, repository.ErrSubscriptionNotFound
		}
		return repository.Subscription{}, err
	}

	slog.Debug("subscription found", "subscription", sub)
	return sub, nil
}

func (r *SubscriptionRepository) GetAllSubscriptions(ctx context.Context) ([]repository.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}

	subs := make([]repository.Subscription, 0)
	for rows.Next() {
		var sub repository.Subscription
		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan all subscriptions: %w", err)
	}

	slog.Debug("all subscriptions fetched", "total", len(subs))
	return subs, nil
}

func (r *SubscriptionRepository) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrSubscriptionNotFound
		}
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	slog.Debug("subscription deleted", "id", id)
	return nil
}

func (r *SubscriptionRepository) UpdateSubscription(ctx context.Context, id uuid.UUID, fields repository.SubscriptionUpdate) (repository.Subscription, error) {
	query := squirrel.Update("subscriptions").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, service_name, price, user_id, start_date, end_date").
		PlaceholderFormat(squirrel.Dollar)

	if fields.ServiceName != nil {
		query = query.Set("service_name", *fields.ServiceName)
	}
	if fields.Price != nil {
		query = query.Set("price", *fields.Price)
	}
	if fields.EndDate != nil {
		query = query.Set("end_date", *fields.EndDate)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return repository.Subscription{}, fmt.Errorf("failed to build update query: %w", err)
	}

	var updatedSub repository.Subscription
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&updatedSub.ID,
		&updatedSub.ServiceName,
		&updatedSub.Price,
		&updatedSub.UserID,
		&updatedSub.StartDate,
		&updatedSub.EndDate,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.Subscription{}, repository.ErrSubscriptionNotFound
		}
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			// Unique constraint failed
			if pgxError.Code == "23505" {
				return repository.Subscription{}, repository.ErrSubscriptionAlreadyExists
			}
		}
		return repository.Subscription{}, fmt.Errorf("failed to update subscription: %w", err)
	}

	slog.Debug("subscription updated", "subscription", updatedSub)
	return updatedSub, nil
}

func (r *SubscriptionRepository) GetTotalCostWithFilters(ctx context.Context, filter repository.SubscriptionFilter) (int, error) {
	query := squirrel.Select("COALESCE(SUM(price), 0)").
		From("subscriptions").
		PlaceholderFormat(squirrel.Dollar)

	if filter.UserID != nil {
		query = query.Where(squirrel.Eq{"user_id": *filter.UserID})
	}
	if filter.ServiceName != nil {
		query = query.Where(squirrel.ILike{"service_name": "%" + *filter.ServiceName + "%"})
	}
	if filter.StartDate != nil {
		query = query.Where(squirrel.GtOrEq{"start_date": *filter.StartDate})
	}
	if filter.EndDate != nil {
		query = query.Where(squirrel.LtOrEq{"end_date": *filter.EndDate})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	var totalCost int
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	slog.Debug("total cost with filters calculated", "total_cost", totalCost, "filter", filter)
	return totalCost, nil
}
