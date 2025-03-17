package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"restrocheck/internal/core"
	"restrocheck/internal/repository"
	"restrocheck/internal/service"
	resp "restrocheck/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type WaiterHandler struct {
	log     *slog.Logger
	service service.WaiterService
}

func NewWaiterHandler(log *slog.Logger, service service.WaiterService) *WaiterHandler {
	return &WaiterHandler{
		log:     log,
		service: service,
	}
}

func (wh *WaiterHandler) SaveWaiter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.waiter.SaveWaiter"

		log := wh.log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Received request to save waiter")
		var req core.CreateWaiterRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			resp.RespondWithError(log, w, r, http.StatusBadRequest, err, "invalid JSON format")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id, err := wh.service.SaveWaiter(ctx, req)
		if err != nil {
			var validateErr validator.ValidationErrors
			switch {
			case errors.As(err, &validateErr):
				resp.ValidationError(err.(validator.ValidationErrors), log, w, r, http.StatusBadRequest, err, "error validator")
				return
			case errors.Is(err, repository.ErrPhoneExists):
				resp.RespondWithError(log, w, r, http.StatusConflict, err, "phone already exists")
				return
			default:
				resp.RespondWithError(log, w, r, http.StatusInternalServerError, err, "failed to save waiter")
				return
			}
		}
		log.Info("waiter added", slog.Int64("id", id))
		render.JSON(w, r, core.SaveWaiterResponse{
			Response: resp.OK(),
			ID:       id,
		})
	}
}

func (wh *WaiterHandler) FetchWaiter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.waiter.FetchWaiter"

		log := wh.log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		pkStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(pkStr, 10, 64)
		if err != nil {
			resp.RespondWithError(log, w, r, http.StatusBadRequest, err, "invalid waiter ID")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		wtr, err := wh.service.FetchWaiter(ctx, id)
		if err != nil {
			if errors.Is(err, repository.ErrWaiterNotFound) {
				resp.RespondWithError(log, w, r, http.StatusNotFound, err, "waiter does not exist")
				return
			}
			resp.RespondWithError(log, w, r, http.StatusInternalServerError, err, "failed to fetch waiter")
			return
		}
		log.Info("waiter found", slog.Int64("waiter_id", id))
		render.JSON(w, r, core.FetchWaiterResponse{
			Response: resp.OK(),
			Waiter:   wtr,
		})
	}
}

func (wh *WaiterHandler) RemoveWaiter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.waiter.RemoveWaiter"

		log := wh.log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		pkStr := chi.URLParam(r, "id")
		pk, err := strconv.ParseInt(pkStr, 10, 64)
		if err != nil {
			resp.RespondWithError(log, w, r, http.StatusBadRequest, err, "invalid waiter ID")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id, err := wh.service.RemoveWaiter(ctx, pk)
		if err != nil {
			if errors.Is(err, repository.ErrWaiterNotFound) {
				resp.RespondWithError(log, w, r, http.StatusNotFound, err, "waiter does not exist")
				return
			}
			resp.RespondWithError(log, w, r, http.StatusInternalServerError, err, "failed to delete waiter")
			return
		}
		log.Info("waiter deleted successfully", slog.Int64("waiter_id", id))
		render.JSON(w, r, core.RemoveWaiterResponse{
			Response: resp.OK(),
			ID:       id,
		})
	}
}

func (wh *WaiterHandler) ChangeWaiter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.waiter.RemoveWaiter"

		log := wh.log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			resp.RespondWithError(log, w, r, http.StatusBadRequest, err, "invalid waiter ID")
			return
		}

		var req core.UpdateWaiterRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			resp.RespondWithError(log, w, r, http.StatusBadRequest, err, "invalid JSON format")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		wtr, err := wh.service.ChangeWaiter(ctx, id, req)
		if err != nil {
			var validateErr validator.ValidationErrors
			switch {
			case errors.As(err, &validateErr):
				resp.ValidationError(err.(validator.ValidationErrors), log, w, r, http.StatusBadRequest, err, "error validator")
				return
			case errors.Is(err, repository.ErrPhoneExists):
				resp.RespondWithError(log, w, r, http.StatusConflict, err, "phone number already exists")
				return
			case errors.Is(err, repository.ErrWaiterNotFound):
				resp.RespondWithError(log, w, r, http.StatusNotFound, err, "waiter not found")
				return
			default:
				resp.RespondWithError(log, w, r, http.StatusInternalServerError, err, "failed to update waiter")
				return
			}
		}
		log.Info("waiter updated successfully", slog.Int64("waiter_id", id))
		render.JSON(w, r, core.ChangeWaiterResponse{
			Response: resp.OK(),
			Waiter:   wtr,
		})
	}
}

func (wh *WaiterHandler) FetchAllWaiters() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.waiter.FetchAllWaiters"

		log := wh.log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		waiters, err := wh.service.FetchAllWaiters(ctx)
		if err != nil {
			if errors.Is(err, repository.ErrEmptyCollectionWaiter) {
				render.JSON(w, r, core.FetchAllWaiterResponse{
					Response: resp.OK(),
					Waiters:  []core.PartialWaiter{},
				})
				return
			}
			resp.RespondWithError(log, w, r, http.StatusInternalServerError, err, "failed to fetch all waiters")
			return
		}

		// Конвертируем полный список официантов в срез PartialWaiter
		waitersResponse := make([]core.PartialWaiter, len(waiters))
		for i, w := range waiters {
			waitersResponse[i] = core.PartialWaiter{
				ID:        w.ID,
				FirstName: w.FirstName,
				LastName:  w.LastName,
			}
		}

		render.JSON(w, r, core.FetchAllWaiterResponse{
			Response: resp.OK(),
			Waiters:  waitersResponse,
		})
	}
}
