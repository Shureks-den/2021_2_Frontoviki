package repository

import (
	"context"
	"fmt"
	"yula/internal/codes"
	"yula/internal/models"
	"yula/internal/pkg/user"

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

func (ur *UserRepository) Insert(user *models.UserData) *codes.DatabaseError {
	row := ur.pool.QueryRow(context.Background(),
		"INSERT INTO users (username, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id;",
		user.Username, user.Email, user.Password, user.CreatedAt)

	var id int64
	if err := row.Scan(&id); err != nil {
		fmt.Println("unable to insert", err.Error())
		return codes.NewDatabaseError(codes.UnexpectedDbError)
	}

	user.Id = id
	return nil
}

func (ur *UserRepository) SelectByEmail(email string) (*models.UserData, *codes.DatabaseError) {
	row := ur.pool.QueryRow(context.Background(),
		"SELECT id, username, email, password, created_at, name, surname, image FROM users WHERE email = $1", email)

	user := models.UserData{}
	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt,
		&user.Name, &user.Surname, &user.Image); err != nil {
		switch err.Error() {
		case "no rows in result set":
			return nil, codes.NewDatabaseError(codes.EmptyRow)
		}
		return nil, codes.NewDatabaseError(codes.UnexpectedDbError)
	}

	return &user, nil
}

func (ur *UserRepository) SelectById(userId int64) (*models.UserData, *codes.DatabaseError) {
	row := ur.pool.QueryRow(context.Background(),
		"SELECT id, username, email, password, created_at, name, surname, image FROM users WHERE id = $1", userId)

	user := models.UserData{}
	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt,
		&user.Name, &user.Surname, &user.Image); err != nil {
		switch err.Error() {
		case "no rows in result set":
			return nil, codes.NewDatabaseError(codes.EmptyRow)
		}
		return nil, codes.NewDatabaseError(codes.UnexpectedDbError)
	}

	return &user, nil
}

func (ur *UserRepository) Update(user *models.UserData) *codes.DatabaseError {
	res, err := ur.pool.Exec(context.Background(),
		"UPDATE users SET username = $2, email = $3, password = $4, name = $5, surname = $6, image = $7 WHERE id = $1",
		user.Id, user.Username, user.Email, user.Password, user.Name, user.Surname, user.Image)

	if err != nil {
		return codes.NewDatabaseError(codes.UnableToUpdate)
	}

	if count := res.RowsAffected(); count != 1 {
		return codes.NewDatabaseError(codes.NotUpdated)
	}

	return nil
}
