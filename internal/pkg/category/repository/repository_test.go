package repository

import (
	"testing"
	"yula/internal/models"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestSelectCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &CategoryRepository{
		DB: db,
	}
	rows := sqlmock.
		NewRows([]string{"name"})
	expect := []*models.Category{
		{"aboba"},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.Name)
	}

	mock.
		ExpectQuery("SELECT").
		WillReturnRows(rows)

	_, err = repo.SelectCategories()
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)

}
