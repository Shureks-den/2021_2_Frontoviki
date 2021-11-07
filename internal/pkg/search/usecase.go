package search

import "yula/internal/models"

type SearchUsecase interface {
	SearchWithFilter(query *models.SearchFilter, page *models.Page) ([]*models.Advert, error)
}
