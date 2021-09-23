package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yula/internal/models"
	"yula/internal/pkg/user"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUse user.UserUsecase
	// сессии ?
}

func NewUserHandler(userUse user.UserUsecase) *UserHandler {
	return &UserHandler{
		userUse: userUse,
	}
}

func (uh *UserHandler) Configurate(r *mux.Router) {
	r.HandleFunc("/signup", uh.SignUpHandler).Methods(http.MethodPost)
}

func (uh *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signUpUser models.UserSignUp

	err := json.NewDecoder(r.Body).Decode(&signUpUser)
	if err != nil {
		status := models.StatusByCode(models.BadRequest)
		status.Message = err.Error()

		w.WriteHeader(status.HttpCode)

		jsonStatus := models.ToJson(status)
		w.Write(jsonStatus)
		return
	}

	fmt.Printf("User %s got\n", signUpUser.Email)

	user, status := uh.userUse.Create(&signUpUser)
	if status != models.StatusByCode(models.Created) {
		w.WriteHeader(status.HttpCode)
		jsonStatus := models.ToJson(status)
		w.Write(jsonStatus)
		return
	}

	fmt.Printf("User %s transformed\n", user.Email)

	w.WriteHeader(status.HttpCode)
	w.Write(models.ToJson(status))
}
