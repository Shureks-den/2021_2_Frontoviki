package usecase

import (
	"testing"
	"yula/internal/models"

	myerr "yula/internal/error"
	"yula/internal/pkg/user/mocks"

	imageloader "yula/internal/pkg/image_loader"

	imageloaderRepo "yula/internal/pkg/image_loader/repository"
	imageloaderUse "yula/internal/pkg/image_loader/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	ilr imageloader.ImageLoaderRepository = imageloaderRepo.NewImageLoaderRepository()
	ilu imageloader.ImageLoaderUsecase    = imageloaderUse.NewImageLoaderUsecase(ilr)
)

func TestCreate(t *testing.T) {

	ur := mocks.UserRepository{}
	uu := NewUserUsecase(&ur, ilu)

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
	uu := NewUserUsecase(&ur, ilu)

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	ur.On("SelectByEmail", reqUser.Email).Return(nil, myerr.EmptyQuery).Once()
	ur.On("Insert", mock.MatchedBy(func(ud *models.UserData) bool { return ud.Email == reqUser.Email })).Return(nil).Once()

	createdUser, error := uu.Create(&reqUser)
	assert.Nil(t, error)

	ur.On("SelectByEmail", reqUser.Email).Return(createdUser, nil)
	user, error := uu.GetByEmail(createdUser.Email)
	assert.Nil(t, error)

	assert.Equal(t, user.Email, createdUser.Email)
}

func TestTwiceCreate(t *testing.T) {
	ur := mocks.UserRepository{}
	uu := NewUserUsecase(&ur, ilu)

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	ur.On("SelectByEmail", reqUser.Email).Return(nil, myerr.EmptyQuery).Once()
	ur.On("Insert", mock.MatchedBy(func(ud *models.UserData) bool { return ud.Email == reqUser.Email })).Return(nil).Once()

	createdUser, error := uu.Create(&reqUser)
	assert.Nil(t, error)

	ur.On("SelectByEmail", reqUser.Email).Return(createdUser, nil)
	usr, error := uu.Create(&reqUser)

	assert.Equal(t, error, myerr.AlreadyExist)
	assert.Nil(t, usr)
}

func TestGetByEmailUserNotExist(t *testing.T) {
	ur := mocks.UserRepository{}
	uu := NewUserUsecase(&ur, ilu)

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
	uu := NewUserUsecase(&ur, ilu)

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
	uu := NewUserUsecase(&ur, ilu)

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
	uu := NewUserUsecase(&ur, ilu)

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
	uu := NewUserUsecase(&ur, ilu)

	ur.On("SelectById", mock.MatchedBy(func(userId int64) bool { return userId < 0 })).Return(nil, myerr.EmptyQuery)

	_, error := uu.GetById(-1)
	assert.Equal(t, error, myerr.NotExist)
}

func TestUpdateUserProfile(t *testing.T) {
	ur := mocks.UserRepository{}
	uu := NewUserUsecase(&ur, ilu)

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	ur.On("SelectByEmail", reqUser.Email).Return(nil, myerr.EmptyQuery).Once()
	ur.On("Insert", mock.MatchedBy(func(ud *models.UserData) bool { return ud.Email == reqUser.Email })).Return(nil).Once()

	createdUser, error := uu.Create(&reqUser)
	assert.Nil(t, error)

	createdUser.Email = "aboba@obama.com"

	ur.On("SelectById", createdUser.Id).Return(createdUser, nil)
	ur.On("SelectByEmail", createdUser.Email).Return(nil, myerr.NotExist).Once()
	ur.On("Update", createdUser).Return(nil).Once()

	newProfile, error := uu.UpdateProfile(createdUser.Id, createdUser)
	assert.Nil(t, error)

	assert.Equal(t, newProfile, createdUser.ToProfile())
}
