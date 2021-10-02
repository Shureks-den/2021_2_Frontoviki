package usecase

import (
	"log"
	"yula/internal/models"
	"yula/internal/pkg/advt"
)

type AdvtUsecase struct {
	advtRepository advt.AdvtRepository
}

func NewAdvtUsecase(advtRepository advt.AdvtRepository) advt.AdvtUsecase {
	return &AdvtUsecase{
		advtRepository: advtRepository,
	}
}

func (au *AdvtUsecase) GetListAdvt(from int64, count int64, newest bool) ([]*models.Advert, error) {
	advts, err := au.advtRepository.SelectListAdvt(newest, from, count)
	if err != nil {
		log.Println("invalid data from SelectListAdvt")
		return nil, err
	}

	if len(advts) == 0 {
		return []*models.Advert{}, nil
	}

	return advts, nil
}
