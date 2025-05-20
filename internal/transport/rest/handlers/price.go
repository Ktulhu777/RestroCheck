package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"restrocheck/internal/core"
	rp "restrocheck/internal/repository"
	sv "restrocheck/internal/service"
	resp "restrocheck/pkg/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type PriceHandler struct {
	log     *slog.Logger
	service sv.PriceService
}

func NewPriceHandler(log *slog.Logger, service sv.PriceService) *PriceHandler {
	return &PriceHandler{
		log:     log,
		service: service,
	}
}

func (ph *PriceHandler) SavePrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.price.SavePrice"

		log := ph.log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Received request to save price")

		var req core.CreatePriceRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			resp.RespondWithError(log, w, r, http.StatusBadRequest, err, "invalid JSON format")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id, err := ph.service.SavePice(ctx, req)
		if err != nil {
			var validateErr validator.ValidationErrors
			switch {
			case errors.As(err, &validateErr):
				resp.ValidationError(err.(validator.ValidationErrors), log, w, r, http.StatusBadRequest, err, "error validator")
				return
			case errors.Is(err, rp.ErrPriceUnique):
				resp.RespondWithError(log, w, r, http.StatusConflict, err, "price must unique")
				return
			case errors.Is(err, rp.ErrMenuIdDoesNotExists):
				resp.RespondWithError(log, w, r, http.StatusNotFound, err, "This dish does not exist")
				return
			case errors.Is(err, rp.ErrPriceInvalidSize):
				resp.RespondWithError(log, w, r, http.StatusConflict, err, "invalid size for price")
				return
			default:
				resp.RespondWithError(log, w, r, http.StatusInternalServerError, err, "server error")
				return
			}
		}

		log.Info("Price added in DataBase: ", slog.Int64("ID", id))
		render.JSON(w, r, core.SavePriceResponse{
			Response: resp.OK(),
			ID:       id,
		})
	}
}
