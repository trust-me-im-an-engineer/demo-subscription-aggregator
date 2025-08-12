package api

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/web/handler"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/service"
)

func NewRouter(s service.SubscriptionService) *http.ServeMux {
	h := handler.NewHandler(s)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /subscriptions", h.Create)
	mux.HandleFunc("GET /subscriptions", h.List)
	mux.HandleFunc("GET /subscriptions/{id}", h.GetByID)
	mux.HandleFunc("PATCH /subscriptions/{id}", h.Update)
	mux.HandleFunc("DELETE /subscriptions/{id}", h.Delete)
	mux.HandleFunc("GET /subscriptions/total-cost", h.GetTotalCost)

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	return mux
}
