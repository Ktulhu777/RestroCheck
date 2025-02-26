package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"restrocheck/internal/core"
	"restrocheck/pkg/storage"
	"time"
)

type WaiterRepo interface {
	SaveWaiter(ctx context.Context, firstName, lastName, phone string, hireDate time.Time, salary float64) (int64, error)
	FetchWaiter(ctx context.Context, pk int64) (*core.Waiter, error)
	ChangeWaiter(ctx context.Context, id int64, firstName, lastName, phone *string, hireDate *time.Time, salary *float64) (*core.Waiter, error)
	RemoveWaiter(ctx context.Context, pk int64) (int64, error)
	FetchAllWaiters(ctx context.Context) ([]core.Waiter, error)
}

type WaiterRepository struct {
	db *sql.DB
}

func NewWaiterRepo(db *sql.DB) *WaiterRepository {
	return &WaiterRepository{db: db}
}

func (w *WaiterRepository) SaveWaiter(
	ctx context.Context,
	firstName, lastName, phone string,
	hireDate time.Time,
	salary float64,
) (int64, error) {
	const fn = "internal.repository.waiter.SaveWaiter"

	const query = `
		INSERT INTO waiters (first_name, last_name, phone, hire_date, salary) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	var pk int64
	err := w.db.QueryRowContext(ctx, query, firstName, lastName, phone, hireDate, salary).Scan(&pk)

	if err != nil {
		if storage.IsDuplicatePhoneError(err) {
			return 0, fmt.Errorf("%s: %w", fn, core.ErrPhoneExists)
		}
		return 0, fmt.Errorf("%s: failed to insert waiter: %w", fn, err)
	}

	return pk, nil
}

func (w *WaiterRepository) FetchWaiter(
	ctx context.Context,
	pk int64,
) (*core.Waiter, error) {
	const fn = "internal.repository.waiter.FetchWaiter"
	const query = `
		SELECT id, first_name, last_name, phone, hire_date, salary
		FROM waiters	
		WHERE id = $1
	`

	var wtr core.Waiter
	err := w.db.QueryRowContext(ctx, query, pk).Scan(
		&wtr.ID,
		&wtr.FirstName,
		&wtr.LastName,
		&wtr.Phone,
		&wtr.HireDate,
		&wtr.Salary,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, core.ErrWaiterNotFound)
		}
		return nil, fmt.Errorf("%s: failed to fetch waiter: %w", fn, err)
	}
	return &wtr, nil
}

func (w *WaiterRepository) ChangeWaiter(
	ctx context.Context,
	id int64,
	firstName,
	lastName,
	phone *string,
	hireDate *time.Time,
	salary *float64,
) (*core.Waiter, error) {
	const fn = "internal.repository.waiter.ChangeWaiter"

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to start transaction: %w", fn, err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const query = `
		UPDATE waiters
		SET 
			first_name = COALESCE($2, first_name),
			last_name = COALESCE($3, last_name),
			phone = COALESCE($4, phone),
			hire_date = COALESCE($5, hire_date),
			salary = COALESCE($6, salary)
		WHERE id = $1
		RETURNING id, first_name, last_name, phone, hire_date, salary;
	`

	var wtr core.Waiter
	err = tx.QueryRowContext(ctx, query, id, firstName, lastName, phone, hireDate, salary).Scan(
		&wtr.ID,
		&wtr.FirstName,
		&wtr.LastName,
		&wtr.Phone,
		&wtr.HireDate,
		&wtr.Salary,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, core.ErrWaiterNotFound)
		}

		if storage.IsDuplicatePhoneError(err) {
			return nil, fmt.Errorf("%s: %w", fn, core.ErrPhoneExists)
		}
		return nil, fmt.Errorf("%s: failed to update waiter: %w", fn, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", fn, err)
	}

	return &wtr, nil
}

func (w *WaiterRepository) RemoveWaiter(ctx context.Context, pk int64) (int64, error) {
	const fn = "internal.repository.waiter.RemoveWaiter"

	const query = `
		DELETE FROM waiters	
		WHERE id = $1
		RETURNING id
	`

	var deletePK int64
	err := w.db.QueryRowContext(ctx, query, pk).Scan(&deletePK)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", fn, core.ErrWaiterNotFound)
		}
		return 0, fmt.Errorf("%s: failed to delete waiter: %w", fn, err)
	}

	return deletePK, nil
}

func (w *WaiterRepository) FetchAllWaiters(ctx context.Context) ([]core.Waiter, error) {
	const fn = "internal.repository.waiter.FetchAllWaiters"

	const query = `
		SELECT id, first_name, last_name
		FROM waiters
	`

	rows, err := w.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to fetch waiters: %w", fn, err)
	}
	defer rows.Close()

	var waiters []core.Waiter

	for rows.Next() {
		var waiter core.Waiter
		if err := rows.Scan(&waiter.ID, &waiter.FirstName, &waiter.LastName); err != nil {
			return nil, fmt.Errorf("%s: failed to scan waiter: %w", fn, err)
		}
		waiters = append(waiters, waiter)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: error iterating over waiters: %w", fn, err)
	}

	if len(waiters) == 0 {
		return nil, fmt.Errorf("%s: %w", fn, core.ErrEmptyCollectionWaiter)
	}

	return waiters, nil
}
