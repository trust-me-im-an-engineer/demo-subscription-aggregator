package postgres

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
