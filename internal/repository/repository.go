package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/models"
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

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, sub Subscription) (uuid.UUID, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (Subscription, error)
	GetAllSubscriptions(ctx context.Context) ([]Subscription, error)
	ReplaceSubscription(ctx context.Context, id uuid.UUID, sub models.ReplaceSubscriptionRequest) error
	UpdateSubscription(ctx context.Context, id uuid.UUID, sub models.UpdateSubscriptionRequest) error
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetTotalCostWithFilters(ctx context.Context, request models.TotalCostRequest) (int, error)
}
