package main

import (
	"fmt"
	"log"
	"net/http"
	"yula/internal/config"
	"yula/internal/database"
	userHttp "yula/internal/pkg/user/delivery/http"
	userRep "yula/internal/pkg/user/repository"
	userUse "yula/internal/pkg/user/usecase"

	"yula/internal/pkg/middleware"
	sessHttp "yula/internal/pkg/session/delivery/http"
	sessRep "yula/internal/pkg/session/repository"
	sessUse "yula/internal/pkg/session/usecase"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer postgres.Close()

	r := mux.NewRouter()

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	//&cnfg.TarantoolCfg

	uu := userUse.NewUserUsecase(ur)
	su := sessUse.NewSessionUsecase(sr)

	uh := userHttp.NewUserHandler(uu, su)
	sh := sessHttp.NewSessionHandler(su, uu)

	sm := middleware.NewSessionMiddleware(su)

	s := r.PathPrefix("/users").Subrouter()
	uh.Routing(s, sm)
	sh.Routing(r)

	fmt.Println("start serving :8080")

	http.ListenAndServe(":8080", r)
}
