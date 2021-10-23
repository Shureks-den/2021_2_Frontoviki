package usecase

import (
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/cart"
)

type CartUsecase struct {
	cartRepository cart.CartRepository
}

func NewCartUsecase(cartRepository cart.CartRepository) cart.CartUsecase {
	return &CartUsecase{
		cartRepository: cartRepository,
	}
}

func (cu *CartUsecase) GetCart(userId int64) (*models.CartList, error) {
	return nil, nil
}

func (cu *CartUsecase) AddToCart(userId int64, singleCart *models.CartHandler) error {
	return nil
}

func (cu *CartUsecase) UpdateCart(userId int64, singleCart *models.CartHandler, maxAmount int64) error {
	_, err := cu.cartRepository.Select(userId, singleCart.AdvertId)
	newOneInCart := models.NewCart(userId, singleCart)

	switch err {
	case nil:
		if newOneInCart.Amount == 0 {
			err = cu.cartRepository.Delete(newOneInCart)
			return err
		} else if newOneInCart.Amount > maxAmount {
			var genErr error = internalError.SetMaxCopies(maxAmount)
			return genErr
		}

		err = cu.cartRepository.Update(newOneInCart)
		return err

	case internalError.EmptyQuery:
		if newOneInCart.Amount == 0 {
			return nil
		} else if newOneInCart.Amount > maxAmount {
			var genErr error = internalError.SetMaxCopies(maxAmount)
			return genErr
		}

		err = cu.cartRepository.Insert(newOneInCart)
		return err

	default:
		return err

	}
}

func (cu *CartUsecase) RemoveFromCart(userId int64, advertId int64) error {
	return nil
}

func (cu *CartUsecase) UpdateAllCart(userId int64) error {
	return nil
}

func (cu *CartUsecase) ClearAllCart(userId int64) error {
	return nil
}
