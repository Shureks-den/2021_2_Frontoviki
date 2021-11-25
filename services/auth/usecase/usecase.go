package usecase

import (
	"time"
	internalError "yula/internal/error"
	"yula/internal/models"
	session "yula/services/auth"

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
	sess := models.Session{Value: uuid.NewString(), UserId: userId, ExpiresAt: time.Now().Add(time.Hour)}

	if err := su.sessionRepo.Set(&sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

func (su *SessionUsecase) Delete(value string) error {
	sess, err := su.sessionRepo.GetByValue(value)
	if err != nil {
		switch err {
		case internalError.EmptyQuery:
			return internalError.NotExist
		default:
			return err
		}
	}

	err = su.sessionRepo.Delete(sess)
	return err
}

func (su *SessionUsecase) Check(value string) (*models.Session, error) {
	sess, err := su.sessionRepo.GetByValue(value)
	if err != nil {
		switch err {
		case internalError.EmptyQuery:
			return nil, internalError.NotExist
		default:
			return nil, err
		}
	}
	return sess, nil
}
