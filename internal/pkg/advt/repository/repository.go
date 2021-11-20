package repository

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"sort"
	"strings"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/advt"
	imageloader "yula/internal/pkg/image_loader"
)

type AdvtRepository struct {
	DB *sql.DB
}

func NewAdvtRepository(DB *sql.DB) advt.AdvtRepository {
	return &AdvtRepository{
		DB: DB,
	}
}

func (ar *AdvtRepository) SelectListAdvt(isSortedByPublichedDate bool, from, count int64) ([]*models.Advert, error) {
	queryStr := `SELECT a.id, a.Name, a.Description, a.price, a.location, a.latitude, a.longitude, a.published_at, 
				 a.date_close, a.is_active, a.views, a.publisher_id, c.name, array_agg(ai.img_path), a.amount, a.is_new FROM advert a
				 JOIN category c ON a.category_id = c.Id
				 LEFT JOIN advert_image ai ON a.id = ai.advert_id
				 WHERE a.is_active 
				 GROUP BY a.id, a.name, a.Description,  a.price, a.location, a.latitude, a.longitude, a.published_at, 
				 a.date_close, a.is_active, a.views, a.publisher_id, c.name %s LIMIT $1 OFFSET $2;`
	if isSortedByPublichedDate {
		queryStr = fmt.Sprintf(queryStr, " ORDER BY a.published_at DESC")
	} else {
		queryStr = fmt.Sprintf(queryStr, "")
	}

	rows, err := ar.DB.QueryContext(context.Background(), queryStr, count, from*count)
	if err != nil {
		return nil, internalError.GenInternalError(err)
	}
	defer rows.Close()

	adverts := make([]*models.Advert, 0)
	for rows.Next() {
		advert := &models.Advert{}
		var images string

		err := rows.Scan(&advert.Id, &advert.Name, &advert.Description, &advert.Price, &advert.Location, &advert.Latitude,
			&advert.Longitude, &advert.PublishedAt, &advert.DateClose, &advert.IsActive, &advert.Views,
			&advert.PublisherId, &advert.Category, &images, &advert.Amount, &advert.IsNew)
		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		advert.Images = make([]string, 0)
		if images[1:len(images)-1] != "NULL" {
			advert.Images = strings.Split(images[1:len(images)-1], ",")
			sort.Strings(advert.Images)
		}

		if len(advert.Images) == 0 {
			advert.Images = append(advert.Images, imageloader.DefaultAdvertImage)
		}

		adverts = append(adverts, advert)
	}

	return adverts, nil
}

func (ar *AdvtRepository) Insert(advert *models.Advert) error {
	tx, err := ar.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	queryStr := `INSERT INTO advert (name, description, category_id, publisher_id, latitude, longitude, location, price, amount, is_new) 
				VALUES ($1, $2, (SELECT id FROM category WHERE lower(name) = lower($3)), $4, $5, $6, $7, $8, $9, $10) RETURNING id;`
	query := ar.DB.QueryRowContext(context.Background(), queryStr,
		advert.Name, advert.Description, advert.Category, advert.PublisherId,
		advert.Latitude, advert.Longitude, advert.Location, advert.Price, advert.Amount, advert.IsNew)

	if err := query.Scan(&advert.Id); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	// вставляем в таблицу просмотров id созданного объявления
	queryStr = "INSERT INTO views_ (advert_id) VALUES ($1);"
	_, err = ar.DB.Exec(queryStr, advert.Id)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}

func (ar *AdvtRepository) SelectById(advertId int64) (*models.Advert, error) {
	queryStr := `
				SELECT a.id, a.Name, a.Description, a.price, a.location, a.latitude, a.longitude, a.published_at, 
				a.date_close, a.is_active, a.views, a.publisher_id, c.name, array_agg(ai.img_path), a.amount, a.is_new FROM advert a
				JOIN category c ON a.category_id = c.Id
				LEFT JOIN advert_image ai ON a.id = ai.advert_id 
				GROUP BY a.id, a.name, a.Description,  a.price, a.location, a.latitude, a.longitude, a.published_at, 
				a.date_close, a.is_active, a.views, a.publisher_id, c.name
				HAVING a.id = $1;`
	queryRow := ar.DB.QueryRowContext(context.Background(), queryStr, advertId)

	var advert models.Advert
	var images string

	err := queryRow.Scan(&advert.Id, &advert.Name, &advert.Description, &advert.Price, &advert.Location, &advert.Latitude,
		&advert.Longitude, &advert.PublishedAt, &advert.DateClose, &advert.IsActive, &advert.Views,
		&advert.PublisherId, &advert.Category, &images, &advert.Amount, &advert.IsNew)

	if err != nil {
		return nil, internalError.EmptyQuery
	}

	advert.Images = make([]string, 0)
	if images[1:len(images)-1] != "NULL" {
		advert.Images = strings.Split(images[1:len(images)-1], ",")
		sort.Strings(advert.Images)
	}

	if len(advert.Images) == 0 {
		advert.Images = append(advert.Images, imageloader.DefaultAdvertImage)
	}

	return &advert, nil
}

func (ar *AdvtRepository) Update(newAdvert *models.Advert) error {
	tx, err := ar.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	queryStr := `UPDATE advert set name = $2, description = $3, category_id = (SELECT c.id FROM category c WHERE lower(c.name) = lower($4)), 
				location = $5, latitude = $6, longitude = $7, price = $8, is_active = $9, date_close = $10, 
				amount = $11, is_new = $12 WHERE id = $1 RETURNING id;`
	query := tx.QueryRowContext(context.Background(), queryStr, newAdvert.Id, newAdvert.Name, newAdvert.Description,
		newAdvert.Category, newAdvert.Location, newAdvert.Latitude, newAdvert.Longitude,
		newAdvert.Price, newAdvert.IsActive, newAdvert.DateClose, newAdvert.Amount, newAdvert.IsNew)

	err = query.Scan(&newAdvert.Id)
	if err != nil {
		if rlbckEr := tx.Rollback(); rlbckEr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}
	return nil
}

func (ar *AdvtRepository) Delete(advertId int64) error {
	tx, err := ar.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = tx.ExecContext(context.Background(), "DELETE FROM advert WHERE id = $1;", advertId)
	if err != nil {
		if rlbckEr := tx.Rollback(); rlbckEr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}

func (ar *AdvtRepository) DeleteImages(images []string, advertId int64) error {
	tx, err := ar.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	for _, img_path := range images {
		_, err = tx.ExecContext(context.Background(),
			"DELETE FROM advert_image WHERE advert_id = $1 AND img_path LIKE $2;",
			advertId, img_path)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return internalError.RollbackError
			}
			return internalError.GenInternalError(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}

func (ar *AdvtRepository) InsertImages(advertId int64, newImages []string) error {
	tx, err := ar.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	// сначала очищаем все картинки у объявления
	// _, err = tx.ExecContext(context.Background(),
	// 	"DELETE FROM advert_image WHERE advert_id = $1;",
	// 	advertId)
	// if err != nil {
	// 	rollbackErr := tx.Rollback()
	// 	if rollbackErr != nil {
	// 		return internalError.RollbackError
	// 	}
	// 	return internalError.GenInternalError(err)
	// }

	// вставляем в базу новые url картинок
	for _, image := range newImages {
		_, err := tx.ExecContext(context.Background(),
			"INSERT INTO advert_image (advert_id, img_path) VALUES ($1, $2);",
			advertId, image)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return internalError.RollbackError
			}
			return internalError.GenInternalError(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}

const (
	defaultAdvertsQueryByPublisherId string = `
		SELECT a.id, a.Name, a.Description, a.price, a.location, a.latitude, a.longitude, a.published_at, 
			a.date_close, a.is_active, a.views, a.publisher_id, c.name, array_agg(ai.img_path), a.amount, a.is_new FROM advert a
		JOIN category c ON a.category_id = c.Id 
		LEFT JOIN advert_image ai ON a.id = ai.advert_id
		GROUP BY a.id, a.name, a.Description,  a.price, a.location, a.latitude, a.longitude, a.published_at, 
			a.date_close, a.is_active, a.views, a.publisher_id, c.name, a.amount, a.is_new
		HAVING a.publisher_id = $1 %s %s 
		LIMIT $2 OFFSET $3;
	`
)

func (ar *AdvtRepository) SelectAdvertsByPublisherId(publisherId int64, is_active bool, offset int64, limit int64) ([]*models.Advert, error) {
	var queryStr string
	if is_active {
		queryStr = fmt.Sprintf(defaultAdvertsQueryByPublisherId, "AND a.is_active = true",
			"ORDER BY a.is_active DESC, a.published_at DESC",
		)
	} else {
		queryStr = fmt.Sprintf(defaultAdvertsQueryByPublisherId, "AND a.is_active = false",
			"ORDER BY a.is_active DESC, a.published_at DESC",
		)
	}

	rows, err := ar.DB.QueryContext(context.Background(), queryStr, publisherId, limit, offset*limit)
	if err != nil {
		return nil, internalError.GenInternalError(err)
	}
	defer rows.Close()

	var adverts []*models.Advert
	for rows.Next() {
		var advert models.Advert
		var images string

		err := rows.Scan(&advert.Id, &advert.Name, &advert.Description, &advert.Price, &advert.Location, &advert.Latitude,
			&advert.Longitude, &advert.PublishedAt, &advert.DateClose, &advert.IsActive, &advert.Views,
			&advert.PublisherId, &advert.Category, &images, &advert.Amount, &advert.IsNew)

		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		advert.Images = make([]string, 0)
		if images[1:len(images)-1] != "NULL" {
			advert.Images = strings.Split(images[1:len(images)-1], ",")
			sort.Strings(advert.Images)
		}

		if len(advert.Images) == 0 {
			advert.Images = append(advert.Images, imageloader.DefaultAdvertImage)
		}

		adverts = append(adverts, &advert)
	}

	return adverts, nil
}

func (ar *AdvtRepository) SelectAdvertsByCategory(categoryName string, from, count int64) ([]*models.Advert, error) {
	queryStr := `
		SELECT a.id, a.Name, a.Description, a.price, a.location, a.latitude, a.longitude, a.published_at, 
			a.date_close, a.is_active, a.views, a.publisher_id, c.name, array_agg(ai.img_path), a.amount, a.is_new 
		FROM (
			SELECT * FROM advert WHERE category_id = (SELECT id FROM category WHERE lower(name) = lower($1))
		) as a 
		JOIN category c ON a.category_id = c.Id
		LEFT JOIN advert_image ai ON a.id = ai.advert_id
		GROUP BY a.id, a.name, a.Description,  a.price, a.location, a.latitude, a.longitude, a.published_at, 
				a.date_close, a.is_active, a.views, a.publisher_id, c.name, a.amount, a.is_new 
		HAVING a.is_active = true
		LIMIT $2 OFFSET $3;
	`
	query, err := ar.DB.QueryContext(context.Background(), queryStr, categoryName, count, from*count)
	if err != nil {
		fmt.Println(err.Error())
		return nil, internalError.InternalError
	}

	defer query.Close()
	adverts := make([]*models.Advert, 0)
	for query.Next() {
		var advert models.Advert
		var images string

		err = query.Scan(&advert.Id, &advert.Name, &advert.Description, &advert.Price, &advert.Location, &advert.Latitude,
			&advert.Longitude, &advert.PublishedAt, &advert.DateClose, &advert.IsActive, &advert.Views,
			&advert.PublisherId, &advert.Category, &images, &advert.Amount, &advert.IsNew)

		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		advert.Images = make([]string, 0)
		if images[1:len(images)-1] != "NULL" {
			advert.Images = strings.Split(images[1:len(images)-1], ",")
			sort.Strings(advert.Images)
		}

		if len(advert.Images) == 0 {
			advert.Images = append(advert.Images, imageloader.DefaultAdvertImage)
		}

		adverts = append(adverts, &advert)
	}
	return adverts, nil
}

func (ar *AdvtRepository) SelectFavoriteAdverts(userId int64, from, count int64) ([]*models.Advert, error) {
	queryStr := `
		SELECT a.id, a.Name, a.Description, a.price, a.location, a.latitude, a.longitude, a.published_at, 
			a.date_close, a.is_active, a.views, a.publisher_id, c.name, array_agg(ai.img_path), a.amount, a.is_new FROM advert a
		JOIN favorite f ON a.id = f.advert_id
		JOIN category c ON a.category_id = c.Id 
		LEFT JOIN advert_image ai ON a.id = ai.advert_id
		WHERE f.user_id = $1
		GROUP BY a.id, a.name, a.Description,  a.price, a.location, a.latitude, a.longitude, a.published_at, 
			a.date_close, a.is_active, a.views, a.publisher_id, c.name, a.amount, a.is_new
		LIMIT $2 OFFSET $3;
	`
	query, err := ar.DB.QueryContext(context.Background(), queryStr, userId, count, from*count)
	if err != nil {
		return nil, internalError.GenInternalError(err)
	}

	defer query.Close()
	adverts := make([]*models.Advert, 0)
	for query.Next() {
		var advert models.Advert
		var images string

		err = query.Scan(&advert.Id, &advert.Name, &advert.Description, &advert.Price, &advert.Location, &advert.Latitude,
			&advert.Longitude, &advert.PublishedAt, &advert.DateClose, &advert.IsActive, &advert.Views,
			&advert.PublisherId, &advert.Category, &images, &advert.Amount, &advert.IsNew)

		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		advert.Images = make([]string, 0)
		if images[1:len(images)-1] != "NULL" {
			advert.Images = strings.Split(images[1:len(images)-1], ",")
			sort.Strings(advert.Images)
		}

		if len(advert.Images) == 0 {
			advert.Images = append(advert.Images, imageloader.DefaultAdvertImage)
		}

		adverts = append(adverts, &advert)
	}
	return adverts, nil
}

func (ar *AdvtRepository) SelectFavorite(userId, advertId int64) (*models.Advert, error) {
	queryStr := `
		SELECT a.id, a.Name, a.Description, a.price, a.location, a.latitude, a.longitude, a.published_at, 
			a.date_close, a.is_active, a.views, a.publisher_id, c.name, array_agg(ai.img_path), a.amount, a.is_new FROM advert a
		JOIN category c ON a.category_id = c.Id
		LEFT JOIN advert_image ai ON a.id = ai.advert_id 
		JOIN favorite f ON a.id = f.advert_id AND f.user_id = $2
		GROUP BY a.id, a.name, a.Description,  a.price, a.location, a.latitude, a.longitude, a.published_at, 
			a.date_close, a.is_active, a.views, a.publisher_id, c.name
		HAVING a.id = $1;
	`
	queryRow := ar.DB.QueryRowContext(context.Background(), queryStr, advertId, userId)

	var advert models.Advert
	var images string

	err := queryRow.Scan(&advert.Id, &advert.Name, &advert.Description, &advert.Price, &advert.Location, &advert.Latitude,
		&advert.Longitude, &advert.PublishedAt, &advert.DateClose, &advert.IsActive, &advert.Views,
		&advert.PublisherId, &advert.Category, &images, &advert.Amount, &advert.IsNew)

	if err != nil {
		res, _ := regexp.Match(".*no rows.*", []byte(err.Error()))
		if res {
			return nil, internalError.EmptyQuery
		} else {
			return nil, internalError.GenInternalError(err)
		}
	}

	advert.Images = make([]string, 0)
	if images[1:len(images)-1] != "NULL" {
		advert.Images = strings.Split(images[1:len(images)-1], ",")
		sort.Strings(advert.Images)
	}

	if len(advert.Images) == 0 {
		advert.Images = append(advert.Images, imageloader.DefaultAdvertImage)
	}

	return &advert, nil
}

func (ar *AdvtRepository) InsertFavorite(userId, advertId int64) error {
	tx, err := ar.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = tx.ExecContext(context.Background(),
		"INSERT INTO favorite(user_id, advert_id) VALUES ($1, $2);",
		userId, advertId)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}

func (ar *AdvtRepository) DeleteFavorite(userId, advertId int64) error {
	tx, err := ar.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = tx.ExecContext(context.Background(),
		"DELETE FROM favorite WHERE user_id = $1 AND advert_id = $2;",
		userId, advertId)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}

func (ar *AdvtRepository) SelectViews(advertId int64) (int64, error) {
	queryStr := "SELECT count FROM views_ WHERE advert_id = $1;"
	queryRow := ar.DB.QueryRowContext(context.Background(), queryStr, advertId)

	var views int64
	err := queryRow.Scan(&views)
	return views, err
}

func (ar *AdvtRepository) UpdateViews(advertId int64) error {
	tx, err := ar.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = tx.ExecContext(context.Background(),
		"UPDATE views_ SET count = count + 1 WHERE advert_id = $1;", advertId)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}
