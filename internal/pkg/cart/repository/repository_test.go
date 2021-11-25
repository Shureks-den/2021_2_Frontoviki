package repository

import (
	"database/sql/driver"
	"testing"
	"time"
	"yula/internal/models"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func ParseTime() time.Time {
	testime := "2014-11-12 11:45:26.371"
	layout := "2006-01-02 15:04:05.000"
	te, _ := time.Parse(layout, testime)
	return te
}

var testadvert = &models.Advert{
	Id:          0,
	Name:        "объявление",
	Description: "описание",
	Price:       100,
	Location:    "москва",
	Latitude:    55.75,
	Longitude:   37.62,
	PublishedAt: ParseTime(),
	DateClose:   ParseTime(),
	IsActive:    true,
	PublisherId: 0,
	Category:    "одежда",
	Images:      []string{"default"},
	Views:       1,
	Amount:      1,
	IsNew:       true,
}

var testuserid = int64(1)

func TestSelectOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	rows := sqlmock.NewRows([]string{"user_id", "advert_id", "amount"}).AddRow(testuserid, testadvert.Id, testadvert.Amount)
	mock.ExpectQuery("SELECT").WithArgs(testuserid, testadvert.Id).WillReturnRows(rows)

	_, err = repo.Select(testuserid, testadvert.Id)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	mock.ExpectQuery("SELECT").WithArgs(testuserid, testadvert.Id)

	_, err = repo.Select(testuserid, testadvert.Id)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectAllOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	rows := sqlmock.NewRows([]string{"user_id", "advert_id", "amount"}).AddRow(testuserid, testadvert.Id, testadvert.Amount)
	mock.ExpectQuery("SELECT").WithArgs(testuserid).WillReturnRows(rows)

	_, err = repo.SelectAll(testuserid)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectAllError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	mock.ExpectQuery("SELECT").WithArgs(testuserid)

	_, err = repo.SelectAll(testuserid)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

var testcart = &models.Cart{
	UserId:   testuserid,
	AdvertId: testadvert.Id,
	Amount:   testadvert.Amount,
}

func TestUpdateOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(testcart.UserId, testcart.AdvertId, testcart.Amount).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.Update(testcart)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(testcart.UserId, testcart.AdvertId, testcart.Amount)
	mock.ExpectRollback()

	err = repo.Update(testcart)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(testcart.UserId, testcart.AdvertId, testcart.Amount).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.Insert(testcart)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(testcart.UserId, testcart.AdvertId, testcart.Amount)
	mock.ExpectRollback()

	err = repo.Insert(testcart)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testcart.UserId, testcart.AdvertId).WillReturnResult(driver.RowsAffected(1))
	mock.ExpectCommit()

	err = repo.Delete(testcart)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testcart.UserId, testcart.AdvertId)
	mock.ExpectRollback()

	err = repo.Delete(testcart)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteAllOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testcart.UserId).WillReturnResult(driver.RowsAffected(1))
	mock.ExpectCommit()

	err = repo.DeleteAll(testcart.UserId)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteAllError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testcart.UserId)
	mock.ExpectRollback()

	err = repo.DeleteAll(testcart.UserId)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteAllError2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testcart.UserId).WillReturnResult(driver.RowsAffected(0))
	mock.ExpectCommit()

	err = repo.DeleteAll(testcart.UserId)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
