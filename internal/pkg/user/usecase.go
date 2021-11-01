package user

import (
	"mime/multipart"
	"yula/internal/models"
)

//go:generate mockery -name=UserUsecase

// определяем интерфейс связи между deliver и repository
type UserUsecase interface {
	Create(user *models.UserSignUp) (*models.UserData, error)
	GetByEmail(email string) (*models.UserData, error)

	CheckPassword(user *models.UserData, gettedPassword string) error
	UpdatePassword(userId int64, changePassword *models.ChangePassword) error

	GetById(id int64) (*models.Profile, error)
	UpdateProfile(userId int64, userNew *models.UserData) (*models.Profile, error)
	UploadAvatar(file *multipart.FileHeader, userId int64) (*models.UserData, error)
}
