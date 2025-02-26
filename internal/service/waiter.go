package service

import (
	"context"
	"restrocheck/internal/core"
	"restrocheck/internal/repository"
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
	repo repository.WaiterRepo
}

func NewWaiterService(waiterRepo repository.WaiterRepo) WaiterService {
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

	id, err := w.repo.SaveWaiter(ctx, req.FirstName, req.LastName, req.Phone, hireDate, req.Salary)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (w *waiterService) FetchWaiter(ctx context.Context, id int64) (*core.Waiter, error) {
	wtr, err := w.repo.FetchWaiter(ctx, id)
	if err != nil {
		return nil, err
	}

	return wtr, nil
}

func (w *waiterService) ChangeWaiter(ctx context.Context, id int64, req core.UpdateWaiterRequest) (*core.Waiter, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	wtr, err := w.repo.ChangeWaiter(ctx, id, req.FirstName, req.LastName, req.Phone, req.ParsedHireDate, req.Salary)
	if err != nil {
		return nil, err
	}

	return wtr, nil
}

func (w *waiterService) RemoveWaiter(ctx context.Context, id int64) (int64, error) {
	id, err := w.repo.RemoveWaiter(ctx, id)
	if err != nil {
		return 0, err
	}
	return id, nil
}


func (w  *waiterService) FetchAllWaiters(ctx context.Context) ([]core.Waiter, error) {
	waiters, err := w.repo.FetchAllWaiters(ctx)
	if err != nil {
		return nil, err
	}
	return waiters, nil
}