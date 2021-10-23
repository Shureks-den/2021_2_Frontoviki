package repository

import (
	"context"
	"log"
	"sync"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/user"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
	m    sync.RWMutex
}

func NewUserRepository(pool *pgxpool.Pool) user.UserRepository {
	return &UserRepository{
		pool: pool,
		m:    sync.RWMutex{},
	}
}

func (ur *UserRepository) Insert(user *models.UserData) error {
	ur.m.Lock()
	log.Println(user)
	row := ur.pool.QueryRow(context.Background(),
		"INSERT INTO users (email, password, created_at, name, surname, image) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;",
		user.Email, user.Password, user.CreatedAt, user.Name, user.Surname, user.Image)
	ur.m.Unlock()

	var id int64
	if err := row.Scan(&id); err != nil {
		return internalError.DatabaseError
	}

	user.Id = id
	return nil
}

func (ur *UserRepository) SelectByEmail(email string) (*models.UserData, error) {
	ur.m.RLock()
	row := ur.pool.QueryRow(context.Background(),
		"SELECT id, email, password, created_at, name, surname, image, rating FROM users WHERE email = $1", email)
	ur.m.RUnlock()

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
	ur.m.RLock()
	row := ur.pool.QueryRow(context.Background(),
		"SELECT id, email, password, created_at, name, surname, image, rating FROM users WHERE id = $1", userId)
	ur.m.RUnlock()

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
	res, err := ur.pool.Exec(context.Background(),
		"UPDATE users SET email = $2, password = $3, name = $4, surname = $5, image = $6 WHERE id = $1",
		user.Id, user.Email, user.Password, user.Name, user.Surname, user.Image)

	if err != nil {
		return internalError.InvalidQuery
	}

	if count := res.RowsAffected(); count != 1 {
		return internalError.NotUpdated
	}

	return nil
}
