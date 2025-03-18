package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"restrocheck/internal/core"
	rp "restrocheck/internal/repository"
	sv "restrocheck/internal/service"
	resp "restrocheck/pkg/response"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type MenuHandler struct {
	log     *slog.Logger
	service sv.MenuService
}

func NewMenuHandler(log *slog.Logger, service sv.MenuService) *MenuHandler {
	return &MenuHandler{
		log:     log,
		service: service,
	}
}

func (mh *MenuHandler) SaveMenu() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.menu.SaveMenu"

		log := mh.log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("Received request to save menu")

		var req core.CreateMenuRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			resp.RespondWithError(log, w, r, http.StatusBadRequest, err, "invalid JSON format")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id, err := mh.service.SaveMenu(ctx, req)

		if err != nil {
			var validateErr validator.ValidationErrors
			switch {
			case errors.As(err, &validateErr):
				resp.ValidationError(err.(validator.ValidationErrors), log, w, r, http.StatusBadRequest, err, "error validator")
				return
			case errors.Is(err, rp.ErrMenuNameExists):
				resp.RespondWithError(log, w, r, http.StatusConflict, err, "menu name already exists")
				return
			default:
				resp.RespondWithError(log, w, r, http.StatusInternalServerError, err, "failed to save menu")
				return
			}
		}

		log.Info("menu added", slog.Int64("id", id))
		render.JSON(w, r, core.SaveMenuResponse{
			Response: resp.OK(),
			ID:       id,
		})
	}
}
