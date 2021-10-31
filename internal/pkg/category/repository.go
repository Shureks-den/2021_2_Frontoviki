package category

import "yula/internal/models"

//go:generate mockery -name=CategoryRepository

type CategoryRepository interface {
	SelectCategories() ([]*models.Category, error)
}
