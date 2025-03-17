package service

import (
	"context"

	"restrocheck/internal/core"
	rp "restrocheck/internal/repository"
)

type CategoryService interface {
	SaveCategory(ctx context.Context, req core.CreateCategoryRequest) (int64, error)
}

type categoryService struct {
	repo rp.CategoryRepo
}

func NewCategoryService(categoryRepo rp.CategoryRepo) CategoryService {
	return &categoryService{repo: categoryRepo}
}

func (c *categoryService) SaveCategory(ctx context.Context, req core.CreateCategoryRequest) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}
	return c.repo.SaveCategory(ctx, req.Name)
}
