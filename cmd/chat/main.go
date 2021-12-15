package main

import (
	"database/sql"
	"log"
	"yula/internal/config"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"

	"yula/internal/pkg/logging"

	chatRep "yula/internal/services/chat/repository"
	chatUse "yula/internal/services/chat/usecase"

	chatServer "yula/internal/services/chat/server"
	// _ "yula/docs"
)

func getPostgres(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
	}
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(10)
	return db
}

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

	sqlDB := getPostgres(config.Cfg.GetPostgresUrl())
	defer sqlDB.Close()

	chr := chatRep.NewChatRepository(sqlDB)

	chu := chatUse.NewChatUsecase(chr)

	grpcChat := chatServer.NewChatGRPCServer(logrus.New(), chu)
	err := grpcChat.NewGRPCServer(config.Cfg.GetChatEndPoint())
	if err != nil {
		logger.Errorf("error with load grpc: %s", err.Error())
		return
	}
}
