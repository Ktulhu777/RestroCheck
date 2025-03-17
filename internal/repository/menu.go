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
		// TODO: Доделать
		// "error": "internal.repository.menu.SaveMenu: failed to insert menu name: pq: INSERT или UPDATE в таблице \"menu\" нарушает ограничение внешнего ключа \"menu_category_id_fkey

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", fn, ErrMenuNameExists)
		}
		return 0, fmt.Errorf("%s: failed to insert menu name: %w", fn, err)
	}
	return id, nil
}
