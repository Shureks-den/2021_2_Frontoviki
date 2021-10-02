package advt

import "yula/internal/models"

type AdvtUsecase interface {
	GetListAdvt(from int64, count int64, newest bool) ([]*models.Advert, error)
}
