package usecase

import (
	"yula/internal/models"
	"yula/internal/pkg/category"
)

type CategoryUsecase struct {
	categoryRepository category.CategoryRepository
}

func NewCategoryUsecase(categoryRepository category.CategoryRepository) category.CategoryUsecase {
	return &CategoryUsecase{
		categoryRepository: categoryRepository,
	}
}

func (cu *CategoryUsecase) GetCategories() ([]*models.Category, error) {
	categories, err := cu.categoryRepository.SelectCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}
