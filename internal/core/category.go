package core

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}


func (r *CreateCategoryRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}