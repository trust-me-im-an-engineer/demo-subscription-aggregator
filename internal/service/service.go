package service

import (
	"context"
	"errors"
	"github.com/google/uuid"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/models"
)

var ErrInvalidDateRange = errors.New("end date cannot be before start date")

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req models.CreateSubscriptionRequest) (models.SubscriptionResponse, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (models.SubscriptionResponse, error)
	GetAllSubscriptions(ctx context.Context) ([]models.SubscriptionResponse, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, req models.UpdateSubscriptionRequest) (models.SubscriptionResponse, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetTotalCost(ctx context.Context, filter models.TotalCostRequest) (models.TotalCostResponse, error)
}
