package delivery

import (
	"context"
	"testing"
	"time"

	mocks "yula/internal/services/chat/mocks"
	proto "yula/proto/generated/chat"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetHistorySuccess(t *testing.T) {
	cu := mocks.ChatUsecase{}
	su := NewChatGRPCServer(logrus.New(), &cu)

	cu.On("GetHistory", mock.Anything, int64(1), int64(10)).Return(nil, nil)

	_, err := su.GetHistory(context.Background(), &proto.GetHistoryArg{
		DI: &proto.DialogIdentifier{
			Id1:   1,
			Id2:   2,
			IdAdv: 3,
		},
		FP: &proto.FilterParams{
			Offset: 1,
			Limit:  10,
		},
	})
	assert.Nil(t, err)
}

func TestCreateSuccess(t *testing.T) {
	cu := mocks.ChatUsecase{}
	su := NewChatGRPCServer(logrus.New(), &cu)

	cu.On("Create", mock.AnythingOfType("*models.Message")).Return(nil)

	_, err := su.Create(context.Background(), &proto.Message{
		MI: &proto.MessageIdentifier{
			IdFrom: 1,
			IdTo:   2,
			IdAdv:  3,
		},
		Msg:       "aboba",
		CreatedAt: timestamppb.New(time.Now()),
	})
	assert.Nil(t, err)
}

func TestCreateDialogSuccess(t *testing.T) {
	cu := mocks.ChatUsecase{}
	su := NewChatGRPCServer(logrus.New(), &cu)

	cu.On("CreateDialog", mock.AnythingOfType("*models.Dialog")).Return(nil)

	_, err := su.CreateDialog(context.Background(), &proto.Dialog{
		DI: &proto.DialogIdentifier{
			Id1:   1,
			Id2:   2,
			IdAdv: 3,
		},
		CreatedAt: timestamppb.New(time.Now()),
	})

	assert.Nil(t, err)
}

func TestClearSuccess(t *testing.T) {
	cu := mocks.ChatUsecase{}
	su := NewChatGRPCServer(logrus.New(), &cu)

	cu.On("Clear", mock.AnythingOfType("*models.IDialog")).Return(nil)

	_, err := su.Clear(context.Background(), &proto.DialogIdentifier{
		Id1:   1,
		Id2:   2,
		IdAdv: 3,
	})

	assert.Nil(t, err)
}

func TestGetDialogsSuccess(t *testing.T) {
	cu := mocks.ChatUsecase{}
	su := NewChatGRPCServer(logrus.New(), &cu)

	cu.On("GetDialogs", int64(10)).Return(nil, nil)

	_, err := su.GetDialogs(context.Background(), &proto.UserIdentifier{
		IdFrom: 10,
	})

	assert.Nil(t, err)
}
