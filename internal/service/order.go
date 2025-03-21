package service

import rp "restrocheck/internal/repository"

type OrderService interface{}

type orderService struct {
	repo rp.OrderRepo
}

func NewOrderService(orderRepo rp.OrderRepo) OrderService {
	return &orderService{repo: orderRepo}
}
