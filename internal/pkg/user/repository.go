package user

import (
	"yula/internal/codes"
	"yula/internal/models"
)

// определяем интерфейс для взаимодействия с бд
type UserRepository interface {
	Insert(user *models.UserData) *codes.DatabaseError
	SelectByEmail(email string) (*models.UserData, *codes.DatabaseError)
	SelectById(userId int64) (*models.UserData, *codes.DatabaseError)
	Update(user *models.UserData) *codes.DatabaseError
}
