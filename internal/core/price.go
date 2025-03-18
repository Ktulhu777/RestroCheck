package core

type CreatePriceRequest struct {
	MenuItemID int64  `json:"menu_item_id" validate:"required"`
	Size       string `json:"size" validate:"required"`
	Price      int64  `json:"price" validate:"required,numeric,gte=0"`
}

func (r *CreatePriceRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}