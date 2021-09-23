package user

import (
	"yula/internal/models"

	"github.com/google/uuid"
)

// определяем интерфейс для взаимодействия с бд
type UserRepository interface {
	Insert(user *models.UserData) error
	SelectByEmail(email string) (*models.UserData, error)
	SelectById(userId uuid.UUID) (*models.UserData, error)
	Update(user *models.UserData) error
}
