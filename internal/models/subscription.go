package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

// Subscription представляет подписку в базе данных и в ответах
type Subscription struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	ServiceName string       `json:"service_name" db:"service_name"`
	Price       int          `json:"price" db:"price"`
	UserID      uuid.UUID    `json:"user_id" db:"user_id"`
	StartDate   time.Time    `json:"start_date" db:"start_date"`
	EndDate     sql.NullTime `json:"end_date" db:"end_date"`
}

// CreateSubscriptionRequest DTO для запроса создания подписки
type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int       `json:"price" validate:"required,min=0"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	StartDate   string    `json:"start_date" validate:"required"`
	EndDate     *string   `json:"end_date,omitempty"`
}

// ReplaceSubscriptionRequest DTO для запроса замены подписки (PUT метод)
// Предполагается что пользователя подписки изменить нельзя
type ReplaceSubscriptionRequest struct {
	ServiceName string `json:"service_name" validate:"required"`
	Price       int    `json:"price" validate:"required,min=0"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date,omitempty"`
}

// UpdateSubscriptionRequest DTO для запроса изменения подписки (PATCH метод)
// Предполагается что пользователя подписки изменить нельзя
type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

// TotalCostRequest DTO для запроса суммарной стоимости
type TotalCostRequest struct {
	UserID      *uuid.UUID `query:"user_id"`
	ServiceName *string    `query:"service_name"`
	StartDate   *string    `query:"start_date"`
	EndDate     *string    `query:"end_date"`
}

// TotalCostResponse DTO для ответа суммарной стоимости
type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}
