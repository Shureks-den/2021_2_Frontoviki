package usecase

import (
	"time"
	"yula/internal/models"
	"yula/internal/pkg/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(repo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRepo: repo,
	}
}

func (uu *UserUsecase) Create(userSU *models.UserSignUp) (*models.UserData, *models.Status) {
	if _, err := uu.GetByEmail(userSU.Email); err != models.StatusByCode(models.UserNotExist) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userSU.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, models.StatusByCode(models.InternalError)
	}

	userSU.Password = string(passwordHash)
	user := uu.UserSignUpToUserData(userSU)
	err = uu.userRepo.Insert(user)

	if err != nil {
		return nil, models.StatusByCode(models.InternalError)
	}

	return user, models.StatusByCode(models.Created)
}

func (uu *UserUsecase) UserSignUpToUserData(userSU *models.UserSignUp) *models.UserData {
	var user models.UserData
	user.Id = uuid.New()
	user.Username = userSU.Username
	user.Email = userSU.Email
	user.Password = userSU.Password
	user.CreatedAt = time.Now()
	return &user
}

func (uu *UserUsecase) GetByEmail(email string) (*models.UserData, *models.Status) {
	user, err := uu.userRepo.SelectByEmail(email)

	switch {
	// ^ other cases
	case err != nil:
		return nil, models.StatusByCode(models.UserNotExist)
	}

	return user, models.StatusByCode(models.OK)
}
