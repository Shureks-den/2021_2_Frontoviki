package user

import (
	"mime/multipart"
	"yula/internal/codes"
	"yula/internal/models"
)

// определяем интерфейс связи между deliver и repository
type UserUsecase interface {
	Create(user *models.UserSignUp) (*models.UserData, *codes.ServerError)
	GetByEmail(email string) (*models.UserData, *codes.ServerError)
	CheckPassword(user *models.UserData, gettedPassword string) *codes.ServerError

	GetById(id int64) (*models.Profile, *codes.ServerError)
	UpdateProfile(userId int64, userNew *models.UserData) (*models.Profile, *codes.ServerError)
	CheckEmail(email string) *codes.ServerError
	UploadAvatar(file *multipart.FileHeader, userId int64) *codes.ServerError // пока не работает
}
