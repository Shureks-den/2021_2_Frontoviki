package repository

import (
	"database/sql/driver"
	"testing"
	"time"
	myerr "yula/internal/error"
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

	message := models.Message{MI: models.IMessage{0, 1, 1}, Msg: "qwerty", CreatedAt: ParseTime()}
	repo := NewChatRepository(db)
	rows := sqlmock.NewRows([]string{"user_from", "user_to", "adv_id", "msg", "created_at"})
	rows.AddRow(message.MI.IdFrom, message.MI.IdTo, message.MI.IdAdv, message.Msg, message.CreatedAt)
	mock.ExpectQuery("SELECT").WithArgs(message.MI.IdFrom, message.MI.IdTo, message.MI.IdAdv, int64(0), int64(10)).WillReturnRows(rows)

	_, err = repo.SelectMessages(&message.MI, int64(0), int64(10))

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertMessageOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	message := models.Message{MI: models.IMessage{0, 1, 1}, Msg: "qwerty", CreatedAt: ParseTime()}
	repo := NewChatRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(message.MI.IdFrom, message.MI.IdTo, message.MI.IdAdv, message.Msg).
		WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.InsertMessage(&message)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertMessageError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	message := models.Message{MI: models.IMessage{0, 1, 1}, Msg: "qwerty", CreatedAt: ParseTime()}
	repo := NewChatRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(message.MI.IdFrom, message.MI.IdTo, message.MI.IdAdv, message.Msg)
	mock.ExpectRollback()

	err = repo.InsertMessage(&message)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteMessagesOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	message := models.Message{MI: models.IMessage{0, 1, 1}, Msg: "qwerty", CreatedAt: ParseTime()}
	repo := NewChatRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(message.MI.IdFrom, message.MI.IdTo, message.MI.IdAdv).
		WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.DeleteMessages(&message.MI)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteMessagesError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	message := models.Message{MI: models.IMessage{0, 1, 1}, Msg: "qwerty", CreatedAt: ParseTime()}
	repo := NewChatRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(message.MI.IdFrom, message.MI.IdTo, message.MI.IdAdv)
	mock.ExpectRollback()

	err = repo.DeleteMessages(&message.MI)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectDialogOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	dialog := models.Dialog{DI: models.IDialog{Id1: 0, Id2: 1, IdAdv: 1}, CreatedAt: ParseTime()}
	repo := NewChatRepository(db)
	rows := sqlmock.NewRows([]string{"user1", "user2", "adv_id", "created_at"})
	rows.AddRow(dialog.DI.Id1, dialog.DI.Id2, dialog.DI.IdAdv, dialog.CreatedAt)
	mock.ExpectQuery("SELECT").WithArgs(dialog.DI.Id1, dialog.DI.Id2, dialog.DI.IdAdv).WillReturnRows(rows)

	_, err = repo.SelectDialog(&dialog.DI)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectDialogError1(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	dialog := models.Dialog{DI: models.IDialog{Id1: 0, Id2: 1, IdAdv: 1}, CreatedAt: ParseTime()}
	repo := NewChatRepository(db)
	rows := sqlmock.NewRows([]string{"user1", "user2", "adv_id", "created_at"})
	rows.AddRow(dialog.DI.Id1, dialog.DI.Id2, dialog.DI.IdAdv, dialog.CreatedAt)
	mock.ExpectQuery("SELECT").WithArgs(dialog.DI.Id1, dialog.DI.Id2, dialog.DI.IdAdv).WillReturnError(myerr.EmptyQuery)

	_, err = repo.SelectDialog(&dialog.DI)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertDialogOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	dialog := models.Dialog{DI: models.IDialog{Id1: 0, Id2: 1, IdAdv: 1}, CreatedAt: ParseTime()}
	repo := NewChatRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(dialog.DI.Id1, dialog.DI.Id2, dialog.DI.IdAdv).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.InsertDialog(&dialog)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertDialogError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	dialog := models.Dialog{DI: models.IDialog{Id1: 0, Id2: 1, IdAdv: 1}, CreatedAt: ParseTime()}
	repo := NewChatRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(dialog.DI.Id1, dialog.DI.Id2, dialog.DI.IdAdv)
	mock.ExpectRollback()

	err = repo.InsertDialog(&dialog)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteDialogOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	dialog := models.Dialog{DI: models.IDialog{Id1: 0, Id2: 1, IdAdv: 1}, CreatedAt: ParseTime()}
	repo := NewChatRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(dialog.DI.Id1, dialog.DI.Id2, dialog.DI.IdAdv).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.DeleteDialog(&dialog.DI)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteDialogError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	dialog := models.Dialog{DI: models.IDialog{Id1: 0, Id2: 1, IdAdv: 1}, CreatedAt: ParseTime()}
	repo := NewChatRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(dialog.DI.Id1, dialog.DI.Id2, dialog.DI.IdAdv)
	mock.ExpectRollback()

	err = repo.DeleteDialog(&dialog.DI)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectAllDialogsOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	dialog := models.Dialog{DI: models.IDialog{Id1: 0, Id2: 1, IdAdv: 1}, CreatedAt: ParseTime()}
	repo := NewChatRepository(db)
	rows := sqlmock.NewRows([]string{"user1", "user2", "adv_id", "created_at"})
	rows.AddRow(dialog.DI.Id1, dialog.DI.Id2, dialog.DI.IdAdv, dialog.CreatedAt)
	mock.ExpectQuery("SELECT").WithArgs(dialog.DI.Id1).WillReturnRows(rows)

	_, err = repo.SelectAllDialogs(dialog.DI.Id1)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectAllDialogsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("canot create mock: %s", err)
	}
	defer db.Close()

	dialog := models.Dialog{DI: models.IDialog{Id1: 0, Id2: 1, IdAdv: 1}, CreatedAt: ParseTime()}
	repo := NewChatRepository(db)
	rows := sqlmock.NewRows([]string{"user1", "user2", "adv_id", "created_at"})
	rows.AddRow(dialog.DI.Id1, dialog.DI.Id2, dialog.DI.IdAdv, dialog.CreatedAt)
	mock.ExpectQuery("SELECT").WithArgs(dialog.DI.Id1).WillReturnError(myerr.InternalError)

	_, err = repo.SelectAllDialogs(dialog.DI.Id1)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
