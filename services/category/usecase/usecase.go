package usecase

import (
	"yula/internal/models"
	"yula/services/category"
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
	return categories, err
}
