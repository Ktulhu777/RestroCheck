package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type MenuRepo interface {
	SaveMenu(ctx context.Context, name, photoURL string, categoryID int64) (id int64, err error)
}

type MenuRepository struct {
	db *sql.DB
}

func NewMenuRepo(db *sql.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (m *MenuRepository) SaveMenu(ctx context.Context, name, photoURL string, categoryID int64) (id int64, err error) {
	const fn = "internal.repository.menu.SaveMenu"

	const query = `
		INSERT INTO menu(name, photo_url, category_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = m.db.QueryRowContext(ctx, query, name, photoURL, categoryID).Scan(&id)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return 0, fmt.Errorf("%s: %w", fn, ErrMenuNameExists)
			case "23503":
				return 0, fmt.Errorf("%s: foreign key violation (menu_category_id_fkey not found): %w", fn, ErrCategoryIdNotExists)
			}
			return 0, fmt.Errorf("%s: %w", fn, ErrMenuNameExists)
		}
		return 0, fmt.Errorf("%s: failed to insert menu name: %w", fn, err)
	}
	return id, nil
}
