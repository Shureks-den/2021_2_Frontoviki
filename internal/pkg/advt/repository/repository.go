package repository

import (
	"context"
	"fmt"
	"sync"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/advt"
	imageloader "yula/internal/pkg/image_loader"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AdvtRepository struct {
	pool *pgxpool.Pool
	m    sync.RWMutex
}

func NewAdvtRepository(pool *pgxpool.Pool) advt.AdvtRepository {
	return &AdvtRepository{
		pool: pool,
		m:    sync.RWMutex{},
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

	rows, err := ar.pool.Query(context.Background(), queryStr, count, from*count)
	if err != nil {
		return nil, internalError.GenInternalError(err)
	}
	defer rows.Close()

	var adverts []*models.Advert
	var advertPathImages []*string
	for rows.Next() {
		advert := &models.Advert{}

		err := rows.Scan(&advert.Id, &advert.Name, &advert.Description, &advert.Price, &advert.Location, &advert.Latitude,
			&advert.Longitude, &advert.PublishedAt, &advert.DateClose, &advert.IsActive, &advert.Views,
			&advert.PublisherId, &advert.Category, &advertPathImages, &advert.Amount, &advert.IsNew)
		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		advert.Images = []string{}
		for _, path := range advertPathImages {
			if path != nil {
				advert.Images = append(advert.Images, *path)
			}
		}

		if len(advert.Images) == 0 {
			advert.Images = append(advert.Images, imageloader.DefaultAdvertImage)
		}

		adverts = append(adverts, advert)
	}

	return adverts, nil
}

func (ar *AdvtRepository) Insert(advert *models.Advert) error {
	tx, err := ar.pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return internalError.GenInternalError(err)
	}

	queryStr := `INSERT INTO advert (name, description, category_id, publisher_id, latitude, longitude, location, price, amount, is_new) 
				VALUES ($1, $2, (SELECT id FROM category WHERE name = $3), $4, $5, $6, $7, $8, $9, $10) RETURNING id;`
	query := ar.pool.QueryRow(context.Background(), queryStr,
		advert.Name, advert.Description, advert.Category, advert.PublisherId,
		advert.Latitude, advert.Longitude, advert.Location, advert.Price, advert.Amount, advert.IsNew)

	if err := query.Scan(&advert.Id); err != nil {
		fmt.Println(err.Error())
		rollbackErr := tx.Rollback(context.Background())
		if rollbackErr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit(context.Background())
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
	queryRow := ar.pool.QueryRow(context.Background(), queryStr, advertId)

	var advert models.Advert
	var advertPathImages []*string

	err := queryRow.Scan(&advert.Id, &advert.Name, &advert.Description, &advert.Price, &advert.Location, &advert.Latitude,
		&advert.Longitude, &advert.PublishedAt, &advert.DateClose, &advert.IsActive, &advert.Views,
		&advert.PublisherId, &advert.Category, &advertPathImages, &advert.Amount, &advert.IsNew)

	if err != nil {
		return nil, internalError.EmptyQuery
	}

	advert.Images = []string{}
	for _, path := range advertPathImages {
		if path != nil {
			advert.Images = append(advert.Images, *path)
		}
	}

	return &advert, nil
}

func (ar *AdvtRepository) Update(newAdvert *models.Advert) error {
	tx, err := ar.pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return internalError.GenInternalError(err)
	}

	queryStr := `UPDATE advert set name = $2, description = $3, category_id = (SELECT c.id FROM category c WHERE c.name = $4), 
				location = $5, latitude = $6, longitude = $7, price = $8, is_active = $9, date_close = $10, 
				amount = $11, is_new = $12 WHERE id = $1 RETURNING id;`
	query := tx.QueryRow(context.Background(), queryStr, newAdvert.Id, newAdvert.Name, newAdvert.Description,
		newAdvert.Category, newAdvert.Location, newAdvert.Latitude, newAdvert.Longitude,
		newAdvert.Price, newAdvert.IsActive, newAdvert.DateClose, newAdvert.Amount, newAdvert.IsNew)

	err = query.Scan(&newAdvert.Id)
	if err != nil {
		if rlbckEr := tx.Rollback(context.Background()); rlbckEr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return internalError.NotCommited
	}
	return nil
}

func (ar *AdvtRepository) Delete(advertId int64) error {
	tx, err := ar.pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = tx.Exec(context.Background(), "DELETE FROM advert WHERE id = $1;", advertId)
	if err != nil {
		if rlbckEr := tx.Rollback(context.Background()); rlbckEr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}

func (ar *AdvtRepository) EditImages(advertId int64, newImages []string) error {
	tx, err := ar.pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return internalError.GenInternalError(err)
	}

	// сначала очищаем все картинки у объявления
	_, err = tx.Exec(context.Background(),
		"DELETE FROM advert_image WHERE advert_id = $1;",
		advertId)
	if err != nil {
		rollbackErr := tx.Rollback(context.Background())
		if rollbackErr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	// вставляем в базу новые url картинок
	for _, image := range newImages {
		_, err := tx.Exec(context.Background(),
			"INSERT INTO advert_image (advert_id, img_path) VALUES ($1, $2);",
			advertId, image)
		if err != nil {
			rollbackErr := tx.Rollback(context.Background())
			if rollbackErr != nil {
				return internalError.RollbackError
			}
			return internalError.GenInternalError(err)
		}
	}

	err = tx.Commit(context.Background())
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

	rows, err := ar.pool.Query(context.Background(), queryStr, publisherId, limit, offset*limit)
	if err != nil {
		return nil, internalError.GenInternalError(err)
	}
	defer rows.Close()

	var adverts []*models.Advert
	for rows.Next() {
		var advert models.Advert
		var advertPathImages []*string

		err := rows.Scan(&advert.Id, &advert.Name, &advert.Description, &advert.Price, &advert.Location, &advert.Latitude,
			&advert.Longitude, &advert.PublishedAt, &advert.DateClose, &advert.IsActive, &advert.Views,
			&advert.PublisherId, &advert.Category, &advertPathImages, &advert.Amount, &advert.IsNew)

		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		advert.Images = []string{}
		for _, path := range advertPathImages {
			if path != nil {
				advert.Images = append(advert.Images, *path)
			}
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
	query, err := ar.pool.Query(context.Background(), queryStr, categoryName, count, from*count)
	if err != nil {
		fmt.Println(err.Error())
		return nil, internalError.InternalError
	}

	defer query.Close()
	adverts := make([]*models.Advert, 0)
	for query.Next() {
		var advert models.Advert
		var advertPathImages []*string

		err = query.Scan(&advert.Id, &advert.Name, &advert.Description, &advert.Price, &advert.Location, &advert.Latitude,
			&advert.Longitude, &advert.PublishedAt, &advert.DateClose, &advert.IsActive, &advert.Views,
			&advert.PublisherId, &advert.Category, &advertPathImages, &advert.Amount, &advert.IsNew)

		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		advert.Images = []string{}
		for _, path := range advertPathImages {
			if path != nil {
				advert.Images = append(advert.Images, *path)
			}
		}

		if len(advert.Images) == 0 {
			advert.Images = append(advert.Images, imageloader.DefaultAdvertImage)
		}

		adverts = append(adverts, &advert)
	}
	return adverts, nil
}
