package repository

import (
	"context"
	"database/sql"
	"regexp"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/services/chat"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) chat.ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (cr *ChatRepository) SelectMessages(iMessage *models.IMessage, offset int64, limit int64) ([]*models.Message, error) {
	query := `SELECT user_from, user_to, adv_id, msg, created_at FROM messages
			  WHERE user_from IN ($1, $2) AND user_to IN ($1, $2) AND adv_id = $3
			  ORDER BY created_at
			  OFFSET $4 LIMIT $5;`

	rows, err := cr.db.QueryContext(context.Background(), query, iMessage.IdFrom, iMessage.IdTo, iMessage.IdAdv, offset, limit)
	if err != nil {
		return nil, internalError.GenInternalError(err)
	}
	defer rows.Close()

	messages := make([]*models.Message, 0)
	for rows.Next() {
		message := &models.Message{}
		var adId sql.NullInt64
		err := rows.Scan(&message.MI.IdFrom, &message.MI.IdTo, &adId, &message.Msg, &message.CreatedAt)

		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		if !adId.Valid {
			adId.Int64 = -1
		}
		message.MI.IdAdv = adId.Int64

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
		message.MI.IdFrom, message.MI.IdTo, message.MI.IdAdv, message.Msg)
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

func (cr *ChatRepository) DeleteMessages(iMessage *models.IMessage) error {
	tx, err := cr.db.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = cr.db.Exec("DELETE FROM messages WHERE user_from = $1 AND user_to = $2 AND adv_id = $3;",
		iMessage.IdFrom, iMessage.IdTo, iMessage.IdAdv)
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

func (cr *ChatRepository) SelectDialog(iDialog *models.IDialog) (*models.Dialog, error) {
	query := `SELECT user1, user2, adv_id, created_at FROM dialogs
			  WHERE user1 = $1 AND user2 = $2 AND adv_id = $3
			  ORDER BY created_at;`

	row := cr.db.QueryRowContext(context.Background(), query, iDialog.Id1, iDialog.Id2, iDialog.IdAdv)

	dialog := &models.Dialog{}
	err := row.Scan(&dialog.DI.Id1, &dialog.DI.Id2, &dialog.DI.IdAdv, &dialog.CreatedAt)
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
		dialog.DI.Id1, dialog.DI.Id2, dialog.DI.IdAdv)
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

func (cr *ChatRepository) DeleteDialog(iDialog *models.IDialog) error {
	tx, err := cr.db.BeginTx(context.Background(), nil)
	if err != nil {
		return internalError.GenInternalError(err)
	}

	_, err = cr.db.Exec("DELETE FROM dialogs WHERE user1 = $1 AND user2 = $2 AND adv_id = $3;",
		iDialog.Id1, iDialog.Id2, iDialog.IdAdv)

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
		var adId sql.NullInt64
		err := rows.Scan(&dialog.DI.Id1, &dialog.DI.Id2, &adId, &dialog.CreatedAt)

		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		if !adId.Valid {
			adId.Int64 = -1
		}
		dialog.DI.IdAdv = adId.Int64

		dialogs = append(dialogs, dialog)
	}

	return dialogs, nil
}
