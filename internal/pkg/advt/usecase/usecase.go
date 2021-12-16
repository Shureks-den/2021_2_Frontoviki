package usecase

import (
	"log"
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
	return advts, err
}

func (au *AdvtUsecase) CreateAdvert(userId int64, advert *models.Advert) error {
	advert.PublisherId = userId
	if advert.Amount == 0 {
		advert.Amount = 1
	}
	err := au.advtRepository.Insert(advert)
	return err
}

func (au *AdvtUsecase) GetAdvert(advertId, userId int64, updateViews bool) (*models.Advert, error) {
	advert, err := au.advtRepository.SelectById(advertId)

	if err != nil {
		return nil, err
	}

	if len(advert.Images) == 0 {
		advert.Images = append(advert.Images, imageloader.DefaultAdvertImage)
	}

	// инкрементировать просмотр каждый раз когда кто-то смотрит под флагом
	// и смотрящий - не владелец объявления
	if updateViews && userId != advert.PublisherId {
		err = au.advtRepository.UpdateViews(advertId)
	}
	return advert, err
}

func (au *AdvtUsecase) UpdateAdvert(advertId int64, newAdvert *models.Advert) error {
	newAdvert.Id = advertId
	oldAdvert, err := au.advtRepository.SelectById(advertId)
	if err != nil {
		return err
	}

	newAdvert.PublishedAt = oldAdvert.PublishedAt
	newAdvert.DateClose = oldAdvert.DateClose
	newAdvert.IsActive = oldAdvert.IsActive

	if newAdvert.Price != oldAdvert.Price {
		if err = au.advtRepository.UpdatePrice(&models.AdvertPrice{
			AdvertId:   newAdvert.Id,
			Price:      int64(newAdvert.Price),
			ChangeTime: time.Now(),
		}); err != nil {
			return err
		}
	}

	err = au.advtRepository.Update(newAdvert)
	if err != nil {
		return err
	}

	newAdvert.Images = oldAdvert.Images
	newAdvert.Views = oldAdvert.Views

	return nil
}

func (au *AdvtUsecase) DeleteAdvert(advertId int64, userId int64) error {
	advert, err := au.GetAdvert(advertId, userId, false)
	if err != nil {
		return err
	}

	if advert.PublisherId != userId {
		return internalError.Conflict
	}

	err = au.advtRepository.Delete(advertId)
	return err
}

func (au *AdvtUsecase) CloseAdvert(advertId int64, userId int64) error {
	advert, err := au.GetAdvert(advertId, userId, false)
	if err != nil {
		return err
	}

	if advert.PublisherId != userId {
		return internalError.Conflict
	}

	advert.IsActive = false
	advert.DateClose = time.Now()

	err = au.advtRepository.Update(advert)
	return err
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
	err = au.advtRepository.InsertImages(advertId, imageUrls)
	if err != nil {
		return nil, err
	}
	advert.Images = append(oldImages, imageUrls...)

	// err = au.imageLoaderUsecase.RemoveAdvertImages(oldImages)
	// if err != nil {
	// 	return nil, err
	// }

	return advert, nil
}

func (au *AdvtUsecase) RemoveImages(images []string, advertId, userId int64) error {
	advert, err := au.advtRepository.SelectById(advertId)
	if err != nil {
		return err
	}

	if advert.PublisherId != userId {
		return internalError.Conflict
	}

	err = au.advtRepository.DeleteImages(images, advertId)
	if err != nil {
		return err
	}

	err = au.imageLoaderUsecase.RemoveAdvertImages(images)
	return err
}

func (au *AdvtUsecase) GetAdvertListByPublicherId(publisherId int64, is_active bool, page *models.Page) ([]*models.Advert, error) {
	adverts, err := au.advtRepository.SelectAdvertsByPublisherId(publisherId, is_active, page.PageNum, page.Count)
	return adverts, err
}

func (au *AdvtUsecase) AdvertsToShort(adverts []*models.Advert) []*models.AdvertShort {
	advertsShort := make([]*models.AdvertShort, 0, len(adverts))
	for _, advert := range adverts {
		advertsShort = append(advertsShort, advert.ToShort())
	}
	return advertsShort
}

func (au *AdvtUsecase) GetAdvertListByCategory(categoryName string, page *models.Page) ([]*models.Advert, error) {
	adverts, err := au.advtRepository.SelectAdvertsByCategory(categoryName, page.PageNum, page.Count)
	if err != nil {
		return nil, err
	}

	return adverts, nil
}

func (au *AdvtUsecase) GetFavoriteList(userId int64, page *models.Page) ([]*models.Advert, error) {
	adverts, err := au.advtRepository.SelectFavoriteAdverts(userId, page.PageNum, page.Count)
	if err == nil || err == internalError.EmptyQuery {
		return adverts, nil
	}

	return nil, err
}

func (au *AdvtUsecase) AddFavorite(userId int64, advertId int64) error {
	_, err := au.advtRepository.SelectFavorite(userId, advertId)
	switch err {
	case internalError.EmptyQuery:
		err = au.advtRepository.InsertFavorite(userId, advertId)
		return err
	}
	return err
}

func (au *AdvtUsecase) RemoveFavorite(userId int64, advertId int64) error {
	_, err := au.advtRepository.SelectFavorite(userId, advertId)
	switch err {
	case nil:
		err = au.advtRepository.DeleteFavorite(userId, advertId)
		return err
	}
	return err
}

func (au *AdvtUsecase) GetAdvertViews(advertId int64) (int64, error) {
	views, err := au.advtRepository.SelectViews(advertId)
	return views, err
}

func (au *AdvtUsecase) UpdateAdvertPrice(userId int64, adPrice *models.AdvertPrice) error {
	advert, err := au.advtRepository.SelectById(adPrice.AdvertId)
	if err != nil {
		return err
	}

	if advert.PublisherId != userId || adPrice.Price < 0 {
		return internalError.Conflict
	}

	if advert.Price != int(adPrice.Price) {
		return internalError.AlreadyExist
	}

	advert.Price = int(adPrice.Price)
	err = au.advtRepository.Update(advert)
	if err != nil {
		return err
	}

	err = au.advtRepository.UpdatePrice(adPrice)
	return err
}

func (au *AdvtUsecase) GetPriceHistory(advertId int64) ([]*models.AdvertPrice, error) {
	priceHistory, err := au.advtRepository.SelectPriceHistory(advertId)
	return priceHistory, err
}

func (au *AdvtUsecase) UpdatePromotion(userId int64, promo *models.Promotion) error {
	_, err := au.advtRepository.SelectById(promo.AdvertId)
	if err != nil {
		return err
	}

	// if userId != advert.PublisherId {
	// 	return internalError.Conflict
	// }

	if promo.PromoLevel < advt.MinPromo || promo.PromoLevel >= advt.MaxPromo {
		return internalError.BadRequest
	}

	promo.UpdateTime = time.Now()
	err = au.advtRepository.UpdatePromo(promo)
	return err
}

func (au *AdvtUsecase) GetFavoriteCount(advertId int64) (int64, error) {
	count, err := au.advtRepository.SelectFavoriteCount(advertId)
	if err == nil || err == internalError.EmptyQuery {
		return count, nil
	}
	return count, err
}

func (au *AdvtUsecase) GetRecomendations(advertId int64, count int64, userId int64) ([]*models.Advert, error) {
	// err := au.advtRepository.RegenerateRecomendations()
	// if err != nil {
	// 	log.Printf("1: %s", err.Error())
	// 	return nil, err
	// }

	adverts, err := au.advtRepository.SelectRecomendations(advertId, count, userId)
	if err != nil {
		log.Printf("2: %s", err.Error())
		return nil, err
	}

	if len(adverts) != 0 {
		return adverts, err
	}

	adverts, err = au.advtRepository.SelectDummyRecomendations(advertId, count)
	return adverts, err
}
