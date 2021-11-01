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

	advts, err := au.GetAdvert(ad.Id)
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
	ar.On("EditImages", ad.Id, []string{"/home/aboba/"}).Return(nil)
	ilu.On("RemoveAdvertImages", ad.Images).Return(nil)

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
	ar.On("EditImages", ad.Id, []string{"/home/aboba/"}).Return(myerr.DatabaseError)

	advt, err := au.UploadImages(files, ad.Id, 1)
	assert.Nil(t, advt)
	assert.NotNil(t, err)
}

func TestUploadImagesFail5(t *testing.T) {
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
	ar.On("EditImages", ad.Id, []string{"/home/aboba/"}).Return(nil)
	illu.On("RemoveAdvertImages", ad.Images).Return(myerr.InternalError)

	advt, err := au.UploadImages(files, ad.Id, 1)
	assert.Nil(t, advt)
	assert.NotNil(t, err)
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
