package advt

import (
	"mime/multipart"
	"yula/internal/models"
)

type AdvtUsecase interface {
	GetListAdvt(from int64, count int64, newest bool) ([]*models.Advert, error)
	GetAdvertListByPublicherId(publisherId int64, page *models.Page) ([]*models.Advert, error)

	AdvertsToShort(adverts []*models.Advert) []*models.AdvertShort

	CreateAdvert(userId int64, advert *models.Advert) error
	GetAdvert(advertId int64) (*models.Advert, error)
	UpdateAdvert(advertId int64, newAdvert *models.Advert) error
	DeleteAdvert(advertId int64, userId int64) error
	CloseAdvert(advertId int64, userId int64) error

	UploadImages(files []*multipart.FileHeader, advertId int64, userId int64) (*models.Advert, error)
}
