package subscription

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"log/slog"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/models"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/repository"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/service"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/validation"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/pkg/monthyear"
)

type Handler struct {
	Service   service.SubscriptionService
	Validator *validation.Validator
}

func NewHandler(service service.SubscriptionService) Handler {
	return Handler{
		Service:   service,
		Validator: validation.New(),
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSubscriptionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := h.Validator.Struct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.Service.CreateSubscription(r.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrSubscriptionAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, service.ErrInvalidDateRange) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Error("service failed to create subscription", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, resp, http.StatusCreated)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	subscriptionID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid subscription ID format", http.StatusBadRequest)
		return
	}

	var req models.UpdateSubscriptionRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := h.Validator.Struct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.Service.UpdateSubscription(r.Context(), subscriptionID, req)
	if err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrInvalidDateRange) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Error("service failed to update subscription", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, resp, http.StatusOK)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.Service.GetAllSubscriptions(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	h.writeJSONResponse(w, resp, http.StatusOK)
}

func (h *Handler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	// Don't bother adding custom converters for uuid and monthyear to use tag parsing just in one place,
	// so parse it manually
	req, err := h.parseTotalCostRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Validator.Struct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.Service.GetTotalCost(r.Context(), req)
	if err != nil {
		slog.Error("service failed to get total cost", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, resp, http.StatusOK)
}

func (h *Handler) parseTotalCostRequest(r *http.Request) (models.TotalCostRequest, error) {
	var req models.TotalCostRequest

	// Parse user_id
	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return req, errors.New("invalid user_id format")
		}
		req.UserID = &userID
	}

	// Parse service_name
	if serviceName := r.URL.Query().Get("service_name"); serviceName != "" {
		req.ServiceName = &serviceName
	}

	// Parse start_date
	if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
		var startDate monthyear.MonthYear
		if err := startDate.UnmarshalJSON([]byte(`"` + startDateStr + `"`)); err != nil {
			return req, errors.New("invalid start_date format, expected MM-YYYY")
		}
		req.StartDate = &startDate
	}

	// Parse end_date
	if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
		var endDate monthyear.MonthYear
		if err := endDate.UnmarshalJSON([]byte(`"` + endDateStr + `"`)); err != nil {
			return req, errors.New("invalid end_date format, expected MM-YYYY")
		}
		req.EndDate = &endDate
	}

	return req, nil
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	subscriptionID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid subscription ID format", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteSubscription(r.Context(), subscriptionID)
	if err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		slog.Error("service failed to delete subscription", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid subscription ID format", http.StatusBadRequest)
		return
	}

	resp, err := h.Service.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		slog.Error("service failed to get subscription", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, resp, http.StatusOK)
}

func (h *Handler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to write JSON response", "error", err)
	}
}
