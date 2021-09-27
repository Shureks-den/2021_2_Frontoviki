package repository

import (
	"context"
	"log"
	"sync"
	"yula/internal/models"
	"yula/internal/pkg/advt"

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

func (ar *AdvtRepository) SelectListAdvt(isSortedByPublichedDate bool, from, count int64) ([]*models.AdvtData, error) {
	queryStr := `SELECT * FROM advts;`
	log.Println(queryStr)

	rows, err := ar.pool.Query(context.Background(), queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var advts []*models.AdvtData
	log.Println("db prefor")
	for rows.Next() {
		advt := &models.AdvtData{}

		err := rows.Scan(&advt.Id, &advt.Name, &advt.Description, &advt.Price, &advt.Location,
			&advt.PublishedAt, &advt.Image, &advt.PublisherId, &advt.IsActive)
		if err != nil {
			return nil, err
		}
		log.Println("db for")
		advts = append(advts, advt)
	}
	log.Println("db afterfor")
	if rows.Err() != nil {
		log.Println(rows.Err().Error())
		return nil, rows.Err()
	}
	return advts, nil
}
