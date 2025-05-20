package service

import (
	"context"
	"restrocheck/internal/core"
	rp "restrocheck/internal/repository"
)

type PriceService interface {
	SavePice(ctx context.Context, req core.CreatePriceRequest) (int64, error)
}

type priceService struct {
	repo rp.PriceRepo
}

func NewPriceService(priceRepo rp.PriceRepo) PriceService {
	return &priceService{repo: priceRepo}
}

func (p *priceService) SavePice(ctx context.Context, req core.CreatePriceRequest) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}
	return p.repo.SavePrice(ctx, req.MenuItemID, req.Size, req.Price)
}
