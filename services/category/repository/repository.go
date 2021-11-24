package repository

import (
	"database/sql"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/services/category"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(DB *sql.DB) category.CategoryRepository {
	return &CategoryRepository{
		DB: DB,
	}
}

func (cr *CategoryRepository) SelectCategories() ([]*models.Category, error) {
	rows, err := cr.DB.Query("SELECT name FROM category")
	if err != nil {
		return nil, internalError.GenInternalError(err)
	}

	defer rows.Close()
	categories := make([]*models.Category, 0)
	for rows.Next() {
		var ctgry models.Category

		err = rows.Scan(&ctgry.Name)
		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		categories = append(categories, &ctgry)
	}

	return categories, nil
}
