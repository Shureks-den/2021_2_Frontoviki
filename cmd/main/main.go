package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"yula/internal/config"
	"yula/internal/database"
	imageloaderRepo "yula/internal/pkg/image_loader/repository"
	imageloaderUse "yula/internal/pkg/image_loader/usecase"
	"yula/internal/pkg/logging"
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

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "yula/docs" // docs is generated by Swag CLI, you have to import it.

	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {
	// loads values from .env into the system
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-2], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		log.Fatal("No .env file found")
	}

	govalidator.SetFieldsRequiredByDefault(true)
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

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		logger.Fatalf("db error instance", err.Error())
		return
	}
	defer postgres.Close()

	r := mux.NewRouter()
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	api := mux.NewRouter()

	api.Use(middleware.CorsMiddleware)
	api.Use(middleware.ContentTypeMiddleware)
	api.Use(middleware.LoggerMiddleware)

	ar := advtRep.NewAdvtRepository(postgres.GetDbPool())
	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	ilr := imageloaderRepo.NewImageLoaderRepository()

	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)
	au := advtUse.NewAdvtUsecase(ar, ilu)
	uu := userUse.NewUserUsecase(ur, ilu)
	su := sessUse.NewSessionUsecase(sr)

	ah := advtHttp.NewAdvertHandler(au, uu, logger)
	uh := userHttp.NewUserHandler(uu, su, logger)
	sh := sessHttp.NewSessionHandler(su, uu, logger)

	sm := middleware.NewSessionMiddleware(su)

	ah.Routing(api, sm)
	uh.Routing(api, sm)
	sh.Routing(api)

	//http
	fmt.Println("start serving ::8080")
	error := http.ListenAndServe(":8080", r)

	// //https
	// fmt.Println("start serving ::5000")
	// error := http.ListenAndServeTLS(":5000", "certificate.crt", "key.key", r)

	logger.Errorf("http serve error %v", error)
}
