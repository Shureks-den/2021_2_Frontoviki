package delivery

import (
	"encoding/json"
	"net/http"
	"yula/internal/codes"
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

func (uh *UserHandler) Routing(r *mux.Router) {
	r.HandleFunc("/signup", uh.SignUpHandler).Methods(http.MethodPost)
}

func (uh *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signUpUser models.UserSignUp

	err := json.NewDecoder(r.Body).Decode(&signUpUser)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")

		response := models.HttpError{Code: http.StatusBadRequest, Message: err.Error()}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	user, servErr := uh.userUse.Create(&signUpUser)
	if servErr != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")

		httpStat := codes.ServerErrorToHttpStatus(servErr)
		response := models.HttpError{Code: httpStat.Code, Message: httpStat.Message}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	response := models.HttpUser{Code: http.StatusCreated, Message: "user created successfully",
		Body: models.HttpBodyUser{User: user.RemovePassword()}}
	js, _ := json.Marshal(response)

	w.Write(js)
}
