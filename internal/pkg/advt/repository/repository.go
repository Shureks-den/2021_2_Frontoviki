package repository

import (
	"context"
	"fmt"
	"log"
	"sync"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/advt"

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
	queryStr := `SELECT a.id, a.name, a.description, a.price, a.location, a.published_at, a.publisher_id, a.is_active FROM advert a
				 %s LIMIT $1 OFFSET $2;`
	if isSortedByPublichedDate {
		queryStr = fmt.Sprintf(queryStr, " ORDER BY a.published_at DESC")
	} else {
		queryStr = fmt.Sprintf(queryStr, "")
	}

	rows, err := ar.pool.Query(context.Background(), queryStr, count, from*count)
	if err != nil {
		return nil, internalError.DatabaseError
	}
	defer rows.Close()

	var advts []*models.Advert
	for rows.Next() {
		advt := &models.Advert{}

		err := rows.Scan(&advt.Id, &advt.Name, &advt.Description, &advt.Price, &advt.Location,
			&advt.PublishedAt, &advt.PublisherId, &advt.IsActive)
		if err != nil {
			return nil, internalError.DatabaseError
		}

		advts = append(advts, advt)
	}

	return advts, nil
}

func (ar *AdvtRepository) Insert(advert *models.Advert) error {
	tx, err := ar.pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return internalError.DatabaseError
	}

	queryStr := `INSERT INTO advert (name, description, category_id, publisher_id, latitude, longitude, location, price) 
				VALUES ($1, $2, (SELECT id FROM category WHERE name = $3), $4, $5, $6, $7, $8) RETURNING id;`
	query := ar.pool.QueryRow(context.Background(), queryStr,
		advert.Name, advert.Description, advert.Category, advert.PublisherId,
		advert.Latitude, advert.Longitude, advert.Location, advert.Price)

	if err := query.Scan(&advert.Id); err != nil {
		rollbackErr := tx.Rollback(context.Background())
		if rollbackErr != nil {
			return internalError.RollbackError
		}
		return internalError.DatabaseError
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
				a.date_close, a.is_active, a.views, a.publisher_id, c.name, array_agg(ai.img_path) FROM advert a
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
		&advert.PublisherId, &advert.Category, &advertPathImages)

	if err != nil {
		log.Println(err.Error())
		return nil, internalError.DatabaseError
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
		return internalError.DatabaseError
	}

	queryStr := `UPDATE advert set name = $2, description = $3, category_id = (SELECT c.id FROM category c WHERE c.name = $4), 
				location = $5, latitude = $6, longitude = $7, price = $8, is_active = $9 WHERE id = $1 RETURNING id;`
	query := tx.QueryRow(context.Background(), queryStr, newAdvert.Id, newAdvert.Name, newAdvert.Description,
		newAdvert.Category, newAdvert.Location, newAdvert.Latitude, newAdvert.Longitude,
		newAdvert.Price, newAdvert.IsActive)

	err = query.Scan(&newAdvert.Id)
	if err != nil {
		if rlbckEr := tx.Rollback(context.Background()); rlbckEr != nil {
			return internalError.RollbackError
		}
		return internalError.DatabaseError
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
		return internalError.DatabaseError
	}

	_, err = tx.Exec(context.Background(), "DELETE FROM advert WHERE id = $1;", advertId)
	if err != nil {
		if rlbckEr := tx.Rollback(context.Background()); rlbckEr != nil {
			return internalError.RollbackError
		}
		return internalError.DatabaseError
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}
