package server

import (
	"context"
	"net"
	sessions "yula/internal/services/auth"
	proto "yula/proto/generated/auth"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServer struct {
	su     sessions.SessionUsecase
	logger *logrus.Logger
}

func NewAuthGRPCServer(logger *logrus.Logger, su sessions.SessionUsecase) *AuthServer {
	server := &AuthServer{
		su:     su,
		logger: logger,
	}
	return server
}

func (server *AuthServer) NewGRPCServer(listenUrl string) error {
	lis, err := net.Listen("tcp", listenUrl)
	server.logger.Infof("AUTH: my listen url %s \n", listenUrl)

	if err != nil {
		server.logger.Errorf("can not listen url: %s err :%v\n", listenUrl, err)
		return err
	}

	serv := grpc.NewServer()
	proto.RegisterAuthServer(serv, server)

	server.logger.Info("Start session service\n")
	return serv.Serve(lis)
}

func (s *AuthServer) Check(ctx context.Context, sessionID *proto.SessionID) (*proto.Result, error) {
	res, err := s.su.Check(sessionID.ID)
	if err != nil {
		s.logger.Errorf("can not check session with sessionID = %s, err = %v", sessionID.ID,
			err)
		return nil, err
	}

	return &proto.Result{
		UserID:    res.UserId,
		SessionID: res.Value,
		ExpireAt:  timestamppb.New(res.ExpiresAt),
	}, nil
}

func (s *AuthServer) Create(ctx context.Context, userID *proto.UserID) (*proto.Result, error) {
	res, err := s.su.Create(userID.ID)
	if err != nil {
		s.logger.Errorf("can not create session with userID = %d, err = %v", userID.ID,
			err)
		return nil, err
	}

	return &proto.Result{
		UserID:    res.UserId,
		SessionID: res.Value,
		ExpireAt:  timestamppb.New(res.ExpiresAt),
	}, nil
}

func (s *AuthServer) Delete(ctx context.Context, sessionID *proto.SessionID) (*proto.Nothing, error) {
	err := s.su.Delete(sessionID.ID)
	if err != nil {
		s.logger.Errorf("can not delete session with sessionID = %s, err = %v", sessionID.ID,
			err)
		return &proto.Nothing{Dummy: false}, err
	}

	return &proto.Nothing{
		Dummy: true,
	}, nil
}
