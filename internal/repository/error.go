package repository

import "errors"

var (
	ErrWaiterNotFound        error = errors.New("waiter not found")     // Официант не найден
	ErrPhoneExists           error = errors.New("phone already exists") // Телефон уже существует
	ErrEmptyCollectionWaiter error = errors.New("waiter list is empty") // Список официантов пуст
	ErrCategoryNameExists    error = errors.New("category name already exists")
	ErrMenuNameExists        error = errors.New("menu name already exists")
	ErrPriceUnique           error = errors.New("price must unique")
	ErrMenuIdDoesNotExists   error = errors.New("this dish does not exist")
)
