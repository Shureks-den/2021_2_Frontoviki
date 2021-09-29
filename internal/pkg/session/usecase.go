package session

import "yula/internal/models"

type SessionUsecase interface {
	Check(value string) (*models.Session, error)
	Create(userId int64) (*models.Session, error)
	Delete(value string) error
}
