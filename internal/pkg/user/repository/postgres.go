package repository

import (
	"context"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/user"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) user.UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (ur *UserRepository) Insert(user *models.UserData) error {
	tx, err := ur.pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return internalError.InternalError
	}

	row := tx.QueryRow(context.Background(),
		"INSERT INTO users (email, password, created_at, name, surname, image) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;",
		user.Email, user.Password, user.CreatedAt, user.Name, user.Surname, user.Image)

	var id int64
	err = row.Scan(&id)

	if err != nil {
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

	user.Id = id
	return nil
}

func (ur *UserRepository) SelectByEmail(email string) (*models.UserData, error) {
	row := ur.pool.QueryRow(context.Background(),
		"SELECT id, email, password, created_at, name, surname, image, rating FROM users WHERE email = $1", email)

	user := models.UserData{}
	if err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt,
		&user.Name, &user.Surname, &user.Image, &user.Rating); err != nil {
		switch err.Error() {
		case "no rows in result set":
			return nil, internalError.EmptyQuery
		}
		return nil, internalError.DatabaseError
	}

	return &user, nil
}

func (ur *UserRepository) SelectById(userId int64) (*models.UserData, error) {
	row := ur.pool.QueryRow(context.Background(),
		"SELECT id, email, password, created_at, name, surname, image, rating FROM users WHERE id = $1", userId)

	user := models.UserData{}
	if err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt,
		&user.Name, &user.Surname, &user.Image, &user.Rating); err != nil {
		switch err.Error() {
		case "no rows in result set":
			return nil, internalError.EmptyQuery
		}
		return nil, internalError.DatabaseError
	}

	return &user, nil
}

func (ur *UserRepository) Update(user *models.UserData) error {
	tx, err := ur.pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return internalError.InternalError
	}

	ct, err := tx.Exec(context.Background(),
		"UPDATE users SET email = $2, password = $3, name = $4, surname = $5, image = $6 WHERE id = $1",
		user.Id, user.Email, user.Password, user.Name, user.Surname, user.Image)

	if ra := ct.RowsAffected(); ra != 1 || err != nil {
		rollbackErr := tx.Rollback(context.Background())
		if rollbackErr != nil {
			return internalError.RollbackError
		}

		if ra != 1 {
			return internalError.NotUpdated
		}

		return internalError.DatabaseError
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}
