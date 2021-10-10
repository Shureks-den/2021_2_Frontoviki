package usecase

import (
	"mime/multipart"
	"time"
	internalError "yula/internal/error"
	"yula/internal/models"
	imageloader "yula/internal/pkg/image_loader"
	"yula/internal/pkg/user"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo       user.UserRepository
	imageLoaderUse imageloader.ImageLoaderUsecase
}

func NewUserUsecase(repo user.UserRepository, imageLoaderUse imageloader.ImageLoaderUsecase) user.UserUsecase {
	return &UserUsecase{
		userRepo:       repo,
		imageLoaderUse: imageLoaderUse,
	}
}

func (uu *UserUsecase) Create(userSU *models.UserSignUp) (*models.UserData, error) {
	if _, err := uu.GetByEmail(userSU.Email); err != internalError.NotExist {
		switch err {
		case nil:
			return nil, internalError.AlreadyExist

		default:
			return nil, err
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userSU.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, internalError.InternalError
	}

	user := models.UserData{}
	user.Email = userSU.Email
	user.Password = string(passwordHash)
	user.Name = userSU.Name
	user.Surname = userSU.Surname
	user.CreatedAt = time.Now()
	user.Image = imageloader.DefaultAvatar

	dbErr := uu.userRepo.Insert(&user)

	if dbErr != nil {
		return nil, dbErr
	}

	return &user, nil
}

func (uu *UserUsecase) GetByEmail(email string) (*models.UserData, error) {
	user, err := uu.userRepo.SelectByEmail(email)

	if err == nil {
		return user, nil
	}

	switch err {
	case internalError.EmptyQuery:
		return nil, internalError.NotExist
	default:
		return nil, internalError.InternalError
	}
}

func (uu *UserUsecase) CheckPassword(user *models.UserData, gettedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(gettedPassword))
	if err != nil {
		return internalError.PasswordMismatch
	}
	return nil
}

func (uu *UserUsecase) GetById(user_id int64) (*models.Profile, error) {
	user, err := uu.userRepo.SelectById(user_id)

	if err == nil {
		return user.ToProfile(), nil
	}

	switch err {
	case internalError.EmptyQuery:
		return nil, internalError.NotExist
	default:
		return nil, internalError.InternalError
	}
}

func (uu *UserUsecase) UpdateProfile(userId int64, userNew *models.UserData) (*models.Profile, error) {
	userActual, err := uu.userRepo.SelectById(userId)
	if err != nil {
		return nil, err
	}
	// userActual.Id != userNew.Id => error

	if userNew.Email != "" && userNew.Email != userActual.Email {
		// проверка на уникальность новой почты
		_, serverErr := uu.GetByEmail(userNew.Email)
		if serverErr != internalError.NotExist {
			return nil, serverErr
		}
	}

	userNew.Id = userId
	userNew.Password = userActual.Password
	userNew.CreatedAt = userActual.CreatedAt
	userNew.Image = userActual.Image

	err = uu.userRepo.Update(userNew)
	if err != nil {
		return nil, err
	}

	return userNew.ToProfile(), nil
}

func (uu *UserUsecase) UploadAvatar(file *multipart.FileHeader, userId int64) (*models.UserData, error) {
	user, err := uu.userRepo.SelectById(userId)
	if err != nil {
		return nil, err
	}

	// физическая загрузка фотки
	imageUrl, err := uu.imageLoaderUse.UploadAvatar(file)
	if err != nil {
		return nil, err
	}

	oldAvatar := user.Image
	if !(oldAvatar == "" || oldAvatar == imageloader.DefaultAvatar) {
		err = uu.imageLoaderUse.RemoveAvatar(oldAvatar)
		if err != nil {
			return nil, err
		}
	}

	user.Image = imageUrl
	err = uu.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
