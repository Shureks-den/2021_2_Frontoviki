package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	myerror "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/search/mocks"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func ParseTime() time.Time {
	testime := "2014-11-12 11:45:26.371"
	layout := "2006-01-02 15:04:05.000"
	te, _ := time.Parse(layout, testime)
	return te
}

func TestSearchHandlerOk(t *testing.T) {
	su := mocks.SearchUsecase{}
	sh := NewSearchHandler(&su)

	router := mux.NewRouter().PathPrefix("/search").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("", http.HandlerFunc(sh.SearchHandler)).Methods(http.MethodGet, http.MethodOptions)

	sf := &models.SearchFilter{
		Query: "query", Date: time.Time{}, TimeDuration: -1, Latitude: -80, Longitude: 80, Radius: -1, SortingDate: false, SortingName: false,
	}
	page := &models.Page{PageNum: 0, Count: 50}

	srv := httptest.NewServer(router)
	defer srv.Close()

	ads := make([]*models.Advert, 0)
	su.On("SearchWithFilter", sf, page).Return(ads, nil)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/search?query=%s", srv.URL, sf.Query), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "adverts found successfully")
}

func TestSearchHandlerError1(t *testing.T) {
	su := mocks.SearchUsecase{}
	sh := NewSearchHandler(&su)

	router := mux.NewRouter().PathPrefix("/search").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("", http.HandlerFunc(sh.SearchHandler)).Methods(http.MethodGet, http.MethodOptions)

	sf := &models.SearchFilter{
		Query: "query", Date: time.Time{}, TimeDuration: -1, Latitude: -80, Longitude: 80, Radius: -1, SortingDate: false, SortingName: false,
	}
	page := &models.Page{PageNum: 0, Count: 50}

	srv := httptest.NewServer(router)
	defer srv.Close()

	// ads := make([]*models.Advert, 0)
	su.On("SearchWithFilter", sf, page).Return(nil, myerror.InternalError)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/search?query=%s", srv.URL, sf.Query), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
}

func TestSearchHandlerError2(t *testing.T) {
	su := mocks.SearchUsecase{}
	sh := NewSearchHandler(&su)

	router := mux.NewRouter().PathPrefix("/search").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("", http.HandlerFunc(sh.SearchHandler)).Methods(http.MethodGet, http.MethodOptions)

	sf := &models.SearchFilter{
		Query: "query", Date: time.Time{}, TimeDuration: -1, Latitude: -80, Longitude: 80, Radius: -1, SortingDate: false, SortingName: false,
	}
	// page := &models.Page{PageNum: 0, Count: 50}

	srv := httptest.NewServer(router)
	defer srv.Close()

	// ads := make([]*models.Advert, 0)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/search?quer=%s", srv.URL, sf.Query), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 400)
}

func TestSearchHandlerError3(t *testing.T) {
	su := mocks.SearchUsecase{}
	sh := NewSearchHandler(&su)

	router := mux.NewRouter().PathPrefix("/search").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("", http.HandlerFunc(sh.SearchHandler)).Methods(http.MethodGet, http.MethodOptions)

	sf := &models.SearchFilter{
		Query: "query", Date: time.Time{}, TimeDuration: -1, Latitude: -80, Longitude: 80, Radius: -1, SortingDate: false, SortingName: false,
	}
	// page := &models.Page{PageNum: 0, Count: 50}

	srv := httptest.NewServer(router)
	defer srv.Close()

	// ads := make([]*models.Advert, 0)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/search?query=%s&time_duration=100", srv.URL, sf.Query), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 400)
}
