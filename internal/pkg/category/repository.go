package category

import "yula/internal/models"

type CategoryRepository interface {
	SelectCategories() ([]*models.Category, error)
}
