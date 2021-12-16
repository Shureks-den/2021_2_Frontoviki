package category

import "yula/internal/models"

//go:generate mockery -name=CategoryUsecase

type CategoryUsecase interface {
	GetCategories() ([]*models.Category, error)
}
