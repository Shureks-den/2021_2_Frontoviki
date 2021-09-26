package usecase

import (
	"time"
	"yula/internal/models"
	"yula/internal/pkg/session"

	"github.com/google/uuid"
)

type SessionUsecase struct {
	sessionRepo session.SessionRepository
}

func NewSessionUsecase(repo session.SessionRepository) session.SessionUsecase {
	return &SessionUsecase{
		sessionRepo: repo,
	}
}

func (su *SessionUsecase) Create(userId int64) (*models.Session, error) {
	sess := models.Session{Value: uuid.NewString(), UserId: userId, ExpiresAt: time.Now().Add(time.Minute)}

	if err := su.sessionRepo.Set(&sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

func (su *SessionUsecase) Delete(value string) error {
	sess, err := su.sessionRepo.GetByValue(value)
	if err != nil {
		return err
	}

	err = su.sessionRepo.Delete(sess)
	return err
}

func (su *SessionUsecase) Check(value string) (*models.Session, error) {
	sess, err := su.sessionRepo.GetByValue(value)
	if err != nil {
		return nil, err
	}
	return sess, nil
}
