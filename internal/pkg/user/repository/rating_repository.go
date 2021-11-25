package repository

import (
	"context"
	"database/sql"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/user"
)

type RatingRepository struct {
	db *sql.DB
}

func NewRatingRepository(db *sql.DB) user.RatingRepository {
	return &RatingRepository{
		db: db,
	}
}

func (rr *RatingRepository) SelectRating(userFrom int64, userTo int64) (*models.Rating, error) {
	rating := &models.Rating{}
	query := rr.db.QueryRow("SELECT rate FROM rating WHERE user_from = $1 AND user_to = $2;", userFrom, userTo)

	err := query.Scan(&rating.Rating)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return nil, internalError.EmptyQuery

		default:
			return nil, internalError.GenInternalError(err)
		}
	}

	rating.UserFrom = userFrom
	rating.UserTo = userTo

	return rating, nil
}

func (rr *RatingRepository) InsertRating(rating *models.Rating) error {
	tx, err := rr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = rr.db.Exec("INSERT INTO rating(user_from, user_to, rate) VALUES ($1, $2, $3);",
		rating.UserFrom, rating.UserTo, rating.Rating)
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

	return nil
}

func (rr *RatingRepository) DeleteRating(rating *models.Rating) error {
	tx, err := rr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = rr.db.Exec("DELETE FROM rating WHERE user_from = $1 AND user_to = $2;",
		rating.UserFrom, rating.UserTo)
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

	return nil
}

func (rr *RatingRepository) UpdateRating(rating *models.Rating) error {
	tx, err := rr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = rr.db.Exec("UPDATE rating SET rate = $3 WHERE user_from = $1 AND user_to = $2;",
		rating.UserFrom, rating.UserTo, rating.Rating)
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

	return nil
}

func (rr *RatingRepository) SelectStat(userId int64) (int64, int64, error) {
	var sum, count int64

	query := rr.db.QueryRow("SELECT sum, count FROM rating_statistics WHERE user_id = $1;", userId)
	err := query.Scan(&sum, &count)
	if err != nil {
		return 0, 0, internalError.GenInternalError(err)
	}

	return sum, count, nil
}

func (rr *RatingRepository) InsertStat(userId int64) error {
	tx, err := rr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = rr.db.Exec("INSERT INTO rating_statistics(user_id) VALUES ($1);", userId)
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

	return nil
}

func (rr *RatingRepository) UpdateStat(userId int64, rate int, count int) error {
	tx, err := rr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = rr.db.Exec("UPDATE rating_statistics SET sum = sum + $2, count = count + $3 WHERE user_id = $1;",
		userId, rate, count)
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

	return nil
}
