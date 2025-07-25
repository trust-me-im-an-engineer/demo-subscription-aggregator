package models

import (
	"github.com/google/uuid"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/pkg/monthyear"
)

type CreateSubscriptionRequest struct {
	ServiceName string               `json:"service_name" validate:"required"`
	Price       int                  `json:"price" validate:"required,min=0"`
	UserID      uuid.UUID            `json:"user_id" validate:"required"`
	StartDate   monthyear.MonthYear  `json:"start_date" validate:"required"`
	EndDate     *monthyear.MonthYear `json:"end_date,omitempty"`
}

// Предполагается что пользователя и дату начала изменить нельзя
type UpdateSubscriptionRequest struct {
	ServiceName *string              `json:"service_name,omitempty"`
	Price       *int                 `json:"price,omitempty"`
	EndDate     *monthyear.MonthYear `json:"end_date,omitempty"`
}

type TotalCostRequest struct {
	UserID      *uuid.UUID           `query:"user_id"`
	ServiceName *string              `query:"service_name"`
	StartDate   *monthyear.MonthYear `query:"start_date"`
	EndDate     *monthyear.MonthYear `query:"end_date"`
}

type SubscriptionResponse struct {
	ID          uuid.UUID            `json:"id"`
	UserID      uuid.UUID            `json:"user_id"`
	ServiceName string               `json:"service_name"`
	Price       int                  `json:"price"`
	StartDate   monthyear.MonthYear  `json:"start_date"`
	EndDate     *monthyear.MonthYear `json:"end_date"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}
