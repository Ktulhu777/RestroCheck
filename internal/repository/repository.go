package repository

import "database/sql"

type Repositories struct {
	Waiters  WaiterRepo
	Category CategoryRepo
	Menu     MenuRepo
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Waiters:  NewWaiterRepo(db),
		Category: NewCategoryRepo(db),
		Menu:     NewMenuRepo(db),
	}
}
