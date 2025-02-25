package service

import (
	"restrocheck/internal/repository"
)

type Service struct {
	Waiter WaiterService
}

type Deps struct {
	Repos *repository.Repositories
}

func NewService(deps Deps) *Service {
	waiterService := NewWaiterService(deps.Repos.Waiters)

	return &Service{
		Waiter: waiterService,
	}
}