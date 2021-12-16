package usecase

import (
	"mime/multipart"
	"testing"
	"yula/internal/models"

	mockAdvt "yula/internal/pkg/advt/mocks"

	myerr "yula/internal/error"

	"yula/internal/pkg/image_loader/mocks"

	"github.com/stretchr/testify/assert"
)

var (
	ilu mocks.ImageLoaderUsecase = mocks.ImageLoaderUsecase{}
)

func TestGetListAdvtSuccess(t *testing.T) {
	ua := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ua, &ilu)
	ads := []*models.Advert{
		{
			Id:     32,
			Name:   "aboba",
			Amount: 10,
		},
	}
	ua.On("SelectListAdvt", false, int64(1), int64(5)).Return(ads, nil)

	advts, error := au.GetListAdvt(1, 5, false)
	assert.Nil(t, error)
	assert.Equal(t, ads, advts)
}

func TestGetListAdvtFail(t *testing.T) {
	ua := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ua, &ilu)
	ua.On("SelectListAdvt", false, int64(1), int64(5)).Return(nil, myerr.DatabaseError)

	advts, err := au.GetListAdvt(1, 5, false)
	assert.Equal(t, err, myerr.DatabaseError)
	assert.Nil(t, advts)
}

func TestCreateAdvert(t *testing.T) {
	ua := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ua, &ilu)
	ad := models.Advert{
		Name:   "aboba",
		Amount: 0,
	}
	ua.On("Insert", &ad).Return(nil)

	err := au.CreateAdvert(int64(22), &ad)
	assert.Nil(t, err)
	assert.Equal(t, ad.PublisherId, int64(22))
	assert.Equal(t, ad.Amount, int64(1))
}

func TestGetAdvert(t *testing.T) {
	ua := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ua, &ilu)
	ad := models.Advert{
		Id:     0,
		Name:   "aboba",
		Amount: 0,
	}
	ua.On("SelectById", int64(0)).Return(&ad, nil)

	advts, err := au.GetAdvert(ad.Id, -1, false)
	assert.Nil(t, err)
	assert.Equal(t, ad, *advts)
}

func TestUpdateAdvert(t *testing.T) {
	ua := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ua, &ilu)
	oldAd := models.Advert{
		Id:     0,
		Name:   "aboba",
		Amount: 0,
	}
	newAd := models.Advert{
		Id:     0,
		Name:   "baobab",
		Amount: 0,
	}
	ua.On("SelectById", int64(0)).Return(&oldAd, nil)
	ua.On("Update", &newAd).Return(nil)

	err := au.UpdateAdvert(oldAd.Id, &newAd)
	assert.Nil(t, err)
}

func TestDeleteAdvertSuccess(t *testing.T) {
	ua := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ua, &ilu)
	ad := models.Advert{
		Id:          0,
		Name:        "aboba",
		Amount:      0,
		PublisherId: 1,
	}
	ua.On("Delete", ad.Id).Return(nil)
	ua.On("SelectById", ad.Id).Return(&ad, nil)

	err := au.DeleteAdvert(ad.Id, 1)
	assert.Nil(t, err)
}

func TestDeleteAdvertFail(t *testing.T) {
	ua := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ua, &ilu)
	ad := models.Advert{
		Id:          0,
		Name:        "aboba",
		Amount:      0,
		PublisherId: 1,
	}
	ua.On("Delete", ad.Id).Return(nil)
	ua.On("SelectById", ad.Id).Return(&ad, nil)

	err := au.DeleteAdvert(ad.Id, 10)
	assert.NotNil(t, err)
}

func TestCloseAdvertSuccess(t *testing.T) {
	ua := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ua, &ilu)
	ad := models.Advert{
		Id:          0,
		Name:        "aboba",
		Amount:      0,
		PublisherId: 1,
	}
	ua.On("SelectById", ad.Id).Return(&ad, nil)
	ua.On("Update", &ad).Return(nil)

	err := au.CloseAdvert(ad.Id, 1)
	assert.Nil(t, err)
}

func TestCloseAdvertFail(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)
	ad := models.Advert{
		Id:          0,
		Name:        "aboba",
		Amount:      0,
		PublisherId: 1,
	}
	ar.On("SelectById", ad.Id).Return(&ad, nil)
	ar.On("Update", &ad).Return(nil)

	err := au.CloseAdvert(ad.Id, 10)
	assert.NotNil(t, err)
}

func TestUploadImagesSuccess(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)
	files := []*multipart.FileHeader{
		{
			Filename: "aboba",
		},
	}
	ad := models.Advert{
		Id:          0,
		Name:        "aboba",
		Amount:      0,
		PublisherId: 1,
	}
	ar.On("SelectById", ad.Id).Return(&ad, nil)
	ilu.On("UploadAdvertImages", files).Return([]string{"/home/aboba/"}, nil)
	ar.On("InsertImages", ad.Id, []string{"/home/aboba/"}).Return(nil)
	// ilu.On("RemoveAdvertImages", ad.Images).Return(nil)

	advt, err := au.UploadImages(files, ad.Id, 1)
	assert.NotNil(t, advt)
	assert.Nil(t, err)
}

func TestUploadImagesFail1(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)
	files := []*multipart.FileHeader{
		{
			Filename: "aboba",
		},
	}
	ad := models.Advert{
		Id:          0,
		Name:        "aboba",
		Amount:      0,
		PublisherId: 1,
	}
	ar.On("SelectById", ad.Id).Return(&ad, myerr.DatabaseError)

	advt, err := au.UploadImages(files, ad.Id, 1)
	assert.Nil(t, advt)
	assert.NotNil(t, err)
}

func TestUploadImagesFail2(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)
	files := []*multipart.FileHeader{
		{
			Filename: "aboba",
		},
	}
	ad := models.Advert{
		Id:          0,
		Name:        "aboba",
		Amount:      0,
		PublisherId: 1,
	}
	ar.On("SelectById", ad.Id).Return(&ad, nil)

	advt, err := au.UploadImages(files, ad.Id, 10)
	assert.Nil(t, advt)
	assert.NotNil(t, err)
}

func TestUploadImagesFail3(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	illu := mocks.ImageLoaderUsecase{}
	au := NewAdvtUsecase(&ar, &illu)
	files := []*multipart.FileHeader{
		{
			Filename: "aboba",
		},
	}
	ad := models.Advert{
		Id:          0,
		Name:        "aboba",
		Amount:      0,
		PublisherId: 1,
	}
	ar.On("SelectById", ad.Id).Return(&ad, nil)
	illu.On("UploadAdvertImages", files).Return([]string{"/home/abobab/"}, myerr.InternalError)

	advt, err := au.UploadImages(files, ad.Id, 1)
	assert.Nil(t, advt)
	assert.NotNil(t, err)
}

func TestUploadImagesFail4(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	illu := mocks.ImageLoaderUsecase{}
	au := NewAdvtUsecase(&ar, &illu)
	files := []*multipart.FileHeader{
		{
			Filename: "aboba",
		},
	}
	ad := models.Advert{
		Id:          0,
		Name:        "aboba",
		Amount:      0,
		PublisherId: 1,
	}
	ar.On("SelectById", ad.Id).Return(&ad, nil)
	illu.On("UploadAdvertImages", files).Return([]string{"/home/aboba/"}, nil)
	ar.On("InsertImages", ad.Id, []string{"/home/aboba/"}).Return(myerr.DatabaseError)

	advt, err := au.UploadImages(files, ad.Id, 1)
	assert.Nil(t, advt)
	assert.NotNil(t, err)
}

func TestRemoveImagesSuccess(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	illu := mocks.ImageLoaderUsecase{}
	au := NewAdvtUsecase(&ar, &illu)

	ad := &models.Advert{
		Id:          0,
		Name:        "aboba",
		Amount:      0,
		PublisherId: 1,
	}

	images := []string{"img1", "img2"}
	ar.On("SelectById", ad.Id).Return(ad, nil)
	ar.On("DeleteImages", images, ad.Id).Return(nil)
	illu.On("RemoveAdvertImages", images).Return(nil)

	err := au.RemoveImages(images, ad.Id, ad.PublisherId)
	assert.NoError(t, err)
}

func TestRemoveImagesError(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	illu := mocks.ImageLoaderUsecase{}
	au := NewAdvtUsecase(&ar, &illu)

	ad := &models.Advert{
		Id:          0,
		Name:        "aboba",
		Amount:      0,
		PublisherId: 1,
	}

	images := []string{"img1", "img2"}
	ar.On("SelectById", ad.Id).Return(ad, nil)
	ar.On("DeleteImages", images, ad.Id).Return(nil)
	illu.On("RemoveAdvertImages", images).Return(nil)

	err := au.RemoveImages(images, ad.Id, int64(0))
	assert.Error(t, err)
}

func TestGetAdvertListByPublicherId(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}

	au := NewAdvtUsecase(&ar, &ilu)

	ar.On("SelectAdvertsByPublisherId", int64(0), false, int64(0), int64(0)).Return(nil, nil)

	advts, err := au.GetAdvertListByPublicherId(0, false, &models.Page{})
	assert.Nil(t, advts)
	assert.Nil(t, err)
}

func TestAdvertsToShort(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}

	au := NewAdvtUsecase(&ar, &ilu)

	advts := []*models.Advert{
		{
			Name: "aboba",
		},
	}

	au.AdvertsToShort(advts)
}

func TestGetAdvertListByCategorySuccess(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	advts := []*models.Advert{
		{
			Name:     "aboba",
			Category: "baobab",
		},
	}
	page := &models.Page{}

	ar.On("SelectAdvertsByCategory", advts[0].Category, page.PageNum, page.Count).Return(advts, nil)

	ads, err := au.GetAdvertListByCategory(advts[0].Category, page)
	assert.Equal(t, advts, ads)
	assert.NoError(t, err)
}

func TestGetAdvertListByCategoryError(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	advts := []*models.Advert{
		{
			Name:     "aboba",
			Category: "baobab",
		},
	}
	page := &models.Page{}

	ar.On("SelectAdvertsByCategory", advts[0].Category, page.PageNum, page.Count).Return(nil, myerr.InternalError)

	ads, err := au.GetAdvertListByCategory(advts[0].Category, page)
	assert.NotEqual(t, advts, ads)
	assert.Error(t, err)
}

func TestGetFavoriteListSuccess(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	userId := int64(1)
	page := &models.Page{}
	advts := []*models.Advert{}

	ar.On("SelectFavoriteAdverts", userId, page.PageNum, page.Count).Return(advts, nil)

	ads, err := au.GetFavoriteList(userId, page)
	assert.Equal(t, advts, ads)
	assert.NoError(t, err)
}

func TestGetFavoriteListError(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	userId := int64(1)
	page := &models.Page{}
	advts := []*models.Advert{}

	ar.On("SelectFavoriteAdverts", userId, page.PageNum, page.Count).Return(nil, myerr.InternalError)

	ads, err := au.GetFavoriteList(userId, page)
	assert.NotEqual(t, advts, ads)
	assert.Error(t, err)
}

func TestAddFavoriteSuccess(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	userId, advertId := int64(0), int64(0)

	ar.On("SelectFavorite", userId, advertId).Return(nil, myerr.EmptyQuery)
	ar.On("InsertFavorite", userId, advertId).Return(nil)

	err := au.AddFavorite(userId, advertId)
	assert.NoError(t, err)
}

func TestAddFavoriteError(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	userId, advertId := int64(0), int64(0)

	ar.On("SelectFavorite", userId, advertId).Return(nil, myerr.InternalError)

	err := au.AddFavorite(userId, advertId)
	assert.Error(t, err)
}

func TestRemoveFavoriteSuccess(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	userId, advertId := int64(0), int64(0)

	adv := &models.Advert{}

	ar.On("SelectFavorite", userId, advertId).Return(adv, nil)
	ar.On("DeleteFavorite", userId, advertId).Return(nil)

	err := au.RemoveFavorite(userId, advertId)
	assert.NoError(t, err)
}

func TestRemoveFavoriteError(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	userId, advertId := int64(0), int64(0)

	ar.On("SelectFavorite", userId, advertId).Return(nil, myerr.InternalError)

	err := au.RemoveFavorite(userId, advertId)
	assert.Error(t, err)
}

func TestGetAdvertViews(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	advertId := int64(0)

	ar.On("SelectViews", advertId).Return(int64(0), nil)

	_, err := au.GetAdvertViews(advertId)
	assert.NoError(t, err)
}

func TestUpdateAdvertPriceOk(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	adPrice := models.AdvertPrice{AdvertId: 0, Price: 100}
	adv := &models.Advert{Id: 0, Price: 100}

	ar.On("SelectById", adPrice.AdvertId).Return(adv, nil)
	ar.On("Update", adv).Return(nil)
	ar.On("UpdatePrice", &adPrice).Return(nil)

	err := au.UpdateAdvertPrice(int64(0), &adPrice)
	assert.NoError(t, err)
}

func TestUpdateAdvertPriceSuccess(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	adPrice := models.AdvertPrice{AdvertId: 0, Price: 100}
	adv := &models.Advert{Id: 0, Price: 100}

	ar.On("SelectById", adPrice.AdvertId).Return(adv, nil)
	ar.On("Update", adv).Return(nil)
	ar.On("UpdatePrice", &adPrice).Return(nil)

	err := au.UpdateAdvertPrice(int64(0), &adPrice)
	assert.NoError(t, err)
}

func TestUpdateAdvertPriceError1(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	adPrice := models.AdvertPrice{AdvertId: 0, Price: 100}

	ar.On("SelectById", adPrice.AdvertId).Return(nil, myerr.InternalError)

	err := au.UpdateAdvertPrice(int64(0), &adPrice)
	assert.Error(t, err)
}

func TestUpdateAdvertPriceError2(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	adPrice := models.AdvertPrice{AdvertId: 0, Price: -100}
	adv := &models.Advert{Id: 0, Price: 100}

	ar.On("SelectById", adPrice.AdvertId).Return(adv, nil)

	err := au.UpdateAdvertPrice(int64(0), &adPrice)
	assert.Error(t, err)
}

func TestUpdateAdvertPriceError3(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	adPrice := models.AdvertPrice{AdvertId: 0, Price: 50}
	adv := &models.Advert{Id: 0, Price: 100}

	ar.On("SelectById", adPrice.AdvertId).Return(adv, nil)

	err := au.UpdateAdvertPrice(int64(0), &adPrice)
	assert.Error(t, err)
}

func TestUpdateAdvertPriceError4(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	adPrice := models.AdvertPrice{AdvertId: 0, Price: 100}
	adv := &models.Advert{Id: 0, Price: 100}

	ar.On("SelectById", adPrice.AdvertId).Return(adv, nil)
	ar.On("Update", adv).Return(myerr.InternalError)

	err := au.UpdateAdvertPrice(int64(0), &adPrice)
	assert.Error(t, err)
}

func TestGetPriceHistorySuccess(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	adPrice := &models.AdvertPrice{AdvertId: 0, Price: 100}
	history := []*models.AdvertPrice{adPrice}

	ar.On("SelectPriceHistory", adPrice.AdvertId).Return(history, nil)

	_, err := au.GetPriceHistory(adPrice.AdvertId)
	assert.NoError(t, err)
}

func TestUpdatePromotionSuccess(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	promo := &models.Promotion{AdvertId: 0, PromoLevel: 1}
	adv := &models.Advert{Id: 0, Price: 100}

	ar.On("SelectById", promo.AdvertId).Return(adv, nil)
	ar.On("UpdatePromo", promo).Return(nil)

	err := au.UpdatePromotion(int64(0), promo)
	assert.NoError(t, err)
}

func TestUpdatePromotionError1(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	promo := &models.Promotion{AdvertId: 0, PromoLevel: 1}

	ar.On("SelectById", promo.AdvertId).Return(nil, myerr.InternalError)

	err := au.UpdatePromotion(int64(0), promo)
	assert.Error(t, err)
}

func TestUpdatePromotionError2(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	promo := &models.Promotion{AdvertId: 0, PromoLevel: -1}
	adv := &models.Advert{Id: 0, Price: 100}

	ar.On("SelectById", promo.AdvertId).Return(adv, nil)

	err := au.UpdatePromotion(int64(0), promo)
	assert.Error(t, err)
}

func TestGetFavoriteCountSuccess(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	advertId := int64(0)
	count := int64(1)

	ar.On("SelectFavoriteCount", advertId).Return(count, nil)

	_, err := au.GetFavoriteCount(advertId)
	assert.NoError(t, err)
}

func TestGetFavoriteCountError(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	advertId := int64(0)

	ar.On("SelectFavoriteCount", advertId).Return(int64(0), myerr.InternalError)

	_, err := au.GetFavoriteCount(advertId)
	assert.Error(t, err)
}

func TestGetRecomendationsOk(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	advert := &models.Advert{Id: int64(0)}
	adverts := make([]*models.Advert, 0)
	adverts = append(adverts, advert)
	count := int64(10)
	userId := int64(0)

	ar.On("SelectRecomendations", advert.Id, count, userId).Return(adverts, nil)
	ar.On("SelectDummyRecomendations", count).Return(adverts, nil)

	_, err := au.GetRecomendations(advert.Id, count, userId)
	assert.NoError(t, err)
}

func TestGetRecomendationsError1(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	advert := &models.Advert{Id: int64(0)}
	count := int64(10)
	userId := int64(0)

	ar.On("SelectRecomendations", advert.Id, count, userId).Return(nil, myerr.InternalError)

	_, err := au.GetRecomendations(advert.Id, count, userId)
	assert.Error(t, err)
}

func TestGetRecomendationsError2(t *testing.T) {
	ar := mockAdvt.AdvtRepository{}
	au := NewAdvtUsecase(&ar, &ilu)

	advert := &models.Advert{Id: int64(0)}
	adverts := make([]*models.Advert, 0)
	count := int64(10)
	userId := int64(0)

	ar.On("SelectRecomendations", advert.Id, count, userId).Return(adverts, nil)
	ar.On("SelectDummyRecomendations", count).Return([]*models.Advert{advert}, nil)

	_, err := au.GetRecomendations(advert.Id, count, userId)
	assert.NoError(t, err)
}
