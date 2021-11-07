package search

import "yula/internal/models"

type SearchRepository interface {
	SelectWithFilter(search *models.SearchFilter, from, count int64) ([]*models.Advert, error)
}
