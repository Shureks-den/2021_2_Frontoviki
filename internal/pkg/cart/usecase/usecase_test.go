package usecase

import (
	"testing"
	"yula/internal/models"
	"yula/internal/pkg/cart/mocks"

	myerr "yula/internal/error"

	"github.com/stretchr/testify/assert"
)

func TestGetOrderFromCartSuccess(t *testing.T) {
	cart := models.Cart{
		UserId:   1,
		AdvertId: 2,
		Amount:   42,
	}
	cr := mocks.CartRepository{}
	cr.On("Select", int64(1), int64(2)).Return(&cart, nil)

	cu := NewCartUsecase(&cr)
	cartRes, err := cu.GetOrderFromCart(1, 2)
	assert.Equal(t, cart, *cartRes)
	assert.Nil(t, err)
}

func TestGetOrderFromCartFail(t *testing.T) {
	cr := mocks.CartRepository{}
	cr.On("Select", int64(-2), int64(2)).Return(nil, myerr.InternalError)

	cu := NewCartUsecase(&cr)
	cartRes, err := cu.GetOrderFromCart(-2, 2)
	assert.Equal(t, err, myerr.InternalError)
	assert.Nil(t, cartRes)
}

func TestGetCartSuccess(t *testing.T) {
	cart := []*models.Cart{
		{
			UserId:   1,
			AdvertId: 2,
			Amount:   9,
		},
		{
			UserId:   1,
			AdvertId: 15,
			Amount:   7,
		},
	}
	cr := mocks.CartRepository{}
	cr.On("SelectAll", int64(1)).Return(cart, nil)

	cu := NewCartUsecase(&cr)
	cartRes, err := cu.GetCart(1)
	assert.Equal(t, cartRes, cart)
	assert.Nil(t, err)
}

func TestGetCartFail(t *testing.T) {
	cr := mocks.CartRepository{}
	cr.On("SelectAll", int64(-21)).Return(nil, myerr.InternalError)

	cu := NewCartUsecase(&cr)
	cartRes, err := cu.GetCart(-21)
	assert.Equal(t, err, myerr.InternalError)
	assert.Nil(t, cartRes)
}

func TestUpdateCartInsertSuccess(t *testing.T) {
	cartHandler := models.CartHandler{
		AdvertId: 2,
		Amount:   5,
	}
	cart := models.Cart{
		UserId:   1,
		AdvertId: 2,
		Amount:   5,
	}
	cr := mocks.CartRepository{}
	cr.On("Select", int64(1), int64(2)).Return(nil, myerr.EmptyQuery)
	cr.On("Insert", &cart).Return(nil)

	cu := NewCartUsecase(&cr)
	cartRes, err := cu.UpdateCart(1, &cartHandler, 6)
	assert.Equal(t, *cartRes, cart)
	assert.Nil(t, err)
}

func TestUpdateCartInsertFailAmount0(t *testing.T) {
	cartHandler := models.CartHandler{
		AdvertId: 2,
		Amount:   0,
	}
	cr := mocks.CartRepository{}
	cr.On("Select", int64(1), int64(2)).Return(nil, myerr.EmptyQuery)

	cu := NewCartUsecase(&cr)
	cartRes, err := cu.UpdateCart(1, &cartHandler, 5)
	assert.Nil(t, cartRes)
	assert.Nil(t, err)
}

func TestUpdateCartInsertFailAmountGTAmountMax(t *testing.T) {
	cartHandler := models.CartHandler{
		AdvertId: 2,
		Amount:   8,
	}
	cart := models.Cart{
		UserId:   1,
		AdvertId: 2,
		Amount:   8,
	}
	cr := mocks.CartRepository{}
	cr.On("Select", int64(1), int64(2)).Return(nil, myerr.EmptyQuery)

	cu := NewCartUsecase(&cr)
	cartRes, err := cu.UpdateCart(1, &cartHandler, 6)
	assert.Equal(t, *cartRes, cart)
	assert.Equal(t, err, myerr.SetMaxCopies(6))
}

func TestUpdateCartUpdateSuccess(t *testing.T) {
	cartHandler := models.CartHandler{
		AdvertId: 2,
		Amount:   5,
	}
	cart := models.Cart{
		UserId:   1,
		AdvertId: 2,
		Amount:   5,
	}
	cr := mocks.CartRepository{}
	cr.On("Select", int64(1), int64(2)).Return(nil, nil)
	cr.On("Update", &cart).Return(nil)

	cu := NewCartUsecase(&cr)
	cartRes, err := cu.UpdateCart(1, &cartHandler, 6)
	assert.Equal(t, *cartRes, cart)
	assert.Nil(t, err)
}

func TestUpdateCartUpdateFailAmount0(t *testing.T) {
	cartHandler := models.CartHandler{
		AdvertId: 2,
		Amount:   0,
	}
	cart := models.Cart{
		UserId:   1,
		AdvertId: 2,
		Amount:   0,
	}
	cr := mocks.CartRepository{}
	cr.On("Select", int64(1), int64(2)).Return(nil, nil)
	cr.On("Delete", &cart).Return(nil)

	cu := NewCartUsecase(&cr)
	cartRes, err := cu.UpdateCart(1, &cartHandler, 5)
	assert.Nil(t, cartRes)
	assert.Nil(t, err)
}

func TestUpdateCartUpdateFailAmountGTAmountMax(t *testing.T) {
	cartHandler := models.CartHandler{
		AdvertId: 2,
		Amount:   8,
	}
	cart := models.Cart{
		UserId:   1,
		AdvertId: 2,
		Amount:   8,
	}
	cr := mocks.CartRepository{}
	cr.On("Select", int64(1), int64(2)).Return(nil, nil)

	cu := NewCartUsecase(&cr)
	cartRes, err := cu.UpdateCart(1, &cartHandler, 6)
	assert.Equal(t, *cartRes, cart)
	assert.Equal(t, err, myerr.SetMaxCopies(6))
}

func TestUpdateAllCartSuccess(t *testing.T) {
	cartH := []*models.CartHandler{
		{
			AdvertId: 2,
			Amount:   8,
		},
	}
	ads := []*models.Advert{
		{
			Id:     32,
			Name:   "aboba",
			Amount: 10,
		},
	}
	cart := models.Cart{
		UserId:   1,
		AdvertId: 2,
		Amount:   8,
	}

	expCart := []*models.Cart{
		{
			UserId:   1,
			AdvertId: 2,
			Amount:   8,
		},
	}
	expAds := []*models.Advert{
		{
			Id:     32,
			Name:   "aboba",
			Amount: 10,
		},
	}

	cr := mocks.CartRepository{}
	cr.On("Select", int64(1), int64(2)).Return(nil, nil)
	cr.On("Update", &cart).Return(nil)

	cu := NewCartUsecase(&cr)
	newCart, newAds, _, _ := cu.UpdateAllCart(1, cartH, ads)
	assert.Equal(t, expCart, newCart)
	assert.Equal(t, expAds, newAds)
}

func TestUpdateAllCartFailSomeAmountIsGTAmountMax(t *testing.T) {
	cart := []*models.CartHandler{
		{
			AdvertId: 2,
			Amount:   12,
		},
	}
	ads := []*models.Advert{
		{
			Id:     32,
			Name:   "aboba",
			Amount: 10,
		},
	}

	expCart := []*models.Cart{
		{
			UserId:   1,
			AdvertId: 2,
			Amount:   12,
		},
	}
	expAds := []*models.Advert{
		{
			Id:     32,
			Name:   "aboba",
			Amount: 10,
		},
	}

	cr := mocks.CartRepository{}
	cr.On("Select", int64(1), int64(2)).Return(nil, nil)

	cu := NewCartUsecase(&cr)
	newCart, newAds, msg, err := cu.UpdateAllCart(1, cart, ads)
	assert.Equal(t, expCart, newCart)
	assert.Equal(t, expAds, newAds)
	assert.Equal(t, msg[0], "not enough copies. max amount: 10")
	assert.Nil(t, err)
}

func TestUpdateAllCartFailSomeAmountIs0(t *testing.T) {
	cartH := []*models.CartHandler{
		{
			AdvertId: 2,
			Amount:   0,
		},
	}
	ads := []*models.Advert{
		{
			Id:     32,
			Name:   "aboba",
			Amount: 10,
		},
	}

	cart := models.Cart{
		UserId:   1,
		AdvertId: 2,
		Amount:   0,
	}

	cr := mocks.CartRepository{}
	cr.On("Select", int64(1), int64(2)).Return(nil, nil)
	cr.On("Delete", &cart).Return(myerr.DatabaseError)

	cu := NewCartUsecase(&cr)
	newCart, newAds, msg, err := cu.UpdateAllCart(1, cartH, ads)
	assert.Nil(t, newCart)
	assert.Nil(t, newAds)
	assert.Nil(t, msg)
	assert.NotNil(t, err)
}

func TestClearAllCartSuccess(t *testing.T) {
	cr := mocks.CartRepository{}
	cr.On("DeleteAll", int64(1)).Return(nil)

	cu := NewCartUsecase(&cr)
	err := cu.ClearAllCart(1)
	assert.Nil(t, err)
}

func TestClearAllCartFail(t *testing.T) {
	cr := mocks.CartRepository{}
	cr.On("DeleteAll", int64(1)).Return(myerr.NotExist)

	cu := NewCartUsecase(&cr)
	err := cu.ClearAllCart(1)
	assert.Equal(t, err, myerr.NotExist)
}

func TestMakeOrderSuccess(t *testing.T) {
	ad := models.Advert{
		Id:     32,
		Name:   "aboba",
		Amount: 10,
	}

	cart := models.Cart{
		UserId:   1,
		AdvertId: 32,
		Amount:   10,
	}

	cr := mocks.CartRepository{}
	cr.On("Delete", &cart).Return(nil)

	cu := NewCartUsecase(&cr)
	err := cu.MakeOrder(&cart, &ad)
	assert.Nil(t, err)
	assert.Equal(t, ad.Amount, int64(0))
	assert.Equal(t, ad.IsActive, false)
}

func TestMakeOrderFailAmountGTAmountMax(t *testing.T) {
	ad := models.Advert{
		Id:     32,
		Name:   "aboba",
		Amount: 2,
	}

	cart := models.Cart{
		UserId:   1,
		AdvertId: 32,
		Amount:   3,
	}

	cr := mocks.CartRepository{}

	cu := NewCartUsecase(&cr)
	err := cu.MakeOrder(&cart, &ad)
	assert.NotNil(t, err)
}

func TestMakeOrderFail(t *testing.T) {
	ad := models.Advert{
		Id:     32,
		Name:   "aboba",
		Amount: 4,
	}

	cart := models.Cart{
		UserId:   1,
		AdvertId: 32,
		Amount:   0,
	}

	cr := mocks.CartRepository{}

	cu := NewCartUsecase(&cr)
	err := cu.MakeOrder(&cart, &ad)
	assert.NotNil(t, err)
}
