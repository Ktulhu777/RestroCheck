package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type CategoryRepo interface {
	SaveCategory(ctx context.Context, name string) (id int64, err error)
}

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (c *CategoryRepository) SaveCategory(ctx context.Context, name string) (id int64, err error) {
	const fn = "internal.repository.category.SaveCategory"

	const query = `
		INSERT INTO categories(name)
		VALUES ($1)
		RETURNING id
	`
	err = c.db.QueryRowContext(ctx, query, name).Scan(&id)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", fn, ErrCategoryNameExists)
		}
		return 0, fmt.Errorf("%s: failed to insert category name: %w", fn, err)
	}
	return id, nil
}
