package repository

import (
	"context"
	"github.com/google/uuid"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/models"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, sub *models.Subscription) error
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	GetAllSubscriptions(ctx context.Context) ([]models.Subscription, error)
	ReplaceSubscription(ctx context.Context, sub *models.ReplaceSubscriptionRequest) error
	UpdateSubscription(ctx context.Context, sub *models.UpdateSubscriptionRequest) error
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetTotalCostWithFilters(ctx context.Context, request *models.TotalCostRequest) (int, error)
}
