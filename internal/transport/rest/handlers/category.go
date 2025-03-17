package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"restrocheck/internal/core"
	"restrocheck/internal/repository"
	"restrocheck/internal/service"
	resp "restrocheck/pkg/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type CategoryHandler struct {
	log     *slog.Logger
	service service.CategoryService
}

func NewCategoryHandler(log *slog.Logger, service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		log: log,
		service: service,
	}
}

func (ch *CategoryHandler) SaveCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.category.SaveCategory"

		log := ch.log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Received request to save category")

		var req core.CreateCategoryRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			resp.RespondWithError(log, w, r, http.StatusBadRequest, err, "invalid JSON format")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id, err := ch.service.SaveCategory(ctx, req)
		if err != nil {
			var validateErr validator.ValidationErrors
			switch {
			case errors.As(err, &validateErr):
				resp.ValidationError(err.(validator.ValidationErrors), log, w, r, http.StatusBadRequest, err, "error validator")
				return
			case errors.Is(err, repository.ErrCategoryNameExists):
				resp.RespondWithError(log, w, r, http.StatusConflict, err, "category already exists")
				return
			default:
				resp.RespondWithError(log, w, r, http.StatusInternalServerError, err, "failed to save category")
				return
			}
		}
		log.Info("category added", slog.Int64("id", id))
		render.JSON(w, r, core.SaveCategoryResponse{
			Response: resp.OK(),
			ID:       id,
		})
	}
}
