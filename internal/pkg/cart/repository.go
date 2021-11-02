package cart

import "yula/internal/models"

//go:generate mockery -name=CartRepository

type CartRepository interface {
	Select(userId int64, advertId int64) (*models.Cart, error)
	Update(cart *models.Cart) error
	Insert(cart *models.Cart) error
	Delete(cart *models.Cart) error

	SelectAll(userId int64) ([]*models.Cart, error)
	DeleteAll(userId int64) error
}
