package models

import (
	"github.com/google/uuid"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/pkg/monthyear"
)

// CreateSubscriptionRequest представляет запрос на создание новой подписки
type CreateSubscriptionRequest struct {
	ServiceName string               `json:"service_name" validate:"required" example:"Netflix" description:"Название сервиса"`
	Price       int                  `json:"price" validate:"required,min=0" example:"299" description:"Стоимость в рублях за месяц"`
	UserID      uuid.UUID            `json:"user_id" validate:"required" example:"123e4567-e89b-12d3-a456-426614174000" description:"ID пользователя"`
	StartDate   *monthyear.MonthYear `json:"start_date" validate:"required" example:"01-2024" description:"Дата начала в формате ММ-ГГГГ"`
	EndDate     *monthyear.MonthYear `json:"end_date,omitempty" example:"12-2024" description:"Дата окончания в формате ММ-ГГГГ (необязательно)"`
}

// UpdateSubscriptionRequest представляет запрос на обновление существующей подписки
// Примечание: ID пользователя и дата начала не могут быть изменены
type UpdateSubscriptionRequest struct {
	ServiceName *string              `json:"service_name,omitempty" example:"Netflix Premium" description:"Обновлённое название сервиса"`
	Price       *int                 `json:"price,omitempty" example:"599" description:"Обновлённая стоимость в рублях за месяц"`
	EndDate     *monthyear.MonthYear `json:"end_date,omitempty" example:"12-2024" description:"Обновлённая дата окончания в формате ММ-ГГГГ"`
}

// TotalCostRequest представляет параметры запроса для расчёта общей стоимости
type TotalCostRequest struct {
	UserID      *uuid.UUID           `json:"user_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000" description:"Фильтр по ID пользователя"`
	ServiceName *string              `json:"service_name,omitempty" example:"Netflix" description:"Фильтр по названию сервиса (частичное совпадение)"`
	StartDate   *monthyear.MonthYear `json:"start_date,omitempty" example:"01-2024" description:"Фильтр по дате начала (подписки, начинающиеся с этой даты)"`
	EndDate     *monthyear.MonthYear `json:"end_date,omitempty" example:"12-2024" description:"Фильтр по дате окончания (подписки, заканчивающиеся до этой даты)"`
}

// SubscriptionResponse представляет подписку в ответах API
type SubscriptionResponse struct {
	ID          uuid.UUID            `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" description:"ID подписки"`
	UserID      uuid.UUID            `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" description:"ID пользователя"`
	ServiceName string               `json:"service_name" example:"Netflix" description:"Название сервиса"`
	Price       int                  `json:"price" example:"299" description:"Стоимость в рублях за месяц"`
	StartDate   *monthyear.MonthYear `json:"start_date" example:"01-2024" description:"Дата начала в формате ММ-ГГГГ"`
	EndDate     *monthyear.MonthYear `json:"end_date,omitempty" example:"12-2024" description:"Дата окончания в формате ММ-ГГГГ (null если активна)"`
}

// TotalCostResponse представляет ответ с расчётом общей стоимости
type TotalCostResponse struct {
	TotalCost int `json:"total_cost" example:"1499" description:"Общая стоимость в рублях"`
}
