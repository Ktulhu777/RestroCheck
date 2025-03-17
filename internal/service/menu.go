package service

import (
	"context"
	"restrocheck/internal/core"
	rp "restrocheck/internal/repository"
)

type MenuService interface {
	SaveMenu(ctx context.Context, req core.CreateMenuRequest) (int64, error)
}

type menuService struct {
	repo rp.MenuRepo
}

func NewMenuService(menuRepo rp.MenuRepo) MenuService {
	return &menuService{repo: menuRepo}
}


func (m *menuService) SaveMenu(ctx context.Context, req core.CreateMenuRequest) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}
	return m.repo.SaveMenu(ctx, req.Name, req.PhotoURL, req.CategoryID)
}