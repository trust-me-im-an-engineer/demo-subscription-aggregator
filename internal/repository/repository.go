package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

// Subscription представляет подписку в базе данных
type Subscription struct {
	ID          uuid.UUID    `db:"id"`
	ServiceName string       `db:"service_name"`
	Price       int          `db:"price"`
	UserID      uuid.UUID    `db:"user_id"`
	StartDate   time.Time    `db:"start_date"`
	EndDate     sql.NullTime `db:"end_date"`
}

type SubscriptionFilter struct {
	ServiceName *string
	UserID      *uuid.UUID
	StartDate   *time.Time
	EndDate     *time.Time
}

type ErrSubscriptionNotFound struct{}

type ErrSubscriptionAlreadyExists struct{}

func (e *ErrSubscriptionNotFound) Error() string {
	return "subscription not found"
}

func (e *ErrSubscriptionAlreadyExists) Error() string {
	return "subscription already exists"
}

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, sub Subscription) (uuid.UUID, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (Subscription, error)
	GetAllSubscriptions(ctx context.Context) ([]Subscription, error)
	UpdateSubscription(ctx context.Context, sub Subscription) (Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetTotalCostWithFilters(ctx context.Context, filter SubscriptionFilter) (int, error)
}
