package repository

import (
	"context"
	"fmt"
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
	queryStr := `SELECT a.id, a.name, a.description, a.price, a.location, a.published_at, a.image, a.publisher_id, a.is_active FROM advts a
				 %s LIMIT $1 OFFSET $2;`
	if isSortedByPublichedDate {
		queryStr = fmt.Sprintf(queryStr, " ORDER BY a.published_at DESC")
	} else {
		queryStr = fmt.Sprintf(queryStr, "")
	}

	rows, err := ar.pool.Query(context.Background(), queryStr, count, from*count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var advts []*models.AdvtData
	for rows.Next() {
		advt := &models.AdvtData{}

		err := rows.Scan(&advt.Id, &advt.Name, &advt.Description, &advt.Price, &advt.Location,
			&advt.PublishedAt, &advt.Image, &advt.PublisherId, &advt.IsActive)
		if err != nil {
			return nil, err
		}

		advts = append(advts, advt)
	}

	return advts, nil
}
