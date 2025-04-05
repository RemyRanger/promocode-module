package handler

import (
	"APIs/internal/common/server"
	"APIs/internal/services/promocode/ports"
	"errors"
	"net/http"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

const handler_name = "storages-handler"

var tracer = otel.Tracer("storages-handler-tracer")

type Handler struct {
	service ports.Service
}

func NewHandler(service ports.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Save Promocode
// (POST /promocodes)
func (h *Handler) SavePromocode(w http.ResponseWriter, r *http.Request) {
	// Start otlp tracer
	ctx, span := tracer.Start(r.Context(), handler_name)
	defer span.End()

	// Decode & Validate request body
	var body ports.PromocodeIn
	if code, err := server.BindAndValidate(r, &body); err != nil {
		server.RespondError(w, r, err, code)
		return
	}

	// Validate Restrictions
	if err := validateRestrictions(body.Restrictions); err != nil {
		server.RespondError(w, r, err, http.StatusBadRequest)
		return
	}

	// Map request body to model
	model, err := ports.DtoToModel(&body)
	if err != nil {
		server.RespondError(w, r, err, http.StatusBadRequest)
		return
	}

	// Call service
	modelCreated, err := h.service.SavePromocode(ctx, model)
	if err != nil {
		server.RespondError(w, r, err, http.StatusInternalServerError)
		return
	}

	// Map model to response body
	respBody, err := ports.ModelToDto(modelCreated)
	if err != nil {
		server.RespondError(w, r, err, http.StatusInternalServerError)
		return
	}

	server.RespondWithBody(w, r, respBody, http.StatusCreated)
}

// Validate Promocode
// (POST /promocodes/_validate)
func (h *Handler) ValidatePromocode(w http.ResponseWriter, r *http.Request) {
	// Start otlp tracer
	ctx, span := tracer.Start(r.Context(), handler_name)
	defer span.End()

	// Decode & Validate request body
	var body ports.PromocodeValidation
	if code, err := server.BindAndValidate(r, &body); err != nil {
		server.RespondError(w, r, err, code)
		return
	}

	// Call service
	reasons, err := h.service.ValidatePromocode(ctx, body.PromocodeName, body.Arguments.Age, body.Arguments.Town)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			server.RespondError(w, r, err, http.StatusNotFound)
			return
		}
		server.RespondError(w, r, err, http.StatusInternalServerError)
		return
	}

	// Map model to response body
	var respBody ports.PromocodeValidationResponse
	if len(reasons) == 0 {
		respBody = ports.PromocodeValidationResponse{
			PromocodeName: body.PromocodeName,
			Status:        ports.Accepted,
		}
	} else {
		respBody = ports.PromocodeValidationResponse{
			PromocodeName: body.PromocodeName,
			Status:        ports.Denied,
			Reasons:       &reasons,
		}
	}

	server.RespondWithBody(w, r, respBody, http.StatusOK)
}
