package validation

import (
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/models"
)

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	v := &Validator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	v.validator.RegisterStructValidation(v.createSubscriptionRequest, models.CreateSubscriptionRequest{})
	v.validator.RegisterStructValidation(v.updateSubscriptionRequest, models.UpdateSubscriptionRequest{})
	v.validator.RegisterStructValidation(v.totalCostRequest, models.TotalCostRequest{})

	return v
}

// Struct validates a struct and returns validation errors
func (v *Validator) Struct(s any) error {
	return v.validator.Struct(s)
}

func (v *Validator) createSubscriptionRequest(sl validator.StructLevel) {
	req := sl.Current().Interface().(models.CreateSubscriptionRequest)

	// Validate that end_date is after start_date if provided
	if nil != req.EndDate {
		startTime := time.Time(*req.StartDate)
		endTime := time.Time(*req.EndDate)

		if endTime.Before(startTime) {
			sl.ReportError(req.EndDate, "end_date", "EndDate", "afterstart", "end date must be after start date")
		}
	}
}

func (v *Validator) updateSubscriptionRequest(sl validator.StructLevel) {
	req := sl.Current().Interface().(models.UpdateSubscriptionRequest)

	// Validate price if provided (must be non-negative)
	if req.Price != nil && *req.Price < 0 {
		sl.ReportError(req.Price, "price", "Price", "min", "price must be non-negative")
	}

	// Validate service name if provided (must not be empty)
	if req.ServiceName != nil && *req.ServiceName == "" {
		sl.ReportError(req.ServiceName, "service_name", "ServiceName", "required", "service name cannot be empty")
	}

	if req.ServiceName == nil && req.Price == nil && req.EndDate == nil {
		sl.ReportError(req, "request_body", "RequestBody", "at_least_one_required", "at least one field must be provided")
	}
}

func (v *Validator) totalCostRequest(sl validator.StructLevel) {
	req := sl.Current().Interface().(models.TotalCostRequest)

	// Validate date range if both dates are provided
	if req.StartDate != nil && req.EndDate != nil {
		startTime := time.Time(*req.StartDate)
		endTime := time.Time(*req.EndDate)

		if endTime.Before(startTime) {
			sl.ReportError(req.EndDate, "end_date", "EndDate", "afterstart", "end date must be after start date")
		}
	}

	// Validate service name if provided (must not be empty)
	if req.ServiceName != nil && *req.ServiceName == "" {
		sl.ReportError(req.ServiceName, "service_name", "ServiceName", "required", "service name cannot be empty")
	}
}
