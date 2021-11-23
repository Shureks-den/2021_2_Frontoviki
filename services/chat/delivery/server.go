package delivery

import (
	"context"
	"net"
	"yula/internal/models"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "yula/proto/generated/chat"
	chat "yula/services/chat"
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

func (server *ChatServer) NewGRPCServer(listenUrl string) error {
	lis, err := net.Listen("tcp", listenUrl)
	server.logger.Infof("CHAT: my listen url %s \n", listenUrl)

	if err != nil {
		server.logger.Errorf("can not listen url: %s err :%v\n", listenUrl, err)
		return err
	}

	serv := grpc.NewServer()
	proto.RegisterChatServer(serv, server)

	server.logger.Info("Start chat service\n")
	return serv.Serve(lis)
}

func (s *ChatServer) GetHistory(ctx context.Context, arg *proto.GetHistoryArg) (*proto.Messages, error) {
	res, err := s.cu.GetHistory(arg.DI.IdFrom, arg.DI.IdTo, arg.DI.IdAdv, arg.FP.Offset, arg.FP.Limit)
	if err != nil {
		s.logger.Errorf("can not get history from %d to %d on %d, err = %v", arg.DI.IdFrom, arg.DI.IdTo, arg.DI.IdAdv, err)
		return nil, err
	}
	var messages *proto.Messages = &proto.Messages{}
	for _, message := range res {
		messages.M = append(messages.M, &proto.Message{
			IdFrom:    message.IdFrom,
			IdTo:      message.IdTo,
			IdAdv:     message.IdAdv,
			Msg:       message.Msg,
			CreatedAt: timestamppb.New(message.CreatedAt),
		})
	}
	return messages, nil
}

func (s *ChatServer) Create(ctx context.Context, message *proto.Message) (*proto.Nothing, error) {
	err := s.cu.Create(&models.Message{
		IdFrom:    message.IdFrom,
		IdTo:      message.IdTo,
		IdAdv:     message.IdAdv,
		Msg:       message.Msg,
		CreatedAt: message.CreatedAt.AsTime(),
	})
	if err != nil {
		s.logger.Errorf("can not create message from %d to %d on %d, err = %v", message.IdFrom, message.IdTo, message.IdAdv, err)
		return &proto.Nothing{Dummy: false}, err
	}

	return &proto.Nothing{
		Dummy: true,
	}, nil
}

func (s *ChatServer) Clear(ctx context.Context, DI *proto.DialogIdentifier) (*proto.Nothing, error) {
	err := s.cu.Clear(DI.IdFrom, DI.IdTo, DI.IdAdv)
	if err != nil {
		s.logger.Errorf("can not clear messages from %d to %d on %d, err = %v", DI.IdFrom, DI.IdTo, DI.IdAdv, err)
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
			Id1:       dialog.Id1,
			Id2:       dialog.Id2,
			IdAdv:     dialog.IdAdv,
			CreatedAt: timestamppb.New(dialog.CreatedAt),
		})
	}
	return dialogs, nil
}
