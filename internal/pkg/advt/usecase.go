package advt

import (
	"mime/multipart"
	"yula/internal/models"
)

//go:generate mockery -name=AdvtUsecase

type AdvtUsecase interface {
	GetListAdvt(from int64, count int64, newest bool) ([]*models.Advert, error)
	GetAdvertListByPublicherId(publisherId int64, is_active bool, page *models.Page) ([]*models.Advert, error)
	GetAdvertListByCategory(categoryName string, page *models.Page) ([]*models.Advert, error)

	AdvertsToShort(adverts []*models.Advert) []*models.AdvertShort

	CreateAdvert(userId int64, advert *models.Advert) error
	GetAdvert(advertId, userId int64, updateViews bool) (*models.Advert, error)
	UpdateAdvert(advertId int64, newAdvert *models.Advert) error
	DeleteAdvert(advertId int64, userId int64) error
	CloseAdvert(advertId int64, userId int64) error

	UploadImages(files []*multipart.FileHeader, advertId int64, userId int64) (*models.Advert, error)
	RemoveImages(images []string, advertId, userId int64) error

	GetFavoriteList(userId int64, page *models.Page) ([]*models.Advert, error)
	AddFavorite(userId int64, advertId int64) error
	RemoveFavorite(userId int64, advertId int64) error

	GetAdvertViews(advertId int64) (int64, error)

	GetPriceHistory(advertId int64) ([]*models.AdvertPrice, error)
	UpdateAdvertPrice(userId int64, adPrice *models.AdvertPrice) error
}
