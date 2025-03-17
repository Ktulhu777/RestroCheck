package repository

import "errors"

var (
	ErrWaiterNotFound        error = errors.New("waiter not found")     // Официант не найден
	ErrPhoneExists           error = errors.New("phone already exists") // Телефон уже существует
	ErrEmptyCollectionWaiter error = errors.New("waiter list is empty") // Список официантов пуст
	ErrCategoryNameExists    error = errors.New("category name already exists")
)
