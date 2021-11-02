package usecase

import (
	"strings"
	"time"
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

func (cu *CartUsecase) GetOrderFromCart(userId int64, advertId int64) (*models.Cart, error) {
	order, err := cu.cartRepository.Select(userId, advertId)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (cu *CartUsecase) GetCart(userId int64) ([]*models.Cart, error) {
	cart, err := cu.cartRepository.SelectAll(userId)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (cu *CartUsecase) AddToCart(userId int64, singleCart *models.CartHandler) error {
	return nil
}

func (cu *CartUsecase) UpdateCart(userId int64, singleCart *models.CartHandler, maxAmount int64) (*models.Cart, error) {
	_, err := cu.cartRepository.Select(userId, singleCart.AdvertId)
	newOneInCart := models.NewCart(userId, singleCart)

	switch err {
	case nil:
		if newOneInCart.Amount == 0 {
			err = cu.cartRepository.Delete(newOneInCart)
			return nil, err
		} else if newOneInCart.Amount > maxAmount {
			var genErr error = internalError.SetMaxCopies(maxAmount)
			return newOneInCart, genErr
		}

		err = cu.cartRepository.Update(newOneInCart)
		return newOneInCart, err

	case internalError.EmptyQuery:
		if newOneInCart.Amount == 0 {
			return nil, nil
		} else if newOneInCart.Amount > maxAmount {
			var genErr error = internalError.SetMaxCopies(maxAmount)
			return newOneInCart, genErr
		}

		err = cu.cartRepository.Insert(newOneInCart)
		return newOneInCart, err

	default:
		return nil, err

	}
}

func (cu *CartUsecase) RemoveFromCart(userId int64, advertId int64) error {
	return nil
}

func (cu *CartUsecase) UpdateAllCart(userId int64, cart []*models.CartHandler,
	adverts []*models.Advert) ([]*models.Cart, []*models.Advert, []string, error) {
	if len(cart) != len(adverts) {
		return nil, nil, nil, internalError.BadRequest
	}

	newCart := make([]*models.Cart, 0)
	newAdvert := make([]*models.Advert, 0)
	messages := make([]string, 0)
	for i := range cart {
		el, err := cu.UpdateCart(userId, cart[i], adverts[i].Amount)
		if err != nil && !strings.Contains(err.Error(), "not enough copies") {
			return nil, nil, nil, err
		}

		if el != nil {
			newCart = append(newCart, el)
			newAdvert = append(newAdvert, adverts[i])
			var msg string
			if err != nil {
				_, msg = internalError.ToMetaStatus(err)
			} else {
				msg = "ok"
			}
			messages = append(messages, msg)
		}
	}
	return newCart, newAdvert, messages, nil
}

func (cu *CartUsecase) ClearAllCart(userId int64) error {
	err := cu.cartRepository.DeleteAll(userId)
	if err == internalError.EmptyQuery {
		return nil
	}
	return err
}

func (cu *CartUsecase) MakeOrder(order *models.Cart, advert *models.Advert) error {
	if order.Amount == 0 || order.Amount > advert.Amount {
		return internalError.InvalidQuery
	}

	advert.Amount -= order.Amount
	if advert.Amount == 0 {
		advert.IsActive = false
		advert.DateClose = time.Now()
	}

	err := cu.cartRepository.Delete(order)
	return err
}
