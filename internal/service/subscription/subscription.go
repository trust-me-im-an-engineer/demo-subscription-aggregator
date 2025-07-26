package subscription

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/models"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/repository"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/service"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/pkg/monthyear"
)

type Service struct {
	repo repository.SubscriptionRepository
}

func NewService(repo repository.SubscriptionRepository) Service {
	return Service{repo: repo}
}

func (s Service) CreateSubscription(ctx context.Context, req models.CreateSubscriptionRequest) (models.SubscriptionResponse, error) {
	sub := repository.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   time.Time(req.StartDate),
	}
	if req.EndDate != nil {
		endDate := time.Time(*req.EndDate)
		if endDate.Before(sub.StartDate) {
			return models.SubscriptionResponse{}, &service.ErrInvalidDateRange{}
		}

		sub.EndDate = sql.NullTime{
			Time:  endDate,
			Valid: true,
		}
	}

	id, err := s.repo.CreateSubscription(ctx, sub)
	if err != nil {
		return models.SubscriptionResponse{}, fmt.Errorf("repo failed to create subcsciption: %w", err)
	}

	return models.SubscriptionResponse{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}, nil
}

func (s Service) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (models.SubscriptionResponse, error) {
	sub, err := s.repo.GetSubscriptionByID(ctx, id)
	if err != nil {
		return models.SubscriptionResponse{}, fmt.Errorf("repo failed to get subcsciption by id: %w", err)
	}

	resp := models.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   monthyear.MonthYear(sub.StartDate),
	}
	if sub.EndDate.Valid {
		endDate := monthyear.MonthYear(sub.EndDate.Time)
		resp.EndDate = &endDate
	}

	return resp, nil
}

func (s Service) GetAllSubscriptions(ctx context.Context) ([]models.SubscriptionResponse, error) {
	subs, err := s.repo.GetAllSubscriptions(ctx)
	if err != nil {
		return []models.SubscriptionResponse{}, fmt.Errorf("repo failed to get all subcsciptions: %w", err)
	}
	responds := make([]models.SubscriptionResponse, len(subs))
	for i, sub := range subs {
		responds[i] = models.SubscriptionResponse{
			ID:          sub.ID,
			ServiceName: sub.ServiceName,
			Price:       sub.Price,
			UserID:      sub.UserID,
			StartDate:   monthyear.MonthYear(sub.StartDate),
		}
		if sub.EndDate.Valid {
			endDate := monthyear.MonthYear(sub.EndDate.Time)
			responds[i].EndDate = &endDate
		}
	}

	return responds, nil
}

func (s Service) UpdateSubscription(ctx context.Context, id uuid.UUID, req models.UpdateSubscriptionRequest) (models.SubscriptionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetTotalCost(ctx context.Context, filter models.TotalCostRequest) (models.TotalCostResponse, error) {
	//TODO implement me
	panic("implement me")
}
