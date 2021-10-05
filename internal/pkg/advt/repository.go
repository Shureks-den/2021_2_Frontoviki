package advt

import "yula/internal/models"

type AdvtRepository interface {
	SelectListAdvt(isSortedByPublichedDate bool, from, count int64) ([]*models.Advert, error)
	Insert(advert *models.Advert) error
	SelectById(advertId int64) (*models.Advert, error)
	Update(newAdvert *models.Advert) error
}
