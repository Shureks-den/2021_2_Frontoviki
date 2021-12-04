package chat

import "yula/internal/models"

type ChatUsecase interface {
	GetHistory(iDialog *models.IDialog, offset int64, limit int64) ([]*models.Message, error)
	Create(message *models.Message) error
	CreateDialog(dialog *models.Dialog) error
	Clear(iDialog *models.IDialog) error

	GetDialogs(idFrom int64) ([]*models.Dialog, error)
}
