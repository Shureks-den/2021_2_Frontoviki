package repository

import (
	"context"
	"database/sql"
	"regexp"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/user"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(DB *sql.DB) user.UserRepository {
	return &UserRepository{
		DB: DB,
	}
}

func (ur *UserRepository) Insert(user *models.UserData) error {
	tx, err := ur.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	row := tx.QueryRow(`INSERT INTO users (email, password, created_at, name, surname, image, phone) 
						VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`,
		user.Email, user.Password, user.CreatedAt, user.Name, user.Surname, user.Image, user.Phone)

	var id int64
	err = row.Scan(&id)

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return internalError.RollbackError
		}

		return internalError.GenInternalError(err)
	}

	_, err = tx.ExecContext(context.Background(), "INSERT INTO rating_statistics(user_id) VALUES ($1);", id)
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return rollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	user.Id = id
	return nil
}

func (ur *UserRepository) SelectByEmail(email string) (*models.UserData, error) {
	row := ur.DB.QueryRowContext(context.Background(),
		"SELECT id, email, phone, password, created_at, name, surname, image FROM users WHERE email = $1",
		email)

	user := models.UserData{}
	if err := row.Scan(&user.Id, &user.Email, &user.Phone, &user.Password, &user.CreatedAt,
		&user.Name, &user.Surname, &user.Image); err != nil {
		res, _ := regexp.Match(".*no rows.*", []byte(err.Error()))
		if res {
			return nil, internalError.EmptyQuery
		}
		return nil, internalError.GenInternalError(err)
	}

	return &user, nil
}

func (ur *UserRepository) SelectById(userId int64) (*models.UserData, error) {
	row := ur.DB.QueryRowContext(context.Background(),
		"SELECT id, email, phone, password, created_at, name, surname, image FROM users WHERE id = $1",
		userId)
	user := models.UserData{}
	if err := row.Scan(&user.Id, &user.Email, &user.Phone, &user.Password, &user.CreatedAt,
		&user.Name, &user.Surname, &user.Image); err != nil {
		res, _ := regexp.Match(".*no rows.*", []byte(err.Error()))
		if res {
			return nil, internalError.EmptyQuery
		}
		return nil, internalError.GenInternalError(err)
	}

	return &user, nil
}

func (ur *UserRepository) Update(user *models.UserData) error {
	tx, err := ur.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	ct, err := tx.ExecContext(context.Background(),
		"UPDATE users SET email = $2, password = $3, name = $4, surname = $5, image = $6, phone = $7 WHERE id = $1",
		user.Id, user.Email, user.Password, user.Name, user.Surname, user.Image, user.Phone)

	if ra, _ := ct.RowsAffected(); ra != 1 || err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return internalError.RollbackError
		}

		if ra != 1 {
			return internalError.NotUpdated
		}

		return internalError.GenInternalError(err)
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}
