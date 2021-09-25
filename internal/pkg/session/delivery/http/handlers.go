package delivery

import (
	"encoding/json"
	"net/http"
	"yula/internal/codes"
	"yula/internal/models"
	"yula/internal/pkg/session"
	"yula/internal/pkg/user"

	"github.com/gorilla/mux"
)

type SessionHandler struct {
	sessionUsecase session.SessionUsecase
	userUsecase    user.UserUsecase
}

func NewSessionHandler(sessionUsecase session.SessionUsecase, userUsecase user.UserUsecase) *SessionHandler {
	return &SessionHandler{
		sessionUsecase: sessionUsecase, userUsecase: userUsecase,
	}
}

func (sh *SessionHandler) Routing(r *mux.Router) {
	r.HandleFunc("/signin", sh.SignInHandler).Methods(http.MethodPost)
}

func (sh *SessionHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var signInUser models.UserSignIn

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&signInUser)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")

		response := models.HttpError{Code: http.StatusBadRequest, Message: err.Error()}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	user, serverErr := sh.userUsecase.GetByEmail(signInUser.Email)
	if serverErr != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")

		httpStat := codes.ServerErrorToHttpStatus(serverErr)
		response := models.HttpError{Code: httpStat.Code, Message: httpStat.Message}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	serverErr = sh.userUsecase.CheckPassword(user, signInUser.Password)
	if serverErr != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")

		httpStat := codes.ServerErrorToHttpStatus(serverErr)
		response := models.HttpError{Code: httpStat.Code, Message: httpStat.Message}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	userSession, err := sh.sessionUsecase.Create(user.Id)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")

		response := models.HttpError{Code: http.StatusInternalServerError, Message: "something went wrong"}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    userSession.Value,
		Expires:  userSession.ExpiresAt,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	response := models.HttpError{Code: http.StatusOK, Message: "successful signin"}
	js, _ := json.Marshal(response)

	w.Write(js)
}
