package service

import (
	"context"

	"restrocheck/internal/core"
	rp "restrocheck/internal/repository"
)

type OrderService interface {
	SaveOrder(ctx context.Context, req core.CreateOrderRequest) (int64, error)
}

type orderService struct {
	repo rp.OrderRepo
}

func NewOrderService(orderRepo rp.OrderRepo) OrderService {
	return &orderService{repo: orderRepo}
}

func (o *orderService) SaveOrder(ctx context.Context, req core.CreateOrderRequest) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}
	return o.repo.SaveOrder(ctx, req.WaiterID, req.TimeCreated, req.TimeActualCompleted, req.Items, req.Comment)
}