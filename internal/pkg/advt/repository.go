package advt

import "yula/internal/models"

type AdvtRepository interface {
	SelectListAdvt(isSortedByPublichedDate bool, from, count int64) ([]*models.Advert, error)
}
