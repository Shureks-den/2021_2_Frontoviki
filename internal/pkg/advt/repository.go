package advt

import "yula/internal/models"

//go:generate mockery -name=AdvtRepository

type AdvtRepository interface {
	SelectListAdvt(isSortedByPublichedDate bool, from, count int64) ([]*models.Advert, error)
	SelectAdvertsByPublisherId(publisherId int64, is_active bool, offset int64, limit int64) ([]*models.Advert, error)
	SelectAdvertsByCategory(categoryName string, from, count int64) ([]*models.Advert, error)
	SelectFavoriteAdverts(userId int64, from, count int64) ([]*models.Advert, error)

	Insert(advert *models.Advert) error
	SelectById(advertId int64) (*models.Advert, error)
	Update(newAdvert *models.Advert) error
	Delete(advertId int64) error

	InsertImages(advertId int64, newImages []string) error
	DeleteImages(images []string, advertId int64) error

	SelectFavorite(userId, advertId int64) (*models.Advert, error)
	InsertFavorite(userId, advertId int64) error
	DeleteFavorite(userId, advertId int64) error
}
