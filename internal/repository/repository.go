package repository

import "database/sql"

type Repositories struct {
	Waiters  WaiterRepo
	Category CategoryRepo
	Menu     MenuRepo
	Price    PriceRepo
	Order    OrderRepo
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Waiters:  NewWaiterRepo(db),
		Category: NewCategoryRepo(db),
		Menu:     NewMenuRepo(db),
		Price:    NewPriceRepo(db),
		Order:    NewOrderRepo(db),
	}
}
