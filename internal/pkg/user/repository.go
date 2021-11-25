package user

import (
	"yula/internal/models"
)

//go:generate mockery -name=UserRepository

// определяем интерфейс для взаимодействия с бд
type UserRepository interface {
	Insert(user *models.UserData) error
	SelectByEmail(email string) (*models.UserData, error)
	SelectById(userId int64) (*models.UserData, error)
	Update(user *models.UserData) error
}

type RatingRepository interface {
	SelectRating(userFrom int64, userTo int64) (*models.Rating, error)
	InsertRating(rating *models.Rating) error
	UpdateRating(rating *models.Rating) error
	DeleteRating(rating *models.Rating) error

	SelectStat(userId int64) (int64, int64, error)
	InsertStat(userId int64) error
	UpdateStat(userId int64, rate int, count int) error
}
