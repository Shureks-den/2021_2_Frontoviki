package chat

import "yula/internal/models"

type ChatUsecase interface {
	GetHistory(idFrom int64, idTo int64, idAdv int64, offset int64, limit int64) ([]*models.Message, error)
	Create(message *models.Message) error
	Clear(idFrom int64, idTo int64, idAdv int64) error

	GetDialogs(idFrom int64) ([]*models.Dialog, error)
}
