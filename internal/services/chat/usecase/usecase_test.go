package usecase

import (
	"testing"
	"time"
	"yula/internal/models"

	myerror "yula/internal/error"
	mocks "yula/internal/services/chat/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateDialogSuccess(t *testing.T) {
	cr := mocks.ChatRepository{}
	cu := ChatUsecase{chatRepo: &cr}

	dialog := models.Dialog{
		DI: models.IDialog{
			Id1:   1,
			Id2:   2,
			IdAdv: 1,
		},
		CreatedAt: time.Now(),
	}

	cr.On("SelectDialog", dialog.ToIDialog()).Return(nil, nil)

	error := cu.CreateDialog(&dialog)
	assert.Nil(t, error)
}

func TestCreateDialogFail(t *testing.T) {
	cr := mocks.ChatRepository{}
	cu := ChatUsecase{chatRepo: &cr}

	dialog := models.Dialog{
		DI: models.IDialog{
			Id1:   1,
			Id2:   2,
			IdAdv: 1,
		},
		CreatedAt: time.Now(),
	}

	cr.On("SelectDialog", dialog.ToIDialog()).Return(nil, myerror.EmptyQuery)
	cr.On("InsertDialog", &dialog).Return(myerror.AlreadyExist)

	error := cu.CreateDialog(&dialog)
	assert.Equal(t, error, myerror.AlreadyExist)
}

func TestCreateSuccess(t *testing.T) {
	cr := mocks.ChatRepository{}
	cu := ChatUsecase{chatRepo: &cr}

	message := models.Message{
		MI: models.IMessage{
			IdFrom: 1,
			IdTo:   2,
			IdAdv:  1,
		},
		CreatedAt: time.Now(),
	}

	cr.On("SelectDialog", mock.AnythingOfType("*models.IDialog")).Return(nil, nil)
	cr.On("InsertMessage", &message).Return(nil)

	error := cu.Create(&message)
	assert.Nil(t, error)
}

func TestClearSuccess(t *testing.T) {
	cr := mocks.ChatRepository{}
	cu := ChatUsecase{chatRepo: &cr}

	DI := models.IDialog{
		Id1:   1,
		Id2:   2,
		IdAdv: 1,
	}

	cr.On("SelectDialog", mock.AnythingOfType("*models.IDialog")).Return(nil, nil)
	cr.On("DeleteMessages", mock.AnythingOfType("*models.IMessage")).Return(nil)

	error := cu.Clear(&DI)
	assert.Nil(t, error)
}

func TestGetHistorySuccess(t *testing.T) {
	cr := mocks.ChatRepository{}
	cu := ChatUsecase{chatRepo: &cr}

	DI := models.IDialog{
		Id1:   1,
		Id2:   2,
		IdAdv: 1,
	}

	cr.On("SelectDialog", mock.AnythingOfType("*models.IDialog")).Return(nil, nil)
	cr.On("SelectMessages", mock.AnythingOfType("*models.IMessage"), int64(1), int64(2)).Return(nil, nil)

	_, error := cu.GetHistory(&DI, int64(1), int64(2))
	assert.Nil(t, error)
}
