package repository

import (
	"context"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/category"

	"github.com/jackc/pgx/v4/pgxpool"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) category.CategoryRepository {
	return &CategoryRepository{
		pool: pool,
	}
}

func (cr *CategoryRepository) SelectCategories() ([]*models.Category, error) {
	queryStr := "SELECT name FROM category;"
	query, err := cr.pool.Query(context.Background(), queryStr)
	if err != nil {
		return nil, internalError.InternalError
	}

	defer query.Close()
	categories := make([]*models.Category, 0)
	for query.Next() {
		var ctgry models.Category

		err = query.Scan(&ctgry.Name)
		if err != nil {
			return nil, internalError.InternalError
		}

		categories = append(categories, &ctgry)
	}

	return categories, nil
}
