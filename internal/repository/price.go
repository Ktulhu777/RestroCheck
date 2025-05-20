package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type PriceRepo interface {
	SavePrice(ctx context.Context, menuItemID int64, size string, price int64) (id int64, err error)
}

type PriceRepository struct {
	db *sql.DB
}

func NewPriceRepo(db *sql.DB) *PriceRepository {
	return &PriceRepository{db: db}
}

func (p *PriceRepository) SavePrice(ctx context.Context, menuItemID int64, size string, price int64) (id int64, err error) {
	const fn = "internal.repository.price.SavePrice"

	const query = `
		INSERT INTO prices(menu_item_id, size, price)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = p.db.QueryRowContext(ctx, query, menuItemID, size, price).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return 0, fmt.Errorf("%s: %w", fn, ErrPriceUnique)
			case "23503":
				return 0, fmt.Errorf("%s: foreign key violation (menu_item_id not found): %w", fn, ErrMenuIdDoesNotExists)
			case "23514":
				return 0, fmt.Errorf("%s: %w (size=%q)", fn, ErrPriceInvalidSize, size)
			}
		}
		return 0, fmt.Errorf("%s: failed to insert price: %w", fn, err)
	}
	return id, nil
}
