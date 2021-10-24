package cart

import "yula/internal/models"

type CartUsecase interface {
	GetCart(userId int64) ([]*models.Cart, error)
	AddToCart(userId int64, singleCart *models.CartHandler) error
	UpdateCart(userId int64, singleCart *models.CartHandler, maxAmount int64) (*models.Cart, error)
	RemoveFromCart(userId int64, advertId int64) error

	UpdateAllCart(userId int64, cart []*models.CartHandler,
		adverts []*models.Advert) ([]*models.Cart, []*models.Advert, []string, error)
	ClearAllCart(userId int64) error
}
