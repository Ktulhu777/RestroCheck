package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restrocheck/internal/core"
	"time"
)

type OrderRepo interface {
	SaveOrder(
		ctx context.Context,
		waiterID int64,
		timeCreated,
		timeActualCompleted time.Time,
		Items []core.OrderItem,
		comment string,
	) (id int64, err error)
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) SaveOrder(
	ctx context.Context,
	waiterID int64,
	timeCreated,
	timeActualCompleted time.Time,
	Items []core.OrderItem,
	comment string,
) (id int64, err error) {
	const fn = "internal,repository.order.SaveOrder"

	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	timeCompleted := timeActualCompleted.Sub(timeCreated)
	timeCompletedStr := fmt.Sprintf("%02d:%02d:%02d",
		int(timeCompleted.Hours()),
		int(timeCompleted.Minutes())%60,
		int(timeCompleted.Seconds())%60,
	)
	const query = `
		INSERT INTO orders(waiter_id, created_at, completed_at, actual_completed_at, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, query, waiterID, timeCreated, timeCompletedStr, timeActualCompleted, comment).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := tx.PrepareContext(ctx, `
    INSERT INTO order_items(order_id, menu_item_id, category, quantity, price)
    SELECT $1, menu.id, $2, $3, prices.price
    FROM menu
    JOIN prices ON menu.id = prices.menu_item_id AND prices.id = $4
    WHERE menu.id = $5
`)

	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: prepare statement error: %w", fn, err)
	}

	for _, item := range Items {
		res, err := stmt.ExecContext(ctx, id, item.Category, item.Quantity, item.PriceID, item.MenuItemID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: exec error for item %v: %w", fn, item, err)
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			tx.Rollback()
			return 0, fmt.Errorf("%s: no rows inserted for item %v", fn, item)
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}
