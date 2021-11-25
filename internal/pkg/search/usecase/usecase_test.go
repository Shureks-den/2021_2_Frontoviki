package usecase

import (
	"testing"
	myerror "yula/internal/error"
	"yula/internal/models"

	mockAdvt "yula/internal/pkg/advt/mocks"
	mockSrch "yula/internal/pkg/search/mocks"

	"github.com/stretchr/testify/assert"
)

func TestSearchWithFilter(t *testing.T) {
	sr := &mockSrch.SearchRepository{}
	ar := &mockAdvt.AdvtRepository{}
	su := NewSearchUsecase(sr, ar)

	filter := &models.SearchFilter{}
	page := &models.Page{}

	sr.On("SelectWithFilter", filter, page.PageNum, page.Count).Return([]*models.Advert{}, nil)

	_, err := su.SearchWithFilter(filter, page)

	assert.NoError(t, err)
}

func TestSearchWithFilterError(t *testing.T) {
	sr := &mockSrch.SearchRepository{}
	ar := &mockAdvt.AdvtRepository{}
	su := NewSearchUsecase(sr, ar)

	filter := &models.SearchFilter{}
	page := &models.Page{}

	sr.On("SelectWithFilter", filter, page.PageNum, page.Count).Return([]*models.Advert{}, myerror.InternalError)

	_, err := su.SearchWithFilter(filter, page)

	assert.Error(t, err)
}
