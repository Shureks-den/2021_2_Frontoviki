package usecase

import (
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/advt"
	"yula/internal/pkg/search"
)

type SearchUsecase struct {
	searchRepository search.SearchRepository
	advertRepository advt.AdvtRepository
}

func NewSearchUsecase(searchRepository search.SearchRepository, advertRepository advt.AdvtRepository) search.SearchUsecase {
	return &SearchUsecase{
		searchRepository: searchRepository,
		advertRepository: advertRepository,
	}
}

func (su *SearchUsecase) SearchWithFilter(filter *models.SearchFilter, page *models.Page) ([]*models.Advert, error) {
	adverts, err := su.searchRepository.SelectWithFilter(filter, page.PageNum, page.Count)
	switch err {
	case nil, internalError.EmptyQuery:
		return adverts, nil
	}
	return nil, err
}
