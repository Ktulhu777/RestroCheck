package core

import "errors"

var (
	ErrWaiterNotFound error = errors.New("waiter not found")  // Официант не найден
	ErrPhoneExists    error = errors.New("phone already exists") // Телефон уже существует
)
