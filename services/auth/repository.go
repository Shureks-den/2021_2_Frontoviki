package session

import "yula/internal/models"

//go:generate mockery -name=SessionRepository

type SessionRepository interface {
	Set(sess *models.Session) error
	Delete(sess *models.Session) error
	GetByValue(value string) (*models.Session, error)
}
