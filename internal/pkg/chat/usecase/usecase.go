package usecase

import (
	"yula/internal/models"
	"yula/internal/pkg/chat"

	internalError "yula/internal/error"
)

type ChatUsecase struct {
	chatRepo chat.ChatRepository
}

func NewChatUsecase(repo chat.ChatRepository) chat.ChatUsecase {
	return &ChatUsecase{
		chatRepo: repo,
	}
}

func (cu *ChatUsecase) Create(message *models.Message) error {
	_, err := cu.chatRepo.SelectDialog(message.IdFrom, message.IdTo)
	if err == internalError.EmptyQuery {
		dialog := &models.Dialog{
			Id1: message.IdFrom,
			Id2: message.IdTo,
		}
		cu.chatRepo.InsertDialog(dialog)
	}
	_, err = cu.chatRepo.SelectDialog(message.IdTo, message.IdFrom)
	if err == internalError.EmptyQuery {
		dialog := &models.Dialog{
			Id1: message.IdTo,
			Id2: message.IdFrom,
		}
		cu.chatRepo.InsertDialog(dialog)
	}
	return cu.chatRepo.InsertMessage(message)
}

func (cu *ChatUsecase) Clear(idFrom int64, idTo int64) error {
	dialog, err := cu.chatRepo.SelectDialog(idFrom, idTo)
	if dialog != nil {
		dialog := &models.Dialog{
			Id1: idFrom,
			Id2: idTo,
		}
		cu.chatRepo.DeleteDialog(dialog)
	}

	if err != nil && err != internalError.EmptyQuery {
		return err
	}

	dialog, err = cu.chatRepo.SelectDialog(idTo, idFrom)
	if dialog == nil && err == internalError.EmptyQuery {
		err = cu.chatRepo.DeleteMessages(idFrom, idTo)
		if err != nil {
			return err
		}
		cu.chatRepo.DeleteMessages(idTo, idFrom)
		if err != nil {
			return err
		}
	}

	return err
}

func (cu *ChatUsecase) GetHistory(idFrom int64, idTo int64, offset int64, limit int64) ([]*models.Message, error) {
	_, err := cu.chatRepo.SelectDialog(idFrom, idTo)
	if err == internalError.EmptyQuery {
		return nil, internalError.NotExist
	}
	return cu.chatRepo.SelectMessages(idFrom, idTo, offset, limit)
}

func (cu *ChatUsecase) GetDialogs(idFrom int64) ([]*models.Dialog, error) {
	return cu.chatRepo.SelectAllDialogs(idFrom)
}
