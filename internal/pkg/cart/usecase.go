package cart

import "yula/internal/models"

type CartUsecase interface {
	GetCart(userId int64) (*models.CartList, error)
	AddToCart(userId int64, singleCart *models.CartHandler) error
	UpdateCart(userId int64, singleCart *models.CartHandler, maxAmount int64) error
	RemoveFromCart(userId int64, advertId int64) error

	UpdateAllCart(userId int64) error
	ClearAllCart(userId int64) error
}
