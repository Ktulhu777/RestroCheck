package repository

import (
	"database/sql"
)

type OrderRepo interface{
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}
