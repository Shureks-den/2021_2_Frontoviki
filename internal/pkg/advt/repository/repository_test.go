package repository

import (
	"database/sql/driver"
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
	PromoLevel:  0,
}

var testimages = fmt.Sprintf("{%s}", strings.Join(testadvert.Images, ", "))

var testpage = &models.Page{
	PageNum: 1,
	Count:   50,
}

func TestSelectListAdvtOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new", "p.promo_level"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew, testadvert.PromoLevel,
	)
	mock.ExpectQuery("SELECT").WithArgs(testpage.Count, testpage.PageNum*testpage.Count).WillReturnRows(rows)

	_, err = repo.SelectListAdvt(true, testpage.PageNum, testpage.Count)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectListAdvtError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew,
	)
	mock.ExpectQuery("SELECT").WithArgs(testpage.Count, testpage.PageNum*testpage.Count)

	_, err = repo.SelectListAdvt(true, testpage.PageNum, testpage.Count)

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

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(testadvert.Id)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WithArgs(testadvert.Name, testadvert.Description, testadvert.Category, testadvert.PublisherId,
		testadvert.Latitude, testadvert.Longitude, testadvert.Location, testadvert.Price, testadvert.Amount, testadvert.IsNew).
		WillReturnRows(rows)
	mock.ExpectExec("INSERT").WithArgs(testadvert.Id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("INSERT").WithArgs(testadvert.Id, testadvert.Price).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("INSERT").WithArgs(testadvert.Id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.Insert(testadvert)

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

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(testadvert.Id)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WithArgs(testadvert.Name, testadvert.Description, testadvert.Category, testadvert.PublisherId,
		testadvert.Latitude, testadvert.Longitude, testadvert.Location, testadvert.Price, testadvert.Amount, testadvert.IsNew).
		WillReturnRows(rows)
	mock.ExpectExec("INSERT").WithArgs(testadvert.Id)
	mock.ExpectRollback()

	err = repo.Insert(testadvert)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertError2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	//rows := sqlmock.NewRows([]string{"id"}).AddRow(testadvert.Id)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WithArgs(testadvert.Name, testadvert.Description, testadvert.Category, testadvert.PublisherId,
		testadvert.Latitude, testadvert.Longitude, testadvert.Location, testadvert.Price, testadvert.Amount, testadvert.IsNew)
	mock.ExpectRollback()

	err = repo.Insert(testadvert)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertError3(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(testadvert.Id)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WithArgs(testadvert.Name, testadvert.Description, testadvert.Category, testadvert.PublisherId,
		testadvert.Latitude, testadvert.Longitude, testadvert.Location, testadvert.Price, testadvert.Amount, testadvert.IsNew).
		WillReturnRows(rows)
	mock.ExpectExec("INSERT").WithArgs(testadvert.Id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("INSERT").WithArgs(testadvert.Id, testadvert.Price)
	mock.ExpectRollback()

	err = repo.Insert(testadvert)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertError4(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(testadvert.Id)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WithArgs(testadvert.Name, testadvert.Description, testadvert.Category, testadvert.PublisherId,
		testadvert.Latitude, testadvert.Longitude, testadvert.Location, testadvert.Price, testadvert.Amount, testadvert.IsNew).
		WillReturnRows(rows)
	mock.ExpectExec("INSERT").WithArgs(testadvert.Id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("INSERT").WithArgs(testadvert.Id, testadvert.Price).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("INSERT").WithArgs(testadvert.Id)
	mock.ExpectRollback()

	err = repo.Insert(testadvert)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectByIdOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new", "p.promo_level"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew, testadvert.PromoLevel,
	)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Id).WillReturnRows(rows)

	_, err = repo.SelectById(testadvert.Id)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectByIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew,
	)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Id)

	_, err = repo.SelectById(testadvert.Id)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(testadvert.Id)
	mock.ExpectBegin()
	mock.ExpectQuery("UPDATE").WithArgs(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Category, testadvert.Location,
		testadvert.Latitude, testadvert.Longitude, testadvert.Price, testadvert.IsActive, testadvert.DateClose,
		testadvert.Amount, testadvert.IsNew).WillReturnRows(rows)
	mock.ExpectCommit()

	err = repo.Update(testadvert)

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

	repo := NewAdvtRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery("UPDATE").WithArgs(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Category, testadvert.Location,
		testadvert.Latitude, testadvert.Longitude, testadvert.Price, testadvert.IsActive, testadvert.DateClose,
		testadvert.Amount, testadvert.IsNew)
	mock.ExpectRollback()

	err = repo.Update(testadvert)

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

	repo := NewAdvtRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testadvert.Id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.Delete(testadvert.Id)

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

	repo := NewAdvtRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testadvert.Id)
	mock.ExpectRollback()

	err = repo.Delete(testadvert.Id)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteImagesOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	img_paths := []string{"image_path"}
	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testadvert.Id, img_paths[0]).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.DeleteImages(img_paths, testadvert.Id)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteImageError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	img_paths := []string{"image_path"}
	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testadvert.Id, img_paths[0])
	mock.ExpectRollback()

	err = repo.DeleteImages(img_paths, testadvert.Id)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertImagesOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	img_paths := []string{"image_path"}
	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(testadvert.Id, img_paths[0]).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.InsertImages(testadvert.Id, img_paths)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertImagesError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	img_paths := []string{"image_path"}
	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(testadvert.Id, img_paths[0])
	mock.ExpectRollback()

	err = repo.InsertImages(testadvert.Id, img_paths)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelecByPublisherIdOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new", "p.promo_level"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew, testadvert.PromoLevel,
	)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.PublisherId, testpage.Count, testpage.PageNum*testpage.Count).WillReturnRows(rows)

	_, err = repo.SelectAdvertsByPublisherId(testadvert.PublisherId, true, testpage.PageNum, testpage.Count)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelecByPublisherIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new", "p.promo_level"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew, testadvert.PromoLevel,
	)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.PublisherId, testpage.Count, testpage.PageNum*testpage.Count)

	_, err = repo.SelectAdvertsByPublisherId(testadvert.PublisherId, true, testpage.PageNum, testpage.Count)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectAdvertsByCategoryOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new", "p.promo_level"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew, testadvert.PromoLevel,
	)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Category, testpage.Count, testpage.PageNum*testpage.Count).WillReturnRows(rows)

	_, err = repo.SelectAdvertsByCategory(testadvert.Category, testpage.PageNum, testpage.Count)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectAdvertsByCategoryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew,
	)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Category, testpage.Count, testpage.PageNum*testpage.Count)

	_, err = repo.SelectAdvertsByCategory(testadvert.Category, testpage.PageNum, testpage.Count)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectFavoriteCountOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)
	count := int64(1)

	rows := sqlmock.NewRows([]string{"cnt"})
	rows.AddRow(count)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Id).WillReturnRows(rows)

	_, err = repo.SelectFavoriteCount(testadvert.Id)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectFavoriteCountError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)
	count := int64(1)

	rows := sqlmock.NewRows([]string{"cnt"})
	rows.AddRow(count)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Id)

	_, err = repo.SelectFavoriteCount(testadvert.Id)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectFavoriteAdvertsOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new", "p.promo_level"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew, testadvert.PromoLevel,
	)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.PublisherId, testpage.Count, testpage.PageNum*testpage.Count).WillReturnRows(rows)

	_, err = repo.SelectFavoriteAdverts(testadvert.PublisherId, testpage.PageNum, testpage.Count)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectFavoriteAdvertsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew,
	)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.PublisherId, testpage.Count, testpage.PageNum*testpage.Count)

	_, err = repo.SelectFavoriteAdverts(testadvert.PublisherId, testpage.PageNum, testpage.Count)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectFavoriteOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new", "p.promo_level"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew, testadvert.PromoLevel,
	)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Id, testadvert.PublisherId).WillReturnRows(rows)

	_, err = repo.SelectFavorite(testadvert.PublisherId, testadvert.Id)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectFavoriteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "location", "latitude", "longitude", "published_at",
		"date_close", "is_active", "views", "publisher_id", "c.name", "array_agg(ai.img_path)", "amount", "a.is_new"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew,
	)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Id, testadvert.PublisherId)

	_, err = repo.SelectFavorite(testadvert.PublisherId, testadvert.Id)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertFavoriteOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(testadvert.PublisherId, testadvert.Id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.InsertFavorite(testadvert.PublisherId, testadvert.Id)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertFavoriteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(testadvert.PublisherId, testadvert.Id)
	mock.ExpectRollback()

	err = repo.InsertFavorite(testadvert.PublisherId, testadvert.Id)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteFavoriteOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testadvert.PublisherId, testadvert.Id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.DeleteFavorite(testadvert.PublisherId, testadvert.Id)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteFavoriteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testadvert.PublisherId, testadvert.Id)
	mock.ExpectRollback()

	err = repo.DeleteFavorite(testadvert.PublisherId, testadvert.Id)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectViewsOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	views := int64(1)
	rows := sqlmock.NewRows([]string{"count"}).AddRow(views)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Id).WillReturnRows(rows)

	_, err = repo.SelectViews(testadvert.Id)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectViewsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	// views := int64(1)
	// rows := sqlmock.NewRows([]string{"count"}).AddRow(views)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Id)

	_, err = repo.SelectViews(testadvert.Id)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateViewsOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(testadvert.Id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.UpdateViews(testadvert.Id)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateViewsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(testadvert.Id)
	mock.ExpectRollback()

	err = repo.UpdateViews(testadvert.Id)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdatePriceOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)
	ap := &models.AdvertPrice{AdvertId: testadvert.Id, Price: int64(testadvert.Price)}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(ap.AdvertId, ap.Price).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.UpdatePrice(ap)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdatePriceError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)
	ap := &models.AdvertPrice{AdvertId: testadvert.Id, Price: int64(testadvert.Price)}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(ap.AdvertId, ap.Price)
	mock.ExpectRollback()

	err = repo.UpdatePrice(ap)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectPriceHistoryOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)
	rows := sqlmock.NewRows([]string{"id", "price", "change_date"})
	rows.AddRow(testadvert.Id, testadvert.Price, testadvert.PublishedAt)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Id).WillReturnRows(rows)

	_, err = repo.SelectPriceHistory(testadvert.Id)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectPriceHistoryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)
	rows := sqlmock.NewRows([]string{"id", "price", "change_date"})
	rows.AddRow(testadvert.Id, testadvert.Price, testadvert.PublishedAt)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Id)

	_, err = repo.SelectPriceHistory(testadvert.Id)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdatePromoOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)
	promo := &models.Promotion{AdvertId: testadvert.Id, PromoLevel: 1, UpdateTime: testadvert.PublishedAt}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(promo.AdvertId, promo.PromoLevel, promo.UpdateTime).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.UpdatePromo(promo)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdatePromoError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)
	promo := &models.Promotion{AdvertId: testadvert.Id, PromoLevel: 1, UpdateTime: testadvert.PublishedAt}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(promo.AdvertId, promo.PromoLevel, promo.UpdateTime)
	mock.ExpectRollback()

	err = repo.UpdatePromo(promo)

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectRecomendationsOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"r1.rec_id", "t1.name", "t1.description", "t1.price", "t1.location",
		"t1.latitude", "t1.longitude", "t1.published_at", "t1.date_close", "t1.is_active", "t1.views",
		"t1.publisher_id", "t1.cat_name", "t1.images", "t1.amount", "t1.is_new", "t1.promo_level"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew, testadvert.PromoLevel,
	)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Id, int64(10), int64(1)).WillReturnRows(rows)

	_, err = repo.SelectRecomendations(testadvert.Id, int64(10), int64(1))

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectRecomendationsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"r1.rec_id", "t1.name", "t1.description", "t1.price", "t1.location",
		"t1.latitude", "t1.longitude", "t1.published_at", "t1.date_close", "t1.is_active", "t1.views",
		"t1.publisher_id", "t1.cat_name", "t1.images", "t1.amount", "t1.is_new", "t1.promo_level"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew, testadvert.PromoLevel,
	)
	mock.ExpectQuery("SELECT").WithArgs(testadvert.Id, int64(10), int64(1))

	_, err = repo.SelectRecomendations(testadvert.Id, int64(10), int64(1))

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectDummyRecomendationsOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"r1.rec_id", "t1.name", "t1.description", "t1.price", "t1.location",
		"t1.latitude", "t1.longitude", "t1.published_at", "t1.date_close", "t1.is_active", "t1.views",
		"t1.publisher_id", "t1.cat_name", "t1.images", "t1.amount", "t1.is_new", "t1.promo_level"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew, testadvert.PromoLevel,
	)
	mock.ExpectQuery("SELECT").WithArgs(int64(0), int64(10)).WillReturnRows(rows)

	_, err = repo.SelectDummyRecomendations(int64(0), int64(10))

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSelectDummyRecomendationsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	rows := sqlmock.NewRows([]string{"r1.rec_id", "t1.name", "t1.description", "t1.price", "t1.location",
		"t1.latitude", "t1.longitude", "t1.published_at", "t1.date_close", "t1.is_active", "t1.views",
		"t1.publisher_id", "t1.cat_name", "t1.images", "t1.amount", "t1.is_new", "t1.promo_level"},
	)
	rows.AddRow(testadvert.Id, testadvert.Name, testadvert.Description, testadvert.Price, testadvert.Location, testadvert.Latitude,
		testadvert.Longitude, testadvert.PublishedAt, testadvert.DateClose, testadvert.IsActive, testadvert.Views, testadvert.PublisherId,
		testadvert.Category, testimages, testadvert.Amount, testadvert.IsNew, testadvert.PromoLevel,
	)
	mock.ExpectQuery("SELECT").WithArgs(int64(0), int64(10))

	_, err = repo.SelectDummyRecomendations(int64(0), int64(10))

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestRegenerateRecomendationsOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs().WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()

	err = repo.RegenerateRecomendations()

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestRegenerateRecomendationsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewAdvtRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs()
	mock.ExpectRollback()

	err = repo.RegenerateRecomendations()

	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
