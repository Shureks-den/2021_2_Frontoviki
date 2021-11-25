package repository

import (
	"database/sql/driver"
	"fmt"
	"testing"
	"time"
	"yula/internal/models"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	testime string = "2014-11-12 11:45:26.371"
	layout  string = "2006-01-02 15:04:05.000"
)

func ParseTime() time.Time {
	te, _ := time.Parse(layout, testime)
	return te
}

var testuser = &models.UserData{
	Id:        0,
	Name:      "Ваня",
	Surname:   "Иванов",
	Email:     "ivan@mail.ru",
	Password:  "1234",
	CreatedAt: ParseTime(),
	Image:     "default_image",
	Phone:     "89999999999",
}

var testrating = &models.Rating{
	UserFrom: 1,
	UserTo:   2,
	Rating:   5,
}

var teststat = &models.RatingStat{
	RatingSum:    9,
	RatingCount:  2,
	RatingAvg:    4.5,
	PersonalRate: 4,
	IsRated:      true,
}

func TimeToString(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.185743", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func TestUserInsertOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(testuser.Id)
	mock.ExpectQuery("INSERT").WithArgs(testuser.Email, testuser.Password, testuser.CreatedAt,
		testuser.Name, testuser.Surname, testuser.Image, testuser.Phone).WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO rating_statistics").WithArgs(testuser.Id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.Insert(testuser)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserInsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(testuser.Email)
	mock.ExpectQuery("INSERT").WithArgs(testuser.Email, testuser.Password, testuser.CreatedAt,
		testuser.Name, testuser.Surname, testuser.Image, testuser.Phone).WillReturnRows(rows)

	mock.ExpectRollback()

	err = repo.Insert(testuser)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserInsertError2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(testuser.Id)
	mock.ExpectQuery("INSERT").WithArgs(testuser.Email, testuser.Password, testuser.CreatedAt,
		testuser.Name, testuser.Surname, testuser.Image, testuser.Phone).WillReturnRows(rows)
	mock.ExpectExec("INSERT").WithArgs(testuser.Id)
	mock.ExpectRollback()

	err = repo.Insert(testuser)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserSelectByEmailOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "email", "phone", "password", "created_at", "name", "surname", "image"})
	rows.AddRow(testuser.Id, testuser.Email, testuser.Phone, testuser.Password, testuser.CreatedAt,
		testuser.Name, testuser.Surname, testuser.Image,
	)
	mock.ExpectQuery("SELECT").WithArgs(testuser.Email).WillReturnRows(rows)

	user, err := repo.SelectByEmail(testuser.Email)

	assert.Equal(t, testuser, user)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserSelectByEmailError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "email", "phone", "password", "created_at", "name", "surname", "image"})
	rows.AddRow(testuser.Id, testuser.Email, testuser.Phone, testuser.Password, testime,
		testuser.Name, testuser.Surname, testuser.Image,
	)
	mock.ExpectQuery("SELECT").WithArgs(testuser.Email).WillReturnRows(rows)

	user, err := repo.SelectByEmail(testuser.Email)

	assert.NotEqual(t, testuser, user)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserSelectByIdOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "email", "phone", "password", "created_at", "name", "surname", "image"})
	rows.AddRow(testuser.Id, testuser.Email, testuser.Phone, testuser.Password, testuser.CreatedAt,
		testuser.Name, testuser.Surname, testuser.Image,
	)
	mock.ExpectQuery("SELECT").WithArgs(testuser.Id).WillReturnRows(rows)

	user, err := repo.SelectById(testuser.Id)

	assert.Equal(t, testuser, user)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserSelectByIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "email", "phone", "password", "created_at", "name", "surname", "image"})
	rows.AddRow(testuser.Id, testuser.Email, testuser.Phone, testuser.Password, testime,
		testuser.Name, testuser.Surname, testuser.Image,
	)
	mock.ExpectQuery("SELECT").WithArgs(testuser.Id).WillReturnRows(rows)

	user, err := repo.SelectById(testuser.Id)

	assert.NotEqual(t, testuser, user)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserUpdateOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(testuser.Id, testuser.Email, testuser.Password,
		testuser.Name, testuser.Surname, testuser.Image, testuser.Phone).WillReturnResult(driver.RowsAffected(1))
	mock.ExpectCommit()

	err = repo.Update(testuser)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserUpdateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(testuser.Id, testuser.Email, testuser.Password,
		testuser.Name, testuser.Surname, testuser.Image, testuser.Phone).WillReturnResult(driver.RowsAffected(0))
	mock.ExpectRollback()

	err = repo.Update(testuser)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserSelectRatingOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	rows := sqlmock.NewRows([]string{"rate"})
	rows.AddRow(testrating.Rating)
	mock.ExpectQuery("SELECT").WithArgs(testrating.UserFrom, testrating.UserTo).WillReturnRows(rows)

	rate, err := repo.SelectRating(testrating.UserFrom, testrating.UserTo)

	assert.Equal(t, rate.Rating, testrating.Rating)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserSelectRatingError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	rows := sqlmock.NewRows([]string{"rate"})
	rows.AddRow(testrating.Rating)
	mock.ExpectQuery("SELECT").WithArgs(testrating.UserFrom, testrating.UserTo)

	_, err = repo.SelectRating(testrating.UserFrom, testrating.UserTo)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserInsertRatingOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(testrating.UserFrom, testrating.UserTo, testrating.Rating).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.InsertRating(testrating)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserInsertRatingError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(testrating.UserFrom, testrating.UserTo, testrating.Rating)
	mock.ExpectRollback()

	err = repo.InsertRating(testrating)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserDeleteRatingOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testrating.UserFrom, testrating.UserTo).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.DeleteRating(testrating)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserDeleteRatingError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testrating.UserFrom, testrating.UserTo)
	mock.ExpectRollback()

	err = repo.DeleteRating(testrating)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserUpdateRatingOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(testrating.UserFrom, testrating.UserTo, testrating.Rating).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.UpdateRating(testrating)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUserUpdateRatingError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(testrating.UserFrom, testrating.UserTo, testrating.Rating)
	mock.ExpectRollback()

	err = repo.UpdateRating(testrating)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectStatOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	rows := sqlmock.NewRows([]string{"sum", "count"})
	rows.AddRow(teststat.RatingSum, teststat.RatingCount)
	mock.ExpectQuery("SELECT").WithArgs(testuser.Id).WillReturnRows(rows)

	_, _, err = repo.SelectStat(testuser.Id)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectStatError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	rows := sqlmock.NewRows([]string{"sum", "count"})
	rows.AddRow(teststat.RatingSum, teststat.RatingCount)
	mock.ExpectQuery("SELECT").WithArgs(testuser.Id)

	_, _, err = repo.SelectStat(testuser.Id)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertStatOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(testuser.Id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.InsertStat(testuser.Id)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertStatError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(testuser.Id)
	mock.ExpectRollback()

	err = repo.InsertStat(testuser.Id)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateStatOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(testuser.Id, teststat.PersonalRate, 1).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.UpdateStat(testuser.Id, teststat.PersonalRate, 1)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateStatError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewRatingRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(testuser.Id, teststat.PersonalRate, 1)
	mock.ExpectRollback()

	err = repo.UpdateStat(testuser.Id, teststat.PersonalRate, 1)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
