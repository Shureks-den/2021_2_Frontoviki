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

func TimeToString(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.185743", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func TestUserInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(testuser.Id)
	mock.ExpectQuery(`INSERT INTO users`).WithArgs(testuser.Email, testuser.Password, testime,
		testuser.Name, testuser.Surname, testuser.Image, testuser.Phone).WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO rating_statistics").WithArgs(testuser.Id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.Insert(testuser)
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func Test(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").
		WithArgs("john", "aboba").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = db.Exec("INSERT INTO users(name, created_at) VALUES (?, ?)", "john", "aboba")
	if err != nil {
		t.Errorf("error '%s' was not expected, while inserting a row", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
