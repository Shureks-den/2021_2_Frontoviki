package repository

import (
	"testing"
	"time"
	"yula/internal/models"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func ParseTime() time.Time {
	testime := "2014-11-12 11:45:26.371"
	layout := "2006-01-02 15:04:05.000"
	te, _ := time.Parse(layout, testime)
	return te
}

func TestSelectMessagesOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	message := models.Message{models.IMessage{0, 1, 1}, "qwerty", ParseTime()}
	repo := NewChatRepository(db)
	rows := sqlmock.NewRows([]string{"user_from", "user_to", "adv_id", "msg", "created_at"})
	rows.AddRow(message.MI.IdFrom, message.MI.IdTo, message.MI.IdAdv, message.Msg, message.CreatedAt)
	mock.ExpectQuery("SELECT").WithArgs(message.MI.IdFrom, message.MI.IdTo, message.MI.IdAdv, int64(0), int64(10)).WillReturnRows(rows)

	_, err = repo.SelectMessages(&message.MI, int64(0), int64(10))

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
