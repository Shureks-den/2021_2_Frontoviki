package repository

import (
	"fmt"
	"strings"
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

var testimages = fmt.Sprintf("{%s}", strings.Join(testadvert.Images, ", "))

var testpage = &models.Page{
	PageNum: 1,
	Count:   50,
}

func TestSelectWithFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sf := &models.SearchFilter{
		Query: "aboba", Category: "boba", Date: ParseTime(), TimeDuration: 30, Latitude: 30.0, Longitude: 30.0,
		Radius: 500, SortingDate: true, SortingName: false,
	}

	repo := NewSearchRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew,
	)
	mock.ExpectQuery("SELECT").WithArgs(sf.Query, sf.Category, sf.Date, sf.TimeDuration, sf.Longitude, sf.Latitude, sf.Radius).
		WillReturnRows(rows)

	ads, err := repo.SelectWithFilter(sf, testpage.PageNum, testpage.Count)

	assert.NoError(t, err)
	assert.NotNil(t, ads)
}

func TestSelectWithFilterError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	sf := &models.SearchFilter{
		Query: "aboba", Category: "boba", Date: ParseTime(), TimeDuration: 30, Latitude: 30.0, Longitude: 30.0,
		Radius: 500, SortingDate: true, SortingName: true,
	}

	repo := NewSearchRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew,
	)
	mock.ExpectQuery("SELECT").WithArgs(sf.Query, sf.Category, sf.Date, sf.TimeDuration, sf.Longitude, sf.Latitude, sf.Radius)

	_, err = repo.SelectWithFilter(sf, testpage.PageNum, testpage.Count)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
