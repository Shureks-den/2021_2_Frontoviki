package chat

import "yula/internal/models"

type ChatRepository interface {
	SelectMessages(idFrom int64, idTo int64, offset int64, limit int64) ([]*models.Message, error)
	InsertMessage(message *models.Message) error
	DeleteMessages(idFrom int64, idTo int64) error

	SelectDialog(idFrom int64, idTo int64) (*models.Dialog, error)
	InsertDialog(dialog *models.Dialog) error
	DeleteDialog(dialog *models.Dialog) error

	SelectAllDialogs(id1 int64) ([]*models.Dialog, error)
}
