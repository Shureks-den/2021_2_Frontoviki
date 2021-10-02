package user

import (
	"yula/internal/models"
)

// определяем интерфейс для взаимодействия с бд
type UserRepository interface {
	Insert(user *models.UserData) error
	SelectByEmail(email string) (*models.UserData, error)
	SelectById(userId int64) (*models.UserData, error)
	Update(user *models.UserData) error
}
