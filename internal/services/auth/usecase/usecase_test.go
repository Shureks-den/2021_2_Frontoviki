package usecase

import (
	"testing"
	"time"
	"yula/internal/models"

	myerr "yula/internal/error"
	"yula/internal/services/auth/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSession_SignInHandler_CreateSuccess(t *testing.T) {
	sr := mocks.SessionRepository{}
	su := NewSessionUsecase(&sr)

	sess := models.Session{Value: uuid.NewString(), UserId: 10529, ExpiresAt: time.Now().Add(time.Minute)}
	sr.On("Set", mock.MatchedBy(func(session *models.Session) bool {
		return session.UserId == sess.UserId
	})).Return(nil)

	session, err := su.Create(sess.UserId)
	assert.Nil(t, err)
	assert.Equal(t, session.UserId, int64(10529))
}

func TestSession_SignInHandler_CreateNotSuccess(t *testing.T) {
	sr := mocks.SessionRepository{}
	su := NewSessionUsecase(&sr)

	sess := models.Session{Value: uuid.NewString(), UserId: -1, ExpiresAt: time.Now().Add(time.Minute)}
	sr.On("Set", mock.MatchedBy(func(session *models.Session) bool {
		return session.UserId == sess.UserId
	})).Return(myerr.DatabaseError)

	session, err := su.Create(sess.UserId)
	assert.Equal(t, err, myerr.DatabaseError)
	assert.Nil(t, session)
}

func TestSession_SignInHandler_DeleteSuccess(t *testing.T) {
	sr := mocks.SessionRepository{}
	su := NewSessionUsecase(&sr)

	sess := models.Session{Value: uuid.NewString(), UserId: 1, ExpiresAt: time.Now().Add(time.Minute)}
	sr.On("GetByValue", sess.Value).Return(&sess, nil)
	sr.On("Delete", &sess).Return(nil)

	err := su.Delete(sess.Value)
	assert.Equal(t, err, nil)
}

func TestSession_SignInHandler_DeleteError(t *testing.T) {
	sr := mocks.SessionRepository{}
	su := NewSessionUsecase(&sr)

	sr.On("GetByValue", "empty").Return(nil, myerr.DatabaseError)

	err := su.Delete("empty")
	assert.Equal(t, err, myerr.DatabaseError)
}

func TestSession_SignInHandler_CheckSuccess(t *testing.T) {
	sr := mocks.SessionRepository{}
	su := NewSessionUsecase(&sr)

	sess := models.Session{Value: uuid.NewString(), UserId: 1, ExpiresAt: time.Now().Add(time.Minute)}
	sr.On("GetByValue", sess.Value).Return(&sess, nil)

	session, err := su.Check(sess.Value)
	assert.Equal(t, *session, sess)
	assert.Nil(t, err)
}

func TestSession_SignInHandler_CheckError(t *testing.T) {
	sr := mocks.SessionRepository{}
	su := NewSessionUsecase(&sr)

	sr.On("GetByValue", "").Return(nil, myerr.DatabaseError)

	session, err := su.Check("")
	assert.Equal(t, err, myerr.DatabaseError)
	assert.Nil(t, session)
}
