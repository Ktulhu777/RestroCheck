package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"restrocheck/internal/core"
	sv "restrocheck/internal/service"
	resp "restrocheck/pkg/response"
)
type OrderHandler struct {
	log     *slog.Logger
	service sv.OrderService
}

func NewOrderHandler(log *slog.Logger, service sv.OrderService) *OrderHandler {
	return &OrderHandler{
		log:     log,
		service: service,
	}
}

func (oh *OrderHandler) SaveOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.order.SaveOrder"

		log := oh.log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Received request to save order")
		var req core.CreateOrderRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			resp.RespondWithError(log, w, r, http.StatusBadRequest, err, "invalid JSON format")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id, err := oh.service.SaveOrder(ctx, req)
		if err != nil {
			var validateErr validator.ValidationErrors
			switch {
			case errors.As(err, &validateErr):
				resp.ValidationError(err.(validator.ValidationErrors), log, w, r, http.StatusBadRequest, err, "error validator")
				return
			default:
				resp.RespondWithError(log, w, r, http.StatusInternalServerError, err, "failed to save order")
				return
			}
		}

		log.Info("order added", slog.Int64("id", id))
		render.JSON(w, r, core.SaveOrderResponse{
			Response: resp.OK(),
			ID:       id,
		})
	}
}
