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
	sub, err := s.repo.GetSubscriptionByID(ctx, id)
	if err != nil {
		return models.SubscriptionResponse{}, err
	}

	if req.ServiceName != nil {
		sub.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		sub.Price = *req.Price
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

	updatedSub, err := s.repo.UpdateSubscription(ctx, sub)
	if err != nil {
		return models.SubscriptionResponse{}, fmt.Errorf("repo failed to update subcsciption: %w", err)
	}

	resp := models.SubscriptionResponse{
		ID:          updatedSub.ID,
		ServiceName: updatedSub.ServiceName,
		Price:       updatedSub.Price,
		UserID:      updatedSub.UserID,
		StartDate:   monthyear.MonthYear(updatedSub.StartDate),
	}
	if updatedSub.EndDate.Valid {
		endDate := monthyear.MonthYear(updatedSub.EndDate.Time)
		resp.EndDate = &endDate
	}

	return resp, nil
}

func (s Service) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteSubscription(ctx, id)
}

func (s Service) GetTotalCost(ctx context.Context, req models.TotalCostRequest) (models.TotalCostResponse, error) {
	filter := repository.SubscriptionFilter{}

	if req.UserID != nil {
		filter.UserID = req.UserID
	}
	if req.ServiceName != nil {
		filter.ServiceName = req.ServiceName
	}
	if req.StartDate != nil {
		startDate := time.Time(*req.StartDate)
		filter.StartDate = &startDate
	}
	if req.EndDate != nil {
		endDate := time.Time(*req.EndDate)
		filter.EndDate = &endDate
	}

	totalCost, err := s.repo.GetTotalCostWithFilters(ctx, filter)
	if err != nil {
		return models.TotalCostResponse{}, fmt.Errorf("repo failed to get total cost: %w", err)
	}

	return models.TotalCostResponse{TotalCost: totalCost}, nil
}
