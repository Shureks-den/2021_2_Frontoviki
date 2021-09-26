package usecase

import (
	"time"
	"yula/internal/codes"
	"yula/internal/models"
	"yula/internal/pkg/user"

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

func (uu *UserUsecase) Create(userSU *models.UserSignUp) (*models.UserData, *codes.ServerError) {
	if _, err := uu.GetByEmail(userSU.Email); err != codes.NewServerError(codes.UserNotExist) {
		switch err {
		case nil:
			return nil, codes.NewServerError(codes.UserAlreadyExist)

		default:
			return nil, codes.NewServerError(codes.InternalError)
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userSU.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, codes.NewServerError(codes.InternalError)
	}

	userSU.Password = string(passwordHash)
	user := models.UserData{}
	user.Username = userSU.Username
	user.Email = userSU.Email
	user.Password = userSU.Password
	user.CreatedAt = time.Now()

	dbErr := uu.userRepo.Insert(&user)

	if dbErr != nil {
		return nil, codes.NewServerError(codes.InternalError)
	}

	return &user, nil
}

func (uu *UserUsecase) GetByEmail(email string) (*models.UserData, *codes.ServerError) {
	user, err := uu.userRepo.SelectByEmail(email)

	if err == nil {
		return user, nil
	}

	switch err.Error {
	case codes.EmptyRow:
		return nil, codes.NewServerError(codes.UserNotExist)
	default:
		return nil, codes.NewServerError(codes.UnexpectedError)
	}
}

func (uu *UserUsecase) CheckPassword(user *models.UserData, gettedPassword string) *codes.ServerError {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(gettedPassword))
	if err != nil {
		return codes.NewServerError(codes.Unauthorized)
	}
	return nil
}

func (uu *UserUsecase) GetById(user_id int64) (*models.Profile, *codes.ServerError) {
	user, err := uu.userRepo.SelectById(user_id)

	if err == nil {
		return user.ToProfile(), nil
	}

	switch err.Error {
	case codes.EmptyRow:
		return nil, codes.NewServerError(codes.UserNotExist)
	default:
		return nil, codes.NewServerError(codes.UnexpectedError)
	}
}
