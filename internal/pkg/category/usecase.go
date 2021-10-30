package category

import "yula/internal/models"

type CategoryUsecase interface {
	GetCategories() ([]*models.Category, error)
}
