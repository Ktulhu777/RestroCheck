package core

import (
	"errors"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	nameRegex = regexp.MustCompile(`^[\p{L}\p{M}’\s-]+$`)
	validate  = validator.New()
)

// Общая структура
type Waiter struct {
	ID        int64   `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Phone     string  `json:"phone"`
	HireDate  string  `json:"hire_date"`
	Salary    float64 `json:"salary"`
}

// Структура для POST запроса
type CreateWaiterRequest struct {
	FirstName string  `json:"first_name" validate:"required"`
	LastName  string  `json:"last_name" validate:"required"`
	Phone     string  `json:"phone" validate:"required,e164"`
	HireDate  string  `json:"hire_date" validate:"required,datetime=2006-01-02"`
	Salary    float64 `json:"salary" validate:"required,numeric,gte=0"`
}

func (r *CreateWaiterRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}

	if !nameRegex.MatchString(r.FirstName) || !nameRegex.MatchString(r.LastName) {
		// TODO: поменяй потом ошибки разделяя на бизнес-логику и работы БД
		return errors.New("invalid characters in first or last name")
	}
	return nil
}

// Структура для PUT запроса
type UpdateWaiterRequest struct {
	FirstName      *string    `json:"first_name,omitempty" validate:"omitempty"`
	LastName       *string    `json:"last_name,omitempty" validate:"omitempty"`
	Phone          *string    `json:"phone,omitempty" validate:"omitempty,e164"`
	HireDate       *string    `json:"hire_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Salary         *float64   `json:"salary,omitempty" validate:"omitempty,numeric,gte=0"`
	ParsedHireDate *time.Time `json:"-"`
}

func (r *UpdateWaiterRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}

	if (r.FirstName != nil && !nameRegex.MatchString(*r.FirstName)) ||
		(r.LastName != nil && !nameRegex.MatchString(*r.LastName)) {
		// TODO: поменяй потом ошибки разделяя на бизнес-логику и работы БД
		return errors.New("invalid characters in name")
	}

	if r.Salary != nil && *r.Salary < 0 {
		// TODO: поменяй потом ошибки разделяя на бизнес-логику и работы БД

		return errors.New("salary must be greater than or equal to 0")
	}

	if r.HireDate != nil {
		loc, _ := time.LoadLocation("UTC")
		// TODO: поменяй потом ошибки разделяя на бизнес-логику и работы БД

		parsedDate, err := time.ParseInLocation("2006-01-02", *r.HireDate, loc)
		if err != nil {
			return errors.New("invalid hire date format, use YYYY-MM-DD")
		}
		r.ParsedHireDate = &parsedDate
	}

	return nil
}
