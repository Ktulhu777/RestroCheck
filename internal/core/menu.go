package core

// Общая структура
type Menu struct {
	Name       string `json:"name"`
	PhotoURL   string `json:"photo_url"`
	CategoryID int64  `json:"category_id"`
}

// Структура для POST запроса
type CreateMenuRequest struct {
	Name       string `json:"name" validate:"required"`
	PhotoURL   string `json:"photo_url" validate:"required"`
	CategoryID int64  `json:"category_id" validate:"required,numeric,gte=0"`
}

func (r *CreateMenuRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}
