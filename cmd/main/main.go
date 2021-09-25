package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"yula/internal/config"
	"yula/internal/database"
	delivery "yula/internal/pkg/user/delivery/http"
	"yula/internal/pkg/user/repository"
	"yula/internal/pkg/user/usecase"

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

	DatabaseUrl, exists := os.LookupEnv("DATABASE_URL")
	if exists {
		fmt.Println(DatabaseUrl)
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	r := mux.NewRouter()

	ur := repository.NewUserRepository(postgres.GetDbPool())
	uu := usecase.NewUserUsecase(ur)
	uh := delivery.NewUserHandler(uu)

	uh.Routing(r)

	fmt.Println("start serving :8080")

	http.ListenAndServe(":8080", r)
}
