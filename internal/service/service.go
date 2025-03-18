package service

import (
	"restrocheck/internal/repository"
)

type Service struct {
	Waiter   WaiterService
	Category CategoryService
	Menu     MenuService
	Price    PriceService
}

type Deps struct {
	Repos *repository.Repositories
}

func NewService(deps Deps) *Service {
	waiterService := NewWaiterService(deps.Repos.Waiters)
	categoryService := NewCategoryService(deps.Repos.Category)
	menuService := NewMenuService(deps.Repos.Menu)
	priceService := NewPriceService(deps.Repos.Price)

	return &Service{
		Waiter:   waiterService,
		Category: categoryService,
		Menu:     menuService,
		Price:    priceService,
	}
}
