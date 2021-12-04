package chat

import (
	"yula/internal/models"
)

type ChatRepository interface {
	SelectMessages(iMessage *models.IMessage, offset int64, limit int64) ([]*models.Message, error)
	InsertMessage(message *models.Message) error
	DeleteMessages(iMessage *models.IMessage) error

	SelectDialog(iDialog *models.IDialog) (*models.Dialog, error)
	InsertDialog(dialog *models.Dialog) error
	DeleteDialog(dialog *models.IDialog) error

	SelectAllDialogs(id1 int64) ([]*models.Dialog, error)
}
