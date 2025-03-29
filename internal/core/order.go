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
    MenuItemID int64  `json:"menu_item_id" validate:"required,gt=0"`
    Category   string `json:"category" validate:"required"`
    Quantity   int64  `json:"quantity" validate:"required,gt=0"`
    PriceID    int64  `json:"price_id" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
    WaiterID            int64       `json:"waiter_id" validate:"required"`
    TimeCreated         time.Time   `json:"created_at" validate:"required"`
    TimeActualCompleted time.Time   `json:"actual_completed_at" validate:"required"`
    Comment             string      `json:"comment" validate:"required,min=5"`
    Items               []OrderItem `json:"items" validate:"required,dive"`
}

func (r *CreateOrderRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}
