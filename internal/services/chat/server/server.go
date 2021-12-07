package delivery

import (
	"context"
	"crypto/tls"
	"net"
	"yula/internal/models"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"

	chat "yula/internal/services/chat"
	proto "yula/proto/generated/chat"
)

type ChatServer struct {
	cu     chat.ChatUsecase
	logger *logrus.Logger
}

func NewChatGRPCServer(logger *logrus.Logger, cu chat.ChatUsecase) *ChatServer {
	server := &ChatServer{
		cu:     cu,
		logger: logger,
	}
	return server
}

func (server *ChatServer) NewGRPCServer(listenUrl string, certFile string, keyFile string) error {
	lis, err := net.Listen("tcp", listenUrl)
	server.logger.Infof("CHAT: my listen url %s \n", listenUrl)

	if err != nil {
		server.logger.Errorf("can not listen url: %s err :%v\n", listenUrl, err)
		return err
	}

	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		server.logger.Error(err.Error())
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	serv := grpc.NewServer(grpc.Creds(credentials.NewTLS(config)))
	proto.RegisterChatServer(serv, server)

	server.logger.Info("Start chat service\n")
	return serv.Serve(lis)
}

func (s *ChatServer) GetHistory(ctx context.Context, arg *proto.GetHistoryArg) (*proto.Messages, error) {
	res, err := s.cu.GetHistory(&models.IDialog{
		Id1:   arg.DI.Id1,
		Id2:   arg.DI.Id2,
		IdAdv: arg.DI.IdAdv,
	}, arg.FP.Offset, arg.FP.Limit)

	if err != nil {
		s.logger.Errorf("can not get history from %d to %d on %d, err = %v", arg.DI.Id1, arg.DI.Id2, arg.DI.IdAdv, err)
		return nil, err
	}
	var messages *proto.Messages = &proto.Messages{}
	for _, message := range res {
		messages.M = append(messages.M, &proto.Message{
			MI: &proto.MessageIdentifier{
				IdFrom: message.MI.IdFrom,
				IdTo:   message.MI.IdTo,
				IdAdv:  message.MI.IdAdv,
			},
			Msg:       message.Msg,
			CreatedAt: timestamppb.New(message.CreatedAt),
		})
	}
	return messages, nil
}

func (s *ChatServer) Create(ctx context.Context, message *proto.Message) (*proto.Nothing, error) {
	err := s.cu.Create(&models.Message{
		MI: models.IMessage{
			IdFrom: message.MI.IdFrom,
			IdTo:   message.MI.IdTo,
			IdAdv:  message.MI.IdAdv,
		},
		Msg:       message.Msg,
		CreatedAt: message.CreatedAt.AsTime(),
	})
	if err != nil {
		s.logger.Errorf("can not create message from %d to %d on %d, err = %v", message.MI.IdFrom, message.MI.IdTo, message.MI.IdAdv, err)
		return &proto.Nothing{Dummy: false}, err
	}

	return &proto.Nothing{
		Dummy: true,
	}, nil
}

func (s *ChatServer) CreateDialog(ctx context.Context, dialog *proto.Dialog) (*proto.Nothing, error) {
	err := s.cu.CreateDialog(&models.Dialog{
		DI: models.IDialog{
			Id1:   dialog.DI.Id1,
			Id2:   dialog.DI.Id2,
			IdAdv: dialog.DI.IdAdv,
		},
		CreatedAt: dialog.CreatedAt.AsTime(),
	})
	if err != nil {
		s.logger.Errorf("can not create dialog from %d to %d on %d, err = %v", dialog.DI.Id1, dialog.DI.Id2, dialog.DI.IdAdv, err)
		return &proto.Nothing{Dummy: false}, err
	}

	return &proto.Nothing{
		Dummy: true,
	}, nil
}

func (s *ChatServer) Clear(ctx context.Context, DI *proto.DialogIdentifier) (*proto.Nothing, error) {
	err := s.cu.Clear(&models.IDialog{
		Id1:   DI.Id1,
		Id2:   DI.Id2,
		IdAdv: DI.IdAdv,
	})
	if err != nil {
		s.logger.Errorf("can not clear messages from %d to %d on %d, err = %v", DI.Id1, DI.Id2, DI.IdAdv, err)
		return &proto.Nothing{Dummy: false}, err
	}

	return &proto.Nothing{
		Dummy: true,
	}, nil
}

func (s *ChatServer) GetDialogs(ctx context.Context, UI *proto.UserIdentifier) (*proto.Dialogs, error) {
	res, err := s.cu.GetDialogs(UI.IdFrom)
	if err != nil {
		s.logger.Errorf("can not get dialogs from %d, err = %v", UI.IdFrom, err)
		return nil, err
	}
	var dialogs *proto.Dialogs = &proto.Dialogs{}
	for _, dialog := range res {
		dialogs.D = append(dialogs.D, &proto.Dialog{
			DI: &proto.DialogIdentifier{
				Id1:   dialog.DI.Id1,
				Id2:   dialog.DI.Id2,
				IdAdv: dialog.DI.IdAdv,
			},
			CreatedAt: timestamppb.New(dialog.CreatedAt),
		})
	}
	return dialogs, nil
}
