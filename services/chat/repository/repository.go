package repository

import (
	"context"
	"database/sql"
	"regexp"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/services/chat"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) chat.ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (cr *ChatRepository) SelectMessages(idFrom int64, idTo int64, idAdv int64, offset int64, limit int64) ([]*models.Message, error) {
	query := `SELECT user_from, user_to, adv_id, msg, created_at FROM messages
			  WHERE user_from IN ($1, $2) AND user_to IN ($1, $2) AND adv_id = $3
			  ORDER BY created_at
			  OFFSET $4 LIMIT $5;`

	rows, err := cr.db.QueryContext(context.Background(), query, idFrom, idTo, idAdv, offset, limit)
	if err != nil {
		return nil, internalError.GenInternalError(err)
	}
	defer rows.Close()

	messages := make([]*models.Message, 0)
	for rows.Next() {
		message := &models.Message{}

		err := rows.Scan(&message.IdFrom, &message.IdTo, &message.IdAdv, &message.Msg, &message.CreatedAt)
		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (cr *ChatRepository) InsertMessage(message *models.Message) error {
	tx, err := cr.db.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = cr.db.Exec("INSERT INTO messages(user_from, user_to, adv_id, msg) VALUES ($1, $2, $3, $4);",
		message.IdFrom, message.IdTo, message.IdAdv, message.Msg)
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

func (cr *ChatRepository) DeleteMessages(idFrom int64, idTo int64, idAdv int64) error {
	tx, err := cr.db.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = cr.db.Exec("DELETE FROM messages WHERE user_from = $1 AND user_to = $2 AND adv_id = $3;",
		idFrom, idTo, idAdv)
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

func (cr *ChatRepository) SelectDialog(id1 int64, id2 int64, idAdv int64) (*models.Dialog, error) {
	query := `SELECT user1, user2, adv_id, created_at FROM dialogs
			  WHERE user1 = $1 AND user2 = $2 AND adv_id = $3
			  ORDER BY created_at;`

	row := cr.db.QueryRowContext(context.Background(), query, id1, id2, idAdv)

	dialog := &models.Dialog{}
	err := row.Scan(&dialog.Id1, &dialog.Id2, &dialog.IdAdv, &dialog.CreatedAt)
	if err != nil {
		res, _ := regexp.Match(".*no rows.*", []byte(err.Error()))
		if res {
			return nil, internalError.EmptyQuery
		} else {
			return nil, internalError.GenInternalError(err)
		}
	}

	return dialog, nil
}

func (cr *ChatRepository) InsertDialog(dialog *models.Dialog) error {
	tx, err := cr.db.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = cr.db.Exec("INSERT INTO dialogs(user1, user2, adv_id) VALUES ($1, $2, $3);",
		dialog.Id1, dialog.Id2, dialog.IdAdv)
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

func (cr *ChatRepository) DeleteDialog(dialog *models.Dialog) error {
	tx, err := cr.db.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = cr.db.Exec("DELETE FROM dialogs WHERE user1 = $1 AND user2 = $2 AND adv_id = $3;",
		dialog.Id1, dialog.Id2, dialog.IdAdv)

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

func (cr *ChatRepository) SelectAllDialogs(id1 int64) ([]*models.Dialog, error) {
	query := `SELECT user1, user2, adv_id, created_at FROM dialogs
			  WHERE user1 = $1
			  ORDER BY created_at DESC;`

	rows, err := cr.db.QueryContext(context.Background(), query, id1)
	if err != nil {
		return nil, internalError.GenInternalError(err)
	}
	defer rows.Close()

	dialogs := make([]*models.Dialog, 0)
	for rows.Next() {
		dialog := &models.Dialog{}

		err := rows.Scan(&dialog.Id1, &dialog.Id2, &dialog.IdAdv, &dialog.CreatedAt)
		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		dialogs = append(dialogs, dialog)
	}

	return dialogs, nil
}
