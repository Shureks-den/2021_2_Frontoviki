package advt

import "yula/internal/models"

type AdvtRepository interface {
	SelectListAdvt(isSortedByPublichedDate bool, from, count int64) ([]*models.Advert, error)
	SelectAdvertsByPublisherId(publisherId int64, offset int64, limit int64) ([]*models.Advert, error)

	Insert(advert *models.Advert) error
	SelectById(advertId int64) (*models.Advert, error)
	Update(newAdvert *models.Advert) error
	Delete(advertId int64) error

	EditImages(advertId int64, newImages []string) error
}
