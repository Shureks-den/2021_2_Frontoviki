package main

import (
	"fmt"
	"net/http"
	"time"
	"yula/internal/models"
	delivery "yula/internal/pkg/user/delivery/http"
	"yula/internal/pkg/user/repository"
	"yula/internal/pkg/user/usecase"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	var users []models.UserData
	users = append(users, models.UserData{uuid.New(), "edin", "edin@kasy.com", "1234", time.Now()})

	ur := repository.NewUserRepository(users)
	uu := usecase.NewUserUsecase(ur)
	uh := delivery.NewUserHandler(uu)

	uh.Configurate(r)

	fmt.Println("start serving :8080")

	http.ListenAndServe(":8080", r)
}
