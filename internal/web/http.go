package web

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/service"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/web/handler/subscription"
)

func NewRouter(s service.SubscriptionService) *http.ServeMux {
	h := subscription.NewHandler(s)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /subscriptions", h.Create)
	mux.HandleFunc("GET /subscriptions", h.GetAll)
	mux.HandleFunc("GET /subscriptions/{id}", h.GetByID)
	mux.HandleFunc("PUT /subscriptions/{id}", h.Update)
	mux.HandleFunc("DELETE /subscriptions/{id}", h.Delete)
	mux.HandleFunc("GET /subscriptions/total-cost", h.GetTotalCost)

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	return mux
}
