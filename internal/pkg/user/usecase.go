package user

import (
	"yula/internal/codes"
	"yula/internal/models"
)

// определяем интерфейс связи между deliver и repository
type UserUsecase interface {
	Create(user *models.UserSignUp) (*models.UserData, *codes.ServerError)
	GetByEmail(email string) (*models.UserData, *codes.ServerError)
	CheckPassword(user *models.UserData, gettedPassword string) *codes.ServerError
}
