package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
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
	r.HandleFunc("/signin", sh.SignInHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/logout", sh.LogOutHandler).Methods(http.MethodPost, http.MethodOptions)
}

func (sh *SessionHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var signInUser models.UserSignIn

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&signInUser)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := models.HttpError{Code: http.StatusBadRequest, Message: err.Error()}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	user, serverErr := sh.userUsecase.GetByEmail(signInUser.Email)
	if serverErr != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		httpStat := codes.ServerErrorToHttpStatus(serverErr)
		response := models.HttpError{Code: httpStat.Code, Message: httpStat.Message}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	serverErr = sh.userUsecase.CheckPassword(user, signInUser.Password)
	if serverErr != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		httpStat := codes.ServerErrorToHttpStatus(serverErr)
		response := models.HttpError{Code: httpStat.Code, Message: httpStat.Message}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	userSession, err := sh.sessionUsecase.Create(user.Id)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.HttpError{Code: http.StatusOK, Message: "signin successfully"}
	js, _ := json.Marshal(response)

	w.Write(js)
}

func (sh *SessionHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		httpStat := codes.ServerErrorToHttpStatus(codes.NewServerError(codes.Unauthorized))
		response := models.HttpError{Code: httpStat.Code, Message: httpStat.Message}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	err = sh.sessionUsecase.Delete(session.Value)
	if err != nil {
		log.Printf("Logout 2 : %s\n", err.Error())
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		httpStat := codes.ServerErrorToHttpStatus(codes.NewServerError(codes.InternalError))
		response := models.HttpError{Code: httpStat.Code, Message: httpStat.Message}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	session.Expires = time.Now().Add(-time.Minute)
	http.SetCookie(w, session)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.HttpError{Code: http.StatusOK, Message: "logout successfully"}
	js, _ := json.Marshal(response)

	w.Write(js)
}
