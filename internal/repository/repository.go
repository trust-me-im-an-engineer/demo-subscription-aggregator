package repository

import (
	"context"
	"database/sql"
	"errors"
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

// At least one field must be provided
type SubscriptionUpdate struct {
	ServiceName *string
	Price       *int
	EndDate     *time.Time
}

type SubscriptionFilter struct {
	ServiceName *string
	UserID      *uuid.UUID
	StartDate   *time.Time
	EndDate     *time.Time
}

var (
	ErrSubscriptionNotFound      = errors.New("subscription not found")
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, sub Subscription) (uuid.UUID, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (Subscription, error)
	GetAllSubscriptions(ctx context.Context) ([]Subscription, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, fields SubscriptionUpdate) (Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetTotalCostWithFilters(ctx context.Context, filter SubscriptionFilter) (int, error)
}
