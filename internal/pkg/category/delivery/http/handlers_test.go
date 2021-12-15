package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	myerr "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/middleware"
	"yula/internal/services/category/mocks"
	"yula/proto/generated/category"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoriesListHandlerOk(t *testing.T) {
	cc := mocks.CategoryClient{}
	ch := NewCategoryHandler(&cc)

	router := mux.NewRouter().PathPrefix("/category").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.HandleFunc("", middleware.SetSCRFToken(http.HandlerFunc(ch.CategoriesListHandler))).Methods(http.MethodGet, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	protocategories := &category.Categories{Categories: []*category.XCategory{&category.XCategory{Name: "qwerty"}}}
	cc.On("GetCategories", mock.Anything, &category.Nothing{Dummy: true}).Return(protocategories, nil)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/category", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&answer)
	assert.NoError(t, err)
}

func TestCategoriesListHandlerError(t *testing.T) {
	cc := mocks.CategoryClient{}
	ch := NewCategoryHandler(&cc)

	router := mux.NewRouter().PathPrefix("/category").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	router.HandleFunc("", middleware.SetSCRFToken(http.HandlerFunc(ch.CategoriesListHandler))).Methods(http.MethodGet, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	cc.On("GetCategories", mock.Anything, &category.Nothing{Dummy: true}).Return(nil, myerr.InternalError)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/category", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&answer)
	assert.NoError(t, err)

	assert.Equal(t, answer.Code, 500)
}
