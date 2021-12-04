package usecase

import (
	"time"
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

func (cu *ChatUsecase) CreateDialog(dialog *models.Dialog) error {
	_, err := cu.chatRepo.SelectDialog(dialog.ToIDialog())
	if err == internalError.EmptyQuery {
		err = cu.chatRepo.InsertDialog(dialog)
		if err != nil {
			return err
		}
	}
	return err
}

func (cu *ChatUsecase) Create(message *models.Message) error {
	dialog12 := &models.IDialog{
		Id1:   message.MI.IdFrom,
		Id2:   message.MI.IdTo,
		IdAdv: message.MI.IdAdv,
	}

	dialog21 := &models.IDialog{
		Id1:   message.MI.IdTo,
		Id2:   message.MI.IdFrom,
		IdAdv: message.MI.IdAdv,
	}

	_, err := cu.chatRepo.SelectDialog(dialog12)
	if err == internalError.EmptyQuery {
		err = cu.chatRepo.InsertDialog(dialog12.ToDialog(time.Now()))
		if err != nil {
			return err
		}
	}

	_, err = cu.chatRepo.SelectDialog(dialog21)
	if err == internalError.EmptyQuery {
		err = cu.chatRepo.InsertDialog(dialog21.ToDialog(time.Now()))
		if err != nil {
			return err
		}
	}

	return cu.chatRepo.InsertMessage(message)
}

func (cu *ChatUsecase) Clear(iDialog *models.IDialog) error {
	dialog, err := cu.chatRepo.SelectDialog(iDialog)
	if dialog != nil {
		cu.chatRepo.DeleteDialog(dialog.ToIDialog())
	}

	if err != nil && err != internalError.EmptyQuery {
		return err
	}

	dialog, err = cu.chatRepo.SelectDialog(&models.IDialog{
		Id1:   iDialog.Id2,
		Id2:   iDialog.Id1,
		IdAdv: iDialog.IdAdv,
	})

	if dialog == nil && err == internalError.EmptyQuery {
		err = cu.chatRepo.DeleteMessages(&models.IMessage{
			IdFrom: iDialog.Id1,
			IdTo:   iDialog.Id2,
			IdAdv:  iDialog.IdAdv,
		})
		if err != nil {
			return err
		}
		err = cu.chatRepo.DeleteMessages(&models.IMessage{
			IdFrom: iDialog.Id2,
			IdTo:   iDialog.Id1,
			IdAdv:  iDialog.IdAdv,
		})
		if err != nil {
			return err
		}
	}

	return err
}

func (cu *ChatUsecase) GetHistory(iDialog *models.IDialog, offset int64, limit int64) ([]*models.Message, error) {
	_, err := cu.chatRepo.SelectDialog(iDialog)
	if err == internalError.EmptyQuery {
		return nil, internalError.NotExist
	}
	return cu.chatRepo.SelectMessages(&models.IMessage{
		IdFrom: iDialog.Id1,
		IdTo:   iDialog.Id2,
		IdAdv:  iDialog.IdAdv,
	}, offset, limit)
}

func (cu *ChatUsecase) GetDialogs(idFrom int64) ([]*models.Dialog, error) {
	return cu.chatRepo.SelectAllDialogs(idFrom)
}
