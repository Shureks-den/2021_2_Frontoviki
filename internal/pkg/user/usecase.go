package user

import (
	"yula/internal/models"
)

// определяем интерфейс связи между deliver и repository
type UserUsecase interface {
	Create(user *models.UserSignUp) (*models.UserData, *models.Status)
	GetByEmail(email string) (*models.UserData, *models.Status)
}
