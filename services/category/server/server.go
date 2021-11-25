package http

import (
	"context"
	"net"
	"yula/services/category"

	proto "yula/proto/generated/category"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type CategoryServer struct {
	cu     category.CategoryUsecase
	logger *logrus.Logger
}

func NewCategoryGRPCServer(logger *logrus.Logger, cu category.CategoryUsecase) *CategoryServer {
	server := &CategoryServer{
		cu:     cu,
		logger: logger,
	}
	return server
}

func (server *CategoryServer) NewGRPCServer(listenUrl string) error {
	lis, err := net.Listen("tcp", listenUrl)
	server.logger.Infof("Category: my listen url %s \n", listenUrl)

	if err != nil {
		server.logger.Errorf("can not listen url: %s err :%v\n", listenUrl, err)
		return err
	}

	serv := grpc.NewServer()
	proto.RegisterCategoryServer(serv, server)

	server.logger.Info("Start category service\n")
	return serv.Serve(lis)
}

func (s *CategoryServer) GetCategories(ctx context.Context, dummy *proto.Nothing) (*proto.Categories, error) {
	res, err := s.cu.GetCategories()
	if err != nil {
		s.logger.Errorf("can not get categories")
		return nil, err
	}
	var categories *proto.Categories = &proto.Categories{}
	for _, category := range res {
		categories.Categories = append(categories.Categories, &proto.XCategory{
			Name: category.Name,
		})
	}
	return categories, nil
}
