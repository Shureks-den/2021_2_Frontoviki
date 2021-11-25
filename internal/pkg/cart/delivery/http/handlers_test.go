package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"yula/internal/models"
	"yula/internal/pkg/middleware"

	advtMock "yula/internal/pkg/advt/mocks"

	cartMock "yula/internal/pkg/cart/mocks"

	userMock "yula/internal/pkg/user/mocks"

	myerr "yula/internal/error"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCartSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.HandleFunc("/one", ch.UpdateOneAdvertHandler).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cartHandler := models.CartHandler{
		AdvertId: 2,
		Amount:   8,
	}

	ad := models.Advert{
		Id:     cartHandler.AdvertId,
		Name:   "aboba",
		Amount: cartHandler.Amount,
	}

	newCart := models.NewCart(10, &cartHandler)

	au.On("GetAdvert", cartHandler.AdvertId, int64(0), false).Return(&ad, nil)
	cu.On("UpdateCart", int64(0), &cartHandler, ad.Amount).Return(newCart, nil)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(cartHandler)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart/one", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "successfully updated")
}

func TestCartFailGetAd(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.HandleFunc("/one", ch.UpdateOneAdvertHandler).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cartHandler := models.CartHandler{
		AdvertId: 2,
		Amount:   8,
	}

	au.On("GetAdvert", cartHandler.AdvertId, int64(0), false).Return(nil, myerr.EmptyQuery)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(cartHandler)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart/one", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 404)
	assert.Equal(t, Answer.Message, "empty rows")
}

func TestCartFailParse(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.HandleFunc("/one", ch.UpdateOneAdvertHandler).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart/one", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 400)
	assert.Equal(t, Answer.Message, "bad request")
}

func TestCartFailUpdateAd(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.HandleFunc("/one", ch.UpdateOneAdvertHandler).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cartHandler := models.CartHandler{
		AdvertId: 2,
		Amount:   8,
	}

	ad := models.Advert{
		Id:     cartHandler.AdvertId,
		Name:   "aboba",
		Amount: cartHandler.Amount,
	}

	au.On("GetAdvert", cartHandler.AdvertId, int64(0), false).Return(&ad, nil)
	cu.On("UpdateCart", int64(0), &cartHandler, ad.Amount).Return(nil, myerr.InternalError)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(cartHandler)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart/one", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestUpdateAllCartSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("", ch.UpdateAllCartHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cartH := []*models.CartHandler{
		{
			AdvertId: 2,
			Amount:   8,
		},
	}

	cart := []*models.Cart{
		{
			UserId:   0,
			AdvertId: 2,
			Amount:   8,
		},
	}

	ads := []*models.Advert{
		{
			Id:     2,
			Name:   "aboba",
			Amount: 10,
		},
	}

	ad := models.Advert{
		Id:     cartH[0].AdvertId,
		Name:   "aboba",
		Amount: cartH[0].Amount,
	}

	au.On("GetAdvert", cartH[0].AdvertId, int64(0), false).Return(&ad, nil)
	cu.On("UpdateAllCart", int64(0), mock.AnythingOfType("[]*models.CartHandler"), mock.AnythingOfType("[]*models.Advert")).Return(cart, ads, []string{"ok"}, nil)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(cartH)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, ((Answer.Body.(map[string]interface{})["hints"]).([]interface{})[0]), "ok")
}

func TestUpdateAllCartFailParse(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("", ch.UpdateAllCartHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 400)
	assert.Equal(t, Answer.Message, "bad request")
}

func TestUpdateAllCartFailGetAd(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("", ch.UpdateAllCartHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cartH := []*models.CartHandler{
		{
			AdvertId: 2,
			Amount:   8,
		},
	}

	au.On("GetAdvert", cartH[0].AdvertId, int64(0), false).Return(nil, myerr.InternalError)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(cartH)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestUpdateAllCartFailUpdateAllCart(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("", ch.UpdateAllCartHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cartH := []*models.CartHandler{
		{
			AdvertId: 2,
			Amount:   8,
		},
	}

	ad := models.Advert{
		Id:     cartH[0].AdvertId,
		Name:   "aboba",
		Amount: cartH[0].Amount,
	}

	au.On("GetAdvert", cartH[0].AdvertId, int64(0), false).Return(&ad, nil)
	cu.On("UpdateAllCart", int64(0), mock.AnythingOfType("[]*models.CartHandler"), mock.AnythingOfType("[]*models.Advert")).Return(nil, nil, nil, myerr.InternalError)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(cartH)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestGetAllCartSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("", ch.GetCartHandler).Methods(http.MethodGet, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cart := []*models.Cart{
		{
			UserId:   int64(0),
			AdvertId: 2,
			Amount:   8,
		},
	}

	ads := []*models.Advert{
		{
			Id:     int64(2),
			Name:   "aboba",
			Amount: 10,
		},
	}

	cu.On("GetCart", cart[0].UserId).Return(cart, nil)
	au.On("GetAdvert", cart[0].AdvertId, int64(0), false).Return(ads[0], nil)

	res, err := http.Get(fmt.Sprintf("%s/cart", srv.URL))
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Body.(map[string]interface{})["adverts"].([]interface{})[0].(map[string]interface{})["amount"], float64(ads[0].Amount))
}

func TestGetAllCartFailGetCart(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("", ch.GetCartHandler).Methods(http.MethodGet, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cart := []*models.Cart{
		{
			UserId:   int64(0),
			AdvertId: 2,
			Amount:   8,
		},
	}

	cu.On("GetCart", cart[0].UserId).Return(nil, myerr.InternalError)

	res, err := http.Get(fmt.Sprintf("%s/cart", srv.URL))
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestGetAllCartFailGetAd(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("", ch.GetCartHandler).Methods(http.MethodGet, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cart := []*models.Cart{
		{
			UserId:   int64(0),
			AdvertId: 2,
			Amount:   8,
		},
	}

	cu.On("GetCart", cart[0].UserId).Return(cart, nil)
	au.On("GetAdvert", cart[0].AdvertId, int64(0), false).Return(nil, myerr.InternalError)

	res, err := http.Get(fmt.Sprintf("%s/cart", srv.URL))
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestClearCartSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("/clear", ch.ClearCartHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cart := []*models.Cart{
		{
			UserId:   int64(0),
			AdvertId: 2,
			Amount:   8,
		},
	}

	cu.On("ClearAllCart", cart[0].UserId).Return(nil)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart/clear", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "cart cleared")
}

func TestClearCartFailed(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("/clear", ch.ClearCartHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cart := []*models.Cart{
		{
			UserId:   int64(0),
			AdvertId: 2,
			Amount:   8,
		},
	}

	cu.On("ClearAllCart", cart[0].UserId).Return(myerr.InternalError)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart/clear", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestCheckoutSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("/{id:[0-9]+}/checkout", ch.CheckoutHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cart := models.Cart{
		UserId:   int64(0),
		AdvertId: 2,
		Amount:   8,
	}

	ad := models.Advert{
		Id:     int64(2),
		Name:   "aboba",
		Amount: 10,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	cu.On("GetOrderFromCart", cart.UserId, cart.AdvertId).Return(&cart, nil)
	au.On("GetAdvert", ad.Id, int64(0), false).Return(&ad, nil)
	uu.On("GetById", cart.UserId).Return(&profile, nil)
	cu.On("MakeOrder", &cart, &ad).Return(nil)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart/2/checkout", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "order made successfully")
}

func TestCheckoutFailParseId(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("/{id:[0-9]+}/checkout", ch.CheckoutHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart/2418594151898483818491/checkout", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 400)
	assert.Equal(t, Answer.Message, "bad request")
}

func TestCheckoutFailGetOrderFromCart(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("/{id:[0-9]+}/checkout", ch.CheckoutHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cart := models.Cart{
		UserId:   int64(0),
		AdvertId: 2,
		Amount:   8,
	}

	cu.On("GetOrderFromCart", cart.UserId, cart.AdvertId).Return(nil, myerr.InternalError)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart/2/checkout", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestCheckoutFailGetAd(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("/{id:[0-9]+}/checkout", ch.CheckoutHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cart := models.Cart{
		UserId:   int64(0),
		AdvertId: 2,
		Amount:   8,
	}

	ad := models.Advert{
		Id:     int64(2),
		Name:   "aboba",
		Amount: 10,
	}

	cu.On("GetOrderFromCart", cart.UserId, cart.AdvertId).Return(&cart, nil)
	au.On("GetAdvert", ad.Id, int64(0), false).Return(nil, myerr.InternalError)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart/2/checkout", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}
func TestCheckoutFailGetById(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("/{id:[0-9]+}/checkout", ch.CheckoutHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cart := models.Cart{
		UserId:   int64(0),
		AdvertId: 2,
		Amount:   8,
	}

	ad := models.Advert{
		Id:     int64(2),
		Name:   "aboba",
		Amount: 10,
	}

	cu.On("GetOrderFromCart", cart.UserId, cart.AdvertId).Return(&cart, nil)
	au.On("GetAdvert", ad.Id, int64(0), false).Return(&ad, nil)
	uu.On("GetById", cart.UserId).Return(nil, myerr.InternalError)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart/2/checkout", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestCheckoutFailMakeOrder(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	cu := cartMock.CartUsecase{}
	ch := NewCartHandler(&cu, &uu, &au)

	router := mux.NewRouter().PathPrefix("/cart").Subrouter()
	router.HandleFunc("/{id:[0-9]+}/checkout", ch.CheckoutHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.LoggerMiddleware)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cart := models.Cart{
		UserId:   int64(0),
		AdvertId: 2,
		Amount:   8,
	}

	ad := models.Advert{
		Id:     int64(2),
		Name:   "aboba",
		Amount: 10,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	cu.On("GetOrderFromCart", cart.UserId, cart.AdvertId).Return(&cart, nil)
	au.On("GetAdvert", ad.Id, int64(0), false).Return(&ad, nil)
	uu.On("GetById", cart.UserId).Return(&profile, nil)
	cu.On("MakeOrder", &cart, &ad).Return(myerr.InternalError)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cart/2/checkout", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}
