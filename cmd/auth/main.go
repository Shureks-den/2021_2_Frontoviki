package main

import (
	"yula/internal/config"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"

	"yula/internal/pkg/logging"

	sessRep "yula/internal/services/auth/repository"
	sessUse "yula/internal/services/auth/usecase"

	authServer "yula/internal/services/auth/server"
	// _ "yula/docs"
)

// @title Volchock's API
// @version 1.0
// @description Advert placement service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath /
func main() {
	logger := logging.GetLogger()

	if err := config.LoadConfig(); err != nil {
		logger.Errorf("error with load config: %s", err.Error())
		return
	}

	sr := sessRep.NewSessionRepository(config.Cfg.GetTarantoolCfg())

	su := sessUse.NewSessionUsecase(sr)

	grpcAuth := authServer.NewAuthGRPCServer(logrus.New(), su)
	grpcAuth.NewGRPCServer(config.Cfg.GetAuthEndPoint())

}
