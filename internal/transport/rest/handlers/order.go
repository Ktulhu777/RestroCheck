package handlers

import (
	"log/slog"
	sv "restrocheck/internal/service"
)

type OrderHandler struct {
	log     *slog.Logger
	service sv.OrderService
}

func NewOrderHandler(log *slog.Logger, service sv.OrderService) *OrderHandler {
	return &OrderHandler{
		log: log,
		service: service,
	}
}

