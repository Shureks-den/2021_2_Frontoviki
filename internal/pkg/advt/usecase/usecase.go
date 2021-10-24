package usecase

import (
	"mime/multipart"
	"time"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/advt"
	imageloader "yula/internal/pkg/image_loader"
)

type AdvtUsecase struct {
	advtRepository     advt.AdvtRepository
	imageLoaderUsecase imageloader.ImageLoaderUsecase
}

func NewAdvtUsecase(advtRepository advt.AdvtRepository, imageLoaderUsecase imageloader.ImageLoaderUsecase) advt.AdvtUsecase {
	return &AdvtUsecase{
		advtRepository:     advtRepository,
		imageLoaderUsecase: imageLoaderUsecase,
	}
}

func (au *AdvtUsecase) GetListAdvt(from int64, count int64, newest bool) ([]*models.Advert, error) {
	advts, err := au.advtRepository.SelectListAdvt(newest, from, count)
	if err != nil {
		return nil, err
	}

	if len(advts) == 0 {
		return []*models.Advert{}, nil
	}

	return advts, nil
}

func (au *AdvtUsecase) CreateAdvert(userId int64, advert *models.Advert) error {
	advert.PublisherId = userId
	if advert.Amount == 0 {
		advert.Amount = 1
	}
	err := au.advtRepository.Insert(advert)
	if err != nil {
		return err
	}
	return nil
}

func (au *AdvtUsecase) GetAdvert(advertId int64) (*models.Advert, error) {
	advert, err := au.advtRepository.SelectById(advertId)
	if err != nil {
		return nil, err
	}

	if len(advert.Images) == 0 {
		advert.Images = append(advert.Images, imageloader.DefaultAdvertImage)
	}

	// инкрементировать просмотр каждый раз когда кто-то смотрит ???

	return advert, nil
}

func (au *AdvtUsecase) UpdateAdvert(advertId int64, newAdvert *models.Advert) error {
	newAdvert.Id = advertId
	err := au.advtRepository.Update(newAdvert)
	if err != nil {
		return err
	}
	return nil
}

func (au *AdvtUsecase) DeleteAdvert(advertId int64, userId int64) error {
	advert, err := au.GetAdvert(advertId)
	if err != nil {
		return err
	}

	if advert.PublisherId != userId {
		return internalError.Conflict
	}

	err = au.advtRepository.Delete(advertId)
	if err != nil {
		return err
	}

	return nil
}

func (au *AdvtUsecase) CloseAdvert(advertId int64, userId int64) error {
	advert, err := au.GetAdvert(advertId)
	if err != nil {
		return err
	}

	if advert.PublisherId != userId {
		return internalError.Conflict
	}

	advert.IsActive = false
	advert.DateClose = time.Now()

	err = au.advtRepository.Update(advert)
	if err != nil {
		return err
	}

	return nil
}

func (au *AdvtUsecase) UploadImages(files []*multipart.FileHeader, advertId int64, userId int64) (*models.Advert, error) {
	advert, err := au.advtRepository.SelectById(advertId)
	if err != nil {
		return nil, err
	}

	if advert.PublisherId != userId {
		return nil, internalError.Conflict
	}

	imageUrls, err := au.imageLoaderUsecase.UploadAdvertImages(files)
	if err != nil {
		return nil, err
	}

	oldImages := advert.Images
	err = au.advtRepository.EditImages(advertId, imageUrls)
	if err != nil {
		return nil, err
	}
	advert.Images = imageUrls

	err = au.imageLoaderUsecase.RemoveAdvertImages(oldImages)
	if err != nil {
		return nil, err
	}

	return advert, nil
}

func (au *AdvtUsecase) GetAdvertListByPublicherId(publisherId int64, is_active bool, page *models.Page) ([]*models.Advert, error) {
	adverts, err := au.advtRepository.SelectAdvertsByPublisherId(publisherId, is_active, page.PageNum, page.Count)
	if err != nil {
		return nil, err
	}

	if len(adverts) == 0 {
		return []*models.Advert{}, nil
	}

	return adverts, nil
}

func (au *AdvtUsecase) AdvertsToShort(adverts []*models.Advert) []*models.AdvertShort {
	advertsShort := []*models.AdvertShort{}
	for _, advert := range adverts {
		advertsShort = append(advertsShort, advert.ToShort())
	}
	return advertsShort
}
