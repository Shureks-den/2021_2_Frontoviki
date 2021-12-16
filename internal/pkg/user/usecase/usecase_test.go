package usecase

import (
	"testing"
	"yula/internal/models"

	myerr "yula/internal/error"
	"yula/internal/pkg/user/mocks"

	imageloader "yula/internal/pkg/image_loader"

	imageloaderMocks "yula/internal/pkg/image_loader/mocks"

	imageloaderRepo "yula/internal/pkg/image_loader/repository"
	imageloaderUse "yula/internal/pkg/image_loader/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

var (
	ilr imageloader.ImageLoaderRepository = imageloaderRepo.NewImageLoaderRepository()
	ilu imageloader.ImageLoaderUsecase    = imageloaderUse.NewImageLoaderUsecase(ilr)
)

func TestCreate(t *testing.T) {

	ur := mocks.UserRepository{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, ilu)

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	ur.On("SelectByEmail", reqUser.Email).Return(nil, myerr.EmptyQuery).Once()
	ur.On("Insert", mock.MatchedBy(func(ud *models.UserData) bool { return ud.Email == reqUser.Email })).Return(nil).Once()

	createdUser, error := uu.Create(&reqUser)
	assert.Nil(t, error)

	assert.Equal(t, reqUser.Email, createdUser.Email)
	assert.NotEqual(t, reqUser.Password, createdUser.Password)
}

func TestGetByEmail(t *testing.T) {
	ur := mocks.UserRepository{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, ilu)

	reqUser := &models.UserSignUp{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	user := &models.UserData{
		Id:       0,
		Password: "aboba",
		Email:    reqUser.Email,
	}

	ur.On("SelectByEmail", reqUser.Email).Return(user, nil)
	userRes, error := uu.GetByEmail(reqUser.Email)
	assert.Nil(t, error)

	assert.Equal(t, userRes.Email, user.Email)
}

func TestTwiceCreate(t *testing.T) {
	ur := mocks.UserRepository{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, ilu)

	reqUser := &models.UserSignUp{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	ur.On("SelectByEmail", reqUser.Email).Return(nil, myerr.EmptyQuery).Once()
	ur.On("Insert", mock.MatchedBy(func(ud *models.UserData) bool { return ud.Email == reqUser.Email })).Return(nil).Once()

	createdUser, error := uu.Create(reqUser)
	assert.Nil(t, error)

	ur.On("SelectByEmail", reqUser.Email).Return(createdUser, nil)
	usr, error := uu.Create(reqUser)

	assert.Equal(t, error, myerr.AlreadyExist)
	assert.Nil(t, usr)
}

func TestGetByEmailUserNotExist(t *testing.T) {
	ur := mocks.UserRepository{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, ilu)

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	ur.On("SelectByEmail", reqUser.Email).Return(nil, myerr.EmptyQuery)

	_, error := uu.GetByEmail(reqUser.Email)
	assert.Equal(t, error, myerr.NotExist)
}

func TestCheckPassword(t *testing.T) {
	ur := mocks.UserRepository{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, ilu)

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	ur.On("SelectByEmail", reqUser.Email).Return(nil, myerr.EmptyQuery).Once()
	ur.On("Insert", mock.MatchedBy(func(ud *models.UserData) bool { return ud.Email == reqUser.Email })).Return(nil).Once()

	createdUser, error := uu.Create(&reqUser)
	assert.Nil(t, error)

	error = uu.CheckPassword(createdUser, reqUser.Password)
	assert.Equal(t, error, nil)
}

func TestCheckPasswordInvalid(t *testing.T) {
	ur := mocks.UserRepository{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, ilu)

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	ur.On("SelectByEmail", reqUser.Email).Return(nil, myerr.EmptyQuery).Once()
	ur.On("Insert", mock.MatchedBy(func(ud *models.UserData) bool { return ud.Email == reqUser.Email })).Return(nil).Once()

	createdUser, error := uu.Create(&reqUser)
	assert.Nil(t, error)

	error = uu.CheckPassword(createdUser, reqUser.Password+"aboba")
	assert.Equal(t, error, myerr.PasswordMismatch)
}

func TestGetById(t *testing.T) {
	ur := mocks.UserRepository{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, ilu)

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	ur.On("SelectByEmail", reqUser.Email).Return(nil, myerr.EmptyQuery).Once()
	ur.On("Insert", mock.MatchedBy(func(ud *models.UserData) bool { return ud.Email == reqUser.Email })).Return(nil).Once()

	createdUser, error := uu.Create(&reqUser)
	assert.Nil(t, error)

	ur.On("SelectById", createdUser.Id).Return(createdUser, nil)
	user, error := uu.GetById(createdUser.Id)
	assert.Nil(t, error)

	assert.Equal(t, user, createdUser.ToProfile())
}

func TestGetByIdUserNotExist(t *testing.T) {
	ur := mocks.UserRepository{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, ilu)

	ur.On("SelectById", mock.MatchedBy(func(userId int64) bool { return userId < 0 })).Return(nil, myerr.EmptyQuery)

	_, error := uu.GetById(-1)
	assert.Equal(t, error, myerr.NotExist)
}

func TestUpdateUserProfile(t *testing.T) {
	ur := mocks.UserRepository{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, ilu)

	reqUser := models.UserData{
		Id:       0,
		Password: "password",
		Email:    "superchel@shibanov.jp",
		Name:     "aboba",
	}

	ur.On("SelectById", reqUser.Id).Return(&reqUser, nil)
	ur.On("SelectByEmail", reqUser.Email).Return(nil, myerr.NotExist)
	ur.On("Update", &reqUser).Return(nil)

	newProfile, error := uu.UpdateProfile(reqUser.Id, &reqUser)
	assert.Nil(t, error)

	assert.Equal(t, newProfile, reqUser.ToProfile())
}

func TestUpdateUserAlreadyExist(t *testing.T) {
	ur := mocks.UserRepository{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, ilu)

	userActual := models.UserData{
		Id:       0,
		Password: "password",
		Email:    "superchel@shibanov.jp",
		Name:     "aboba",
	}

	userNew := models.UserData{
		Id:       0,
		Password: "password",
		Email:    "aaaaa@shibanov.jp",
		Name:     "baobab",
	}

	userOther := models.UserData{
		Id:       150,
		Password: "kabalfmbfal",
		Email:    "aaaaa@shibanov.jp",
		Name:     "kgrmwgwmgklwg",
	}

	ur.On("SelectById", userNew.Id).Return(&userActual, nil)
	ur.On("SelectByEmail", userNew.Email).Return(&userOther, nil)

	newProfile, error := uu.UpdateProfile(userNew.Id, &userNew)
	assert.Equal(t, error, myerr.AlreadyExist)

	assert.Nil(t, newProfile)
}

func TestUpdatePasswordSuccess(t *testing.T) {
	ur := mocks.UserRepository{}
	mockedILU := imageloaderMocks.ImageLoaderUsecase{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, &mockedILU)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	user := models.UserData{
		Id:       0,
		Password: string(passwordHash),
		Email:    "superchel@shibanov.jp",
		Name:     "aboba",
		Image:    "not default",
	}

	assert.Nil(t, err)

	cp := models.ChangePassword{
		Email:       user.Email,
		Password:    "password",
		NewPassword: "newpassword",
	}

	ur.On("SelectById", user.Id).Return(&user, nil)
	ur.On("Update", &user).Return(nil)

	error := uu.UpdatePassword(user.Id, &cp)
	assert.Nil(t, error)
}

func TestSetRatingOk(t *testing.T) {
	ur := mocks.UserRepository{}
	mockedILU := imageloaderMocks.ImageLoaderUsecase{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, &mockedILU)

	rating := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 0}
	selrat := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 0}

	rr.On("SelectRating", rating.UserFrom, rating.UserTo).Return(selrat, myerr.EmptyQuery)

	err := uu.SetRating(rating)
	assert.NoError(t, err)
}

func TestSetRatingOk2(t *testing.T) {
	ur := mocks.UserRepository{}
	mockedILU := imageloaderMocks.ImageLoaderUsecase{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, &mockedILU)

	rating := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 3}
	selrat := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 0}

	rr.On("SelectRating", rating.UserFrom, rating.UserTo).Return(selrat, myerr.EmptyQuery)
	rr.On("InsertRating", rating).Return(nil)
	rr.On("UpdateStat", rating.UserTo, rating.Rating, 1).Return(nil)

	err := uu.SetRating(rating)
	assert.NoError(t, err)
}

func TestSetRatingOk3(t *testing.T) {
	ur := mocks.UserRepository{}
	mockedILU := imageloaderMocks.ImageLoaderUsecase{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, &mockedILU)

	rating := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 0}
	selrat := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 0}

	rr.On("SelectRating", rating.UserFrom, rating.UserTo).Return(selrat, nil)
	rr.On("DeleteRating", rating).Return(nil)
	rr.On("UpdateStat", rating.UserTo, rating.Rating, -1).Return(nil)

	err := uu.SetRating(rating)
	assert.NoError(t, err)
}

func TestSetRatingOk4(t *testing.T) {
	ur := mocks.UserRepository{}
	mockedILU := imageloaderMocks.ImageLoaderUsecase{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, &mockedILU)

	rating := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 3}
	selrat := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 0}

	rr.On("SelectRating", rating.UserFrom, rating.UserTo).Return(selrat, nil)
	rr.On("UpdateRating", rating).Return(nil)
	rr.On("UpdateStat", rating.UserTo, rating.Rating, 0).Return(nil)

	err := uu.SetRating(rating)
	assert.NoError(t, err)
}

func TestSetRatingError(t *testing.T) {
	ur := mocks.UserRepository{}
	mockedILU := imageloaderMocks.ImageLoaderUsecase{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, &mockedILU)

	rating := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 0}
	selrat := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 0}

	rr.On("SelectRating", rating.UserFrom, rating.UserTo).Return(selrat, myerr.InternalError)

	err := uu.SetRating(rating)
	assert.Error(t, err)
}

func TestSetRatingError2(t *testing.T) {
	ur := mocks.UserRepository{}
	mockedILU := imageloaderMocks.ImageLoaderUsecase{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, &mockedILU)

	rating := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 3}
	selrat := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 0}

	rr.On("SelectRating", rating.UserFrom, rating.UserTo).Return(selrat, myerr.EmptyQuery)
	rr.On("InsertRating", rating).Return(myerr.InternalError)

	err := uu.SetRating(rating)
	assert.Error(t, err)
}

func TestGetRatingOk(t *testing.T) {
	ur := mocks.UserRepository{}
	mockedILU := imageloaderMocks.ImageLoaderUsecase{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, &mockedILU)

	rating := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 3}
	selrat := &models.RatingStat{RatingSum: 8, RatingCount: 2, RatingAvg: 4.0, PersonalRate: 3, IsRated: true}

	rr.On("SelectStat", rating.UserTo).Return(selrat.RatingSum, selrat.RatingCount, nil)
	rr.On("SelectRating", rating.UserFrom, rating.UserTo).Return(rating, myerr.EmptyQuery)

	_, err := uu.GetRating(rating.UserFrom, rating.UserTo)
	assert.NoError(t, err)
}

func TestGetRatingOk2(t *testing.T) {
	ur := mocks.UserRepository{}
	mockedILU := imageloaderMocks.ImageLoaderUsecase{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, &mockedILU)

	rating := &models.Rating{UserFrom: 1, UserTo: 2, Rating: 3}
	selrat := &models.RatingStat{RatingSum: 8, RatingCount: 2, RatingAvg: 4.0, PersonalRate: 3, IsRated: true}

	rr.On("SelectStat", rating.UserTo).Return(selrat.RatingSum, selrat.RatingCount, nil)
	rr.On("SelectRating", rating.UserFrom, rating.UserTo).Return(rating, nil)

	_, err := uu.GetRating(rating.UserFrom, rating.UserTo)
	assert.NoError(t, err)
}

func TestGetRatingError(t *testing.T) {
	ur := mocks.UserRepository{}
	mockedILU := imageloaderMocks.ImageLoaderUsecase{}
	rr := mocks.RatingRepository{}
	uu := NewUserUsecase(&ur, &rr, &mockedILU)

	rating := &models.Rating{UserFrom: 0, UserTo: 1, Rating: 3}
	selrat := &models.RatingStat{RatingSum: 8, RatingCount: 2, RatingAvg: 4.0, PersonalRate: 3, IsRated: true}

	rr.On("SelectStat", rating.UserTo).Return(selrat.RatingSum, selrat.RatingCount, nil)
	rr.On("SelectRating", rating.UserFrom, rating.UserTo).Return(rating, myerr.InternalError)

	_, err := uu.GetRating(rating.UserFrom, rating.UserTo)
	assert.Error(t, err)
}
