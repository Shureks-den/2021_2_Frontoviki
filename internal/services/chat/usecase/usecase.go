package usecase

import (
	"yula/internal/models"
	"yula/internal/services/chat"

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
	_, err := cu.chatRepo.SelectDialog(message.IdFrom, message.IdTo, message.IdAdv)
	if err == internalError.EmptyQuery {
		dialog := &models.Dialog{
			Id1:   message.IdFrom,
			Id2:   message.IdTo,
			IdAdv: message.IdAdv,
		}
		cu.chatRepo.InsertDialog(dialog)
	}
	_, err = cu.chatRepo.SelectDialog(message.IdTo, message.IdFrom, message.IdAdv)
	if err == internalError.EmptyQuery {
		dialog := &models.Dialog{
			Id1:   message.IdTo,
			Id2:   message.IdFrom,
			IdAdv: message.IdAdv,
		}
		cu.chatRepo.InsertDialog(dialog)
	}
	return cu.chatRepo.InsertMessage(message)
}

func (cu *ChatUsecase) Clear(idFrom int64, idTo int64, idAdv int64) error {
	dialog, err := cu.chatRepo.SelectDialog(idFrom, idTo, idAdv)
	if dialog != nil {
		dialog := &models.Dialog{
			Id1:   idFrom,
			Id2:   idTo,
			IdAdv: idAdv,
		}
		cu.chatRepo.DeleteDialog(dialog)
	}

	if err != nil && err != internalError.EmptyQuery {
		return err
	}

	dialog, err = cu.chatRepo.SelectDialog(idTo, idFrom, idAdv)
	if dialog == nil && err == internalError.EmptyQuery {
		err = cu.chatRepo.DeleteMessages(idFrom, idTo, idAdv)
		if err != nil {
			return err
		}
		cu.chatRepo.DeleteMessages(idTo, idFrom, idAdv)
		if err != nil {
			return err
		}
	}

	return err
}

func (cu *ChatUsecase) GetHistory(idFrom int64, idTo int64, idAdv int64, offset int64, limit int64) ([]*models.Message, error) {
	_, err := cu.chatRepo.SelectDialog(idFrom, idTo, idAdv)
	if err == internalError.EmptyQuery {
		return nil, internalError.NotExist
	}
	return cu.chatRepo.SelectMessages(idFrom, idTo, idAdv, offset, limit)
}

func (cu *ChatUsecase) GetDialogs(idFrom int64) ([]*models.Dialog, error) {
	return cu.chatRepo.SelectAllDialogs(idFrom)
}
