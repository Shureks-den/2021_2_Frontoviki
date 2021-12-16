package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"yula/internal/models"
	"yula/internal/pkg/middleware"

	myerr "yula/internal/error"

	advtMock "yula/internal/pkg/advt/mocks"

	userMock "yula/internal/pkg/user/mocks"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	faker "github.com/jaswdr/faker"
)

func TestCreateAdvert(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("", http.HandlerFunc(ah.CreateAdvertHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	au.On("CreateAdvert", int64(0), &ad).Return(nil)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(ad)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 201)
	assert.Equal(t, Answer.Message, "advert created successfully")
}

func TestCreateFailCreateAd(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("", http.HandlerFunc(ah.CreateAdvertHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	au.On("CreateAdvert", int64(0), &ad).Return(myerr.InternalError)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(ad)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestCreateFail(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("", http.HandlerFunc(ah.CreateAdvertHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	// assert.Equal(t, Answer.Code, 400)
	// assert.Equal(t, Answer.Message, "bad request")
}

func TestAdvertDetailSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(ah.AdvertDetailHandler)).Methods(http.MethodGet, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	history := make([]*models.AdvertPrice, 0)
	history = append(history, &models.AdvertPrice{AdvertId: ad.Id, Price: int64(ad.Price)})

	au.On("GetAdvert", ad.Id, ad.PublisherId, true).Return(&ad, nil)
	au.On("GetAdvertViews", ad.Id).Return(ad.Views, nil)
	uu.On("GetById", ad.PublisherId).Return(&profile, nil)
	uu.On("GetRating", ad.PublisherId, int64(0)).Return(&models.RatingStat{}, nil)
	au.On("GetPriceHistory", ad.Id).Return(history, nil)
	au.On("GetFavoriteCount", ad.Id).Return(int64(1), nil)

	res, err := http.Get(fmt.Sprintf("%s/adverts/2", srv.URL))
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	// assert.Equal(t, Answer.Code, 200)
	// assert.Equal(t, Answer.Message, "advert found successfully")
}

func TestAdvertDetailFailParseId(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(ah.AdvertDetailHandler)).Methods(http.MethodGet, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s/adverts/2429854528447491428842528458245813", srv.URL))
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	// assert.Equal(t, Answer.Code, 400)
	// assert.Equal(t, Answer.Message, "bad request")
}

func TestAdvertDetailFailGetAd(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(ah.AdvertDetailHandler)).Methods(http.MethodGet, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	au.On("GetAdvert", ad.Id, int64(0), true).Return(&ad, myerr.InternalError)

	res, err := http.Get(fmt.Sprintf("%s/adverts/2", srv.URL))
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	// assert.Equal(t, Answer.Code, 500)
	// assert.Equal(t, Answer.Message, "internal error")
}

func TestAdvertDetailFailGetPublisher(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(ah.AdvertDetailHandler)).Methods(http.MethodGet, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	au.On("GetAdvert", ad.Id, int64(0), true).Return(&ad, nil)
	au.On("GetAdvertViews", ad.Id).Return(ad.Views, nil)
	uu.On("GetById", ad.PublisherId).Return(&profile, myerr.InternalError)

	res, err := http.Get(fmt.Sprintf("%s/adverts/2", srv.URL))
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	// assert.Equal(t, Answer.Code, 500)
	// assert.Equal(t, Answer.Message, "internal error")
}

func TestAdUpdateSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(ah.AdvertUpdateHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	newAd := models.Advert{
		Id:          ad.Id,
		Name:        "baobab",
		Amount:      ad.Amount,
		PublisherId: ad.PublisherId,
	}

	// profile := models.Profile{
	// 	Id:        0,
	// 	Email:     "aboba@baobab.com",
	// 	CreatedAt: time.Now(),
	// 	RatingSum:    5,
	// }

	au.On("UpdateAdvert", ad.Id, &newAd).Return(nil)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(newAd)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/2", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	// assert.Equal(t, Answer.Code, 201)
	// assert.Equal(t, Answer.Message, "advert updated successfully")
}

func TestAdUpdateFailParse(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(ah.AdvertUpdateHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/2428584275427828577824285427", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	// assert.Equal(t, Answer.Code, 400)
	// assert.Equal(t, Answer.Message, "bad request")
}

func TestAdUpdateCantDecode(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(ah.AdvertUpdateHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/2", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	// assert.Equal(t, Answer.Code, 400)
	// assert.Equal(t, Answer.Message, "bad request")
}

func TestAdUpdateFailUpdateAd(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(ah.AdvertUpdateHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	newAd := models.Advert{
		Id:          ad.Id,
		Name:        "baobab",
		Amount:      ad.Amount,
		PublisherId: ad.PublisherId,
	}

	// profile := models.Profile{
	// 	Id:        0,
	// 	Email:     "aboba@baobab.com",
	// 	CreatedAt: time.Now(),
	// 	RatingSum:    5,
	// }

	au.On("UpdateAdvert", ad.Id, &newAd).Return(myerr.InternalError)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(newAd)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/2", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestDeleteAdSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(ah.DeleteAdvertHandler)).Methods(http.MethodDelete, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	au.On("DeleteAdvert", ad.Id, profile.Id).Return(nil)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/adverts/2", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "advert deleted successfully")
}

func TestDeleteAdFailParseId(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(ah.DeleteAdvertHandler)).Methods(http.MethodDelete, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/adverts/242782824881398318183193", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 400)
	assert.Equal(t, Answer.Message, "bad request")
}

func TestDeleteFail(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(ah.DeleteAdvertHandler)).Methods(http.MethodDelete, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	au.On("DeleteAdvert", ad.Id, profile.Id).Return(myerr.InternalError)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/adverts/2", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestDeleteFailParse(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}", http.HandlerFunc(ah.DeleteAdvertHandler)).Methods(http.MethodDelete, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	au.On("DeleteAdvert", ad.Id, profile.Id).Return(myerr.InternalError)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/adverts/29858252288582858285828585284888284882", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 400)
	assert.Equal(t, Answer.Message, "bad request")
}

func TestCloseAdSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}/close", http.HandlerFunc(ah.CloseAdvertHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	au.On("CloseAdvert", ad.Id, profile.Id).Return(nil)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/2/close", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "advert closed successfully")
}

func TestCloseAdFailParseId(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}/close", http.HandlerFunc(ah.CloseAdvertHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/24783887781771817781717713813/close", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 400)
	assert.Equal(t, Answer.Message, "bad request")
}

func TestCloseAdFail(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}/close", http.HandlerFunc(ah.CloseAdvertHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	au.On("CloseAdvert", ad.Id, profile.Id).Return(myerr.InternalError)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/2/close", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestUploadImageSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}/upload", http.HandlerFunc(ah.UploadImageHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	faker := faker.New()
	p := faker.Person()
	fakeimg := p.Image()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	au.On("UploadImages", mock.AnythingOfType("[]*multipart.FileHeader"), ad.Id, profile.Id).Return(&ad, nil)
	uu.On("GetById", ad.PublisherId).Return(&profile, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("images", fakeimg.Name())
	assert.Nil(t, err)
	file, err := os.Open(fakeimg.Name())
	assert.Nil(t, err)
	_, err = io.Copy(fw, file)
	assert.Nil(t, err)
	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/2/upload", srv.URL), bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "images uploaded successfully")
}

func TestUploadImageFailParseId(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}/upload", http.HandlerFunc(ah.UploadImageHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/2772587782747842848283882383282/upload", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 400)
	assert.Equal(t, Answer.Message, "bad request")
}

func TestUploadImageFailUpload(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}/upload", http.HandlerFunc(ah.UploadImageHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	faker := faker.New()
	p := faker.Person()
	fakeimg := p.Image()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	au.On("UploadImages", mock.AnythingOfType("[]*multipart.FileHeader"), ad.Id, profile.Id).Return(nil, myerr.InternalError)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("images", fakeimg.Name())
	assert.Nil(t, err)
	file, err := os.Open(fakeimg.Name())
	assert.Nil(t, err)
	_, err = io.Copy(fw, file)
	assert.Nil(t, err)
	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/2/upload", srv.URL), bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestUploadImageFailGetById(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}/upload", http.HandlerFunc(ah.UploadImageHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	faker := faker.New()
	p := faker.Person()
	fakeimg := p.Image()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	au.On("UploadImages", mock.AnythingOfType("[]*multipart.FileHeader"), ad.Id, profile.Id).Return(&ad, nil)
	uu.On("GetById", ad.PublisherId).Return(nil, myerr.InternalError)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("images", fakeimg.Name())
	assert.Nil(t, err)
	file, err := os.Open(fakeimg.Name())
	assert.Nil(t, err)
	_, err = io.Copy(fw, file)
	assert.Nil(t, err)
	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/2/upload", srv.URL), bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestUploadImageFailParse(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}/upload", http.HandlerFunc(ah.UploadImageHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/2/upload", srv.URL), nil)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestRemoveImageSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/{id:[0-9]+}/images", http.HandlerFunc(ah.RemoveImageHandler)).Methods(http.MethodDelete, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	profile := models.Profile{
		Id:        0,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	au.On("RemoveImages", mock.AnythingOfType("[]string"), ad.Id, profile.Id).Return(nil)
	body := bytes.NewBuffer(nil)
	buf := models.AdvertImages{
		ImagesPath: []string{"aboba"},
	}
	err := json.NewEncoder(body).Encode(buf)
	if err != nil {
		logger.Warnf("trouble with encoder %s", err.Error())
	}

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/adverts/2/images", srv.URL), bytes.NewReader(body.Bytes()))
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)

	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "images removed successfully")
}

func TestSalesmanPageSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/salesman/{id:[0-9]+}", http.HandlerFunc(ah.SalesmanPageHandler)).Methods(http.MethodGet, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	profile := models.Profile{
		Id:        1,
		Email:     "aboba@baobab.com",
		CreatedAt: time.Now(),
	}

	uu.On("GetById", profile.Id).Return(&profile, nil)
	au.On("GetAdvertListByPublicherId", profile.Id, true, &models.Page{PageNum: 0, Count: 50}).Return([]*models.Advert{&ad}, nil)
	au.On("AdvertsToShort", []*models.Advert{&ad}).Return([]*models.AdvertShort{ad.ToShort()}, nil)
	uu.On("GetRating", int64(0), profile.Id).Return(&models.RatingStat{}, nil)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/adverts/salesman/1?page=1&count=50", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)

	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "salesman profile provided")
}

func TestArchiveSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/archive", http.HandlerFunc(ah.ArchiveHandler)).Methods(http.MethodGet, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	au.On("GetAdvertListByPublicherId", int64(0), false, &models.Page{PageNum: 0, Count: 50}).Return([]*models.Advert{&ad}, nil)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/adverts/archive?page=1&count=50", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)

	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "archive got")
}

func TestAdvertListByCategorySuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/category/aboba", http.HandlerFunc(ah.AdvertListByCategoryHandler)).Methods(http.MethodGet, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	au.On("GetAdvertListByCategory", "aboba", &models.Page{PageNum: 0, Count: 50}).Return([]*models.Advert{&ad}, nil)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/adverts/category/aboba?page=1&count=50", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)

	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "adverts got successfully")
}

func TestFavoriteListSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/favorite", http.HandlerFunc(ah.FavoriteListHandler)).Methods(http.MethodGet, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	au.On("GetFavoriteList", int64(0), &models.Page{PageNum: 0, Count: 50}).Return([]*models.Advert{&ad}, nil)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/adverts/favorite", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)

	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "favorite adverts got successfully")
}

func TestRemoveFavoriteSuccess(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/favorite/{id:[0-9]+}", http.HandlerFunc(ah.RemoveFavoriteHandler)).Methods(http.MethodDelete, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	au.On("RemoveFavorite", int64(0), int64(142)).Return(nil)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/adverts/favorite/142", srv.URL), http.NoBody)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)

	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "removed from favorite")
}

func TestUpdatePriceHistory(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/price_history", http.HandlerFunc(ah.UpdatePriceHistory)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	adPrice := &models.AdvertPrice{
		AdvertId: 1,
		Price:    1313,
	}

	b := bytes.NewBuffer(nil)
	err := json.NewEncoder(b).Encode(adPrice)
	if err != nil {
		logger.Warnf("trouble with encoder %s", err.Error())
	}

	au.On("UpdateAdvertPrice", int64(0), adPrice).Return(nil)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/adverts/price_history", srv.URL), bytes.NewReader(b.Bytes()))
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)

	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "price of 1 advert updated")
}

func TestGetPriceHistory(t *testing.T) {
	au := advtMock.AdvtUsecase{}
	uu := userMock.UserUsecase{}
	ah := NewAdvertHandler(&au, &uu)

	router := mux.NewRouter().PathPrefix("/adverts").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/price_history/{id:[0-9]+}", http.HandlerFunc(ah.GetPriceHistory)).Methods(http.MethodGet, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	ad := models.Advert{
		Id:          2,
		Name:        "aboba",
		Amount:      8,
		PublisherId: 0,
	}

	adPrice := &models.AdvertPrice{
		AdvertId: 1,
		Price:    1313,
	}

	au.On("GetPriceHistory", ad.Id).Return([]*models.AdvertPrice{adPrice}, nil)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/adverts/price_history/2", srv.URL), http.NoBody)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)

	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "favorite adverts got successfully")
}
