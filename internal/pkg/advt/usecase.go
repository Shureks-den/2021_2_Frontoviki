package advt

import "yula/internal/models"

type AdvtUsecase interface {
	GetListAdvt(from int64, count int64, newest bool) ([]*models.Advert, error)

	CreateAdvert(userId int64, advert *models.Advert) error
	GetAdvert(advertId int64) (*models.Advert, error)
	UpdateAdvert(advertId int64, newAdvert *models.Advert) error
	DeleteAdvert(advertId int64, userId int64) error
	CloseAdvert(advertId int64, userId int64) error
}
