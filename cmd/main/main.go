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

	advtHttp "yula/internal/pkg/advt/delivery/http"
	advtRep "yula/internal/pkg/advt/repository"
	advtUse "yula/internal/pkg/advt/usecase"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
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
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.JsonMiddleware)

	ar := advtRep.NewAdvtRepository(postgres.GetDbPool())
	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	//&cnfg.TarantoolCfg
	au := advtUse.NewAdvtUsecase(ar)
	uu := userUse.NewUserUsecase(ur)
	su := sessUse.NewSessionUsecase(sr)

	ah := advtHttp.NewAdvtHandler(au)
	uh := userHttp.NewUserHandler(uu, su)
	sh := sessHttp.NewSessionHandler(su, uu)

	sm := middleware.NewSessionMiddleware(su)

	ah.Routing(r)
	uh.Routing(r, sm)
	sh.Routing(r)

	//http
	// fmt.Println("start serving ::80")
	// error := http.ListenAndServe(":80", r)

	//https
	fmt.Println("start serving ::443")
	error := http.ListenAndServeTLS(":443", "certificate.crt", "key.key", r)

	fmt.Println(error)
}
