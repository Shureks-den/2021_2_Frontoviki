package usecase

import (
	"log"
	"mime/multipart"
	"time"
	internalError "yula/internal/error"
	"yula/internal/models"
	imageloader "yula/internal/pkg/image_loader"
	"yula/internal/pkg/user"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo             user.UserRepository
	userRatingRepository user.RatingRepository
	imageLoaderUse       imageloader.ImageLoaderUsecase
}

func NewUserUsecase(repo user.UserRepository, userRatingRepository user.RatingRepository,
	imageLoaderUse imageloader.ImageLoaderUsecase) user.UserUsecase {
	return &UserUsecase{
		userRepo:             repo,
		userRatingRepository: userRatingRepository,
		imageLoaderUse:       imageLoaderUse,
	}
}

func (uu *UserUsecase) Create(userSU *models.UserSignUp) (*models.UserData, error) {
	if _, err := uu.GetByEmail(userSU.Email); err != internalError.NotExist {
		switch err {
		case nil:
			return nil, internalError.AlreadyExist

		default:
			log.Fatal(err)

			return nil, err
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userSU.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, internalError.InternalError
	}

	user := models.UserData{}
	user.Email = userSU.Email
	user.Phone = ""
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

		switch serverErr {
		case nil:
			return nil, internalError.AlreadyExist

		case internalError.InternalError:
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
	if oldAvatar != "" && oldAvatar != imageloader.DefaultAvatar {
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

func (uu *UserUsecase) UpdatePassword(userId int64, changePassword *models.ChangePassword) error {
	user, err := uu.userRepo.SelectById(userId)
	if err != nil {
		return err
	}

	err = uu.CheckPassword(user, changePassword.Password)
	if err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(changePassword.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(passwordHash)
	err = uu.userRepo.Update(user)
	return err
}

func (uu *UserUsecase) SetRating(rating *models.Rating) error {
	lastRating, err := uu.userRatingRepository.SelectRating(rating.UserFrom, rating.UserTo)
	var count int
	var rate int

	// если рейтинга нет и намерение удалить рейтинг
	if err == internalError.EmptyQuery && rating.Rating == 0 {
		// отвечаем, что все ок, удалено
		return nil

		// если рейтинга нет и намерение поставить оценку
	} else if err == internalError.EmptyQuery && rating.Rating > 0 {
		err = uu.userRatingRepository.InsertRating(rating)
		rate = rating.Rating
		count = 1

		// если рейтинг уже есть и намерение удалить его
	} else if err == nil && rating.Rating == 0 {
		err = uu.userRatingRepository.DeleteRating(rating)
		// поставим отрицательную оценку, чтобы бд корректно удалила
		rate = -lastRating.Rating
		count = -1

		// если рейтинг уже есть и намерение обновить его
	} else if err == nil && rating.Rating > 0 {
		err = uu.userRatingRepository.UpdateRating(rating)
		rate = rating.Rating - lastRating.Rating
		count = 0
	} else {
		return err
	}

	if err != nil {
		return err
	}

	// после обновления rating необходимо обновить статистику
	err = uu.userRatingRepository.UpdateStat(rating.UserTo, rate, count)
	return err
}

func (uu *UserUsecase) GetRating(userFrom int64, userTo int64) (*models.RatingStat, error) {
	sum, count, err := uu.userRatingRepository.SelectStat(userTo)
	if err != nil {
		return nil, err
	}

	var avg float32 = 0.0
	if count > 0 {
		avg = float32(sum) / float32(count)
	}
	ratingStat := &models.RatingStat{
		RatingSum:   sum,
		RatingCount: count,
		RatingAvg:   avg,
	}

	rating, err := uu.userRatingRepository.SelectRating(userFrom, userTo)
	if userFrom == 0 || userFrom == userTo || err != nil {
		ratingStat.IsRated = false
		ratingStat.PersonalRate = 0

		switch err {
		case internalError.EmptyQuery:
			return ratingStat, nil

		default:
			return nil, err
		}
	}

	ratingStat.IsRated = true
	ratingStat.PersonalRate = int(rating.Rating)
	return ratingStat, nil
}
