package repository

import (
	"context"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/cart"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CartRepository struct {
	pool *pgxpool.Pool
}

func NewCartRepository(pool *pgxpool.Pool) cart.CartRepository {
	return &CartRepository{
		pool: pool,
	}
}

func (cr *CartRepository) Select(userId int64, advertId int64) (*models.Cart, error) {
	queryStr := "SELECT user_id, advert_id, amount FROM cart WHERE user_id = $1 AND advert_id = $2;"
	query := cr.pool.QueryRow(context.Background(), queryStr, userId, advertId)
	var oneInCart models.Cart
	err := query.Scan(&oneInCart.UserId, &oneInCart.AdvertId, &oneInCart.Amount)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return nil, internalError.EmptyQuery

		default:
			return nil, internalError.InternalError
		}
	}
	return &oneInCart, nil
}

func (cr *CartRepository) SelectAll(userId int64) ([]*models.Cart, error) {
	queryStr := "SELECT user_id, advert_id, amount FROM cart WHERE user_id = $1;"
	query, err := cr.pool.Query(context.Background(), queryStr, userId)
	if err != nil {
		return nil, internalError.InternalError
	}

	defer query.Close()
	cart := make([]*models.Cart, 0)
	for query.Next() {
		var oneInCart models.Cart

		err = query.Scan(&oneInCart.UserId, &oneInCart.AdvertId, &oneInCart.Amount)
		if err != nil {
			return nil, internalError.InternalError
		}

		cart = append(cart, &oneInCart)
	}

	return cart, nil
}

func (cr *CartRepository) Update(cart *models.Cart) error {
	tx, err := cr.pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return internalError.InternalError
	}

	queryStr := "UPDATE cart SET amount = $3 WHERE user_id = $1 AND advert_id = $2;"
	ct, err := tx.Exec(context.Background(), queryStr, cart.UserId, cart.AdvertId, cart.Amount)

	if ct.RowsAffected() == 0 || err != nil {
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

func (cr *CartRepository) Insert(cart *models.Cart) error {
	tx, err := cr.pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return internalError.InternalError
	}

	queryStr := "INSERT INTO cart (user_id, advert_id, amount) VALUES ($1, $2, $3);"
	ct, err := tx.Exec(context.Background(), queryStr, cart.UserId, cart.AdvertId, cart.Amount)

	if ct.RowsAffected() == 0 || err != nil {
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

func (cr *CartRepository) Delete(cart *models.Cart) error {
	tx, err := cr.pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return internalError.InternalError
	}

	queryStr := "DELETE FROM cart WHERE user_id = $1 AND advert_id = $2;"
	ct, err := tx.Exec(context.Background(), queryStr, cart.UserId, cart.AdvertId)

	if ct.RowsAffected() == 0 || err != nil {
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

func (cr *CartRepository) DeleteAll(userId int64) error {
	tx, err := cr.pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return internalError.InternalError
	}

	queryStr := "DELETE FROM cart WHERE user_id = $1;"
	ct, err := tx.Exec(context.Background(), queryStr, userId)

	if ct.RowsAffected() == 0 || err != nil {
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
