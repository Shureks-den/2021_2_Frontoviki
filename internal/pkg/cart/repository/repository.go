package repository

import (
	"context"
	"database/sql"
	"regexp"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/cart"
)

type CartRepository struct {
	DB *sql.DB
}

func NewCartRepository(DB *sql.DB) cart.CartRepository {
	return &CartRepository{
		DB: DB,
	}
}

func (cr *CartRepository) Select(userId int64, advertId int64) (*models.Cart, error) {
	queryStr := "SELECT user_id, advert_id, amount FROM cart WHERE user_id = $1 AND advert_id = $2;"
	query := cr.DB.QueryRowContext(context.Background(), queryStr, userId, advertId)
	var oneInCart models.Cart
	err := query.Scan(&oneInCart.UserId, &oneInCart.AdvertId, &oneInCart.Amount)
	if err != nil {
		res, _ := regexp.Match(".*no rows.*", []byte(err.Error()))
		if res {
			return nil, internalError.EmptyQuery
		} else {
			return nil, internalError.GenInternalError(err)
		}
	}
	return &oneInCart, nil
}

func (cr *CartRepository) SelectAll(userId int64) ([]*models.Cart, error) {
	queryStr := "SELECT user_id, advert_id, amount FROM cart WHERE user_id = $1;"
	query, err := cr.DB.QueryContext(context.Background(), queryStr, userId)
	if err != nil {
		return nil, internalError.GenInternalError(err)
	}

	defer query.Close()
	cart := make([]*models.Cart, 0)
	for query.Next() {
		var oneInCart models.Cart

		err = query.Scan(&oneInCart.UserId, &oneInCart.AdvertId, &oneInCart.Amount)
		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		cart = append(cart, &oneInCart)
	}

	return cart, nil
}

func (cr *CartRepository) Update(cart *models.Cart) error {
	tx, err := cr.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	queryStr := "UPDATE cart SET amount = $3 WHERE user_id = $1 AND advert_id = $2;"
	_, err = tx.ExecContext(context.Background(), queryStr, cart.UserId, cart.AdvertId, cart.Amount)

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}

func (cr *CartRepository) Insert(cart *models.Cart) error {
	tx, err := cr.DB.BeginTx(context.Background(), nil)
	if err != nil {
		internalError.GenInternalError(err)
	}

	queryStr := "INSERT INTO cart (user_id, advert_id, amount) VALUES ($1, $2, $3);"
	_, err = tx.ExecContext(context.Background(), queryStr, cart.UserId, cart.AdvertId, cart.Amount)

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}

func (cr *CartRepository) Delete(cart *models.Cart) error {
	tx, err := cr.DB.BeginTx(context.Background(), nil)
	if err != nil {
		internalError.GenInternalError(err)
	}

	queryStr := "DELETE FROM cart WHERE user_id = $1 AND advert_id = $2;"
	_, err = tx.ExecContext(context.Background(), queryStr, cart.UserId, cart.AdvertId)

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	return nil
}

func (cr *CartRepository) DeleteAll(userId int64) error {
	tx, err := cr.DB.BeginTx(context.Background(), nil)
	if err != nil {
		internalError.GenInternalError(err)
	}

	queryStr := "DELETE FROM cart WHERE user_id = $1;"
	ct, err := tx.ExecContext(context.Background(), queryStr, userId)

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return internalError.RollbackError
		}
		return internalError.GenInternalError(err)
	}

	err = tx.Commit()
	if err != nil {
		return internalError.NotCommited
	}

	if ra, _ := ct.RowsAffected(); ra == 0 {
		return internalError.EmptyQuery
	}

	return nil
}
