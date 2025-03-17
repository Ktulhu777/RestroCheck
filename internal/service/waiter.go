package service

import (
	"context"
	"restrocheck/internal/core"
	rp "restrocheck/internal/repository"
	"time"
)

type WaiterService interface {
	SaveWaiter(ctx context.Context, req core.CreateWaiterRequest) (int64, error)
	FetchWaiter(ctx context.Context, id int64) (*core.Waiter, error)
	ChangeWaiter(ctx context.Context, id int64, req core.UpdateWaiterRequest) (*core.Waiter, error)
	RemoveWaiter(ctx context.Context, id int64) (int64, error)
	FetchAllWaiters(ctx context.Context) ([]core.Waiter, error)
}

type waiterService struct {
	repo rp.WaiterRepo
}

func NewWaiterService(waiterRepo rp.WaiterRepo) WaiterService {
	return &waiterService{repo: waiterRepo}
}

func (w *waiterService) SaveWaiter(ctx context.Context, req core.CreateWaiterRequest) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}

	hireDate, err := time.Parse("2006-01-02", req.HireDate)
	if err != nil {
		return 0, ErrInvalidDateFormat
	}

	return w.repo.SaveWaiter(ctx, req.FirstName, req.LastName, req.Phone, hireDate, req.Salary)
}

func (w *waiterService) FetchWaiter(ctx context.Context, id int64) (*core.Waiter, error) {
	return w.repo.FetchWaiter(ctx, id)
}

func (w *waiterService) ChangeWaiter(ctx context.Context, id int64, req core.UpdateWaiterRequest) (*core.Waiter, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return w.repo.ChangeWaiter(ctx, id, req.FirstName, req.LastName, req.Phone, req.ParsedHireDate, req.Salary)
}

func (w *waiterService) RemoveWaiter(ctx context.Context, id int64) (int64, error) {
	return w.repo.RemoveWaiter(ctx, id)
}

func (w *waiterService) FetchAllWaiters(ctx context.Context) ([]core.Waiter, error) {
	return w.repo.FetchAllWaiters(ctx)
}
