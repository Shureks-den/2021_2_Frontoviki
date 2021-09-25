package repository

import (
	"errors"
	"time"
	"yula/internal/models"
	"yula/internal/pkg/session"

	"github.com/google/uuid"
)

type SessionRepository struct {
	pool []models.Session
}

func NewSessionRepository() session.SessionRepository {
	var pool []models.Session
	pool = append(pool, models.Session{Value: uuid.NewString(), UserId: 1, ExpiresAt: time.Now().Add(time.Hour)})

	return &SessionRepository{
		pool: pool,
	}
}

func (sr *SessionRepository) Set(sess *models.Session) error {
	sr.pool = append(sr.pool, *sess)
	return nil
}

func (sr *SessionRepository) Delete(sess *models.Session) error {
	for ind, val := range sr.pool {
		if val.Value == sess.Value && val.UserId == sess.UserId && val.ExpiresAt == sess.ExpiresAt {
			sr.pool = append(sr.pool[:ind], sr.pool[ind+1:]...)
			return nil
		}
	}
	return errors.New("error")
}

func (sr *SessionRepository) GetByValue(value string) (*models.Session, error) {
	for _, val := range sr.pool {
		if val.Value == value {
			return &val, nil
		}
	}
	return nil, errors.New("no session")
}
