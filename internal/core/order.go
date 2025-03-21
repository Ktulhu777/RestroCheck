package core

import (
	"time"
)

type Order struct {
	ID                  int64     `json:"id"`
	WaiterID            int64     `json:"waiter_id"`
	TimeCreated         time.Time `json:"created_at"`
	TimeCompleted       time.Time `json:"completed_at"`      
	TimeActualCompleted time.Time `json:"actual_completed_at"`
	Comment             string    `json:"comment"`
}

type OrderItem struct {
	ID       int64  `json:"id"`
	OrderID  int64  `json:"order_id"`
	MenuID   int64  `json:"menu_item_id"`
	Category string `json:"category"`
	Quantity int32  `json:"quantity"`
	Price    int64  `json:"price"`
}

type CreateOrderRequest struct {
	WaiterID            int64     `json:"waiter_id" validate:"required"`
	TimeCreated         time.Time `json:"created_at" validate:"required"`
	TimeActualCompleted time.Time `json:"actual_completed_at" validate:"required"`
	Comment             string    `json:"comment" validate:"required,min=5"`
	MenuID              []int64   `json:"menu_id" validate:"required,dive,required,gt=0"`
	Category            string    `json:"category" validate:"required"`
}
