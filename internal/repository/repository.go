package repository

import "database/sql"

type Repositories struct {
	Waiters WaiterRepo
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Waiters: NewWaiterRepo(db),
	}
}
