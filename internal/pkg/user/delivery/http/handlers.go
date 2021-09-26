package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"yula/internal/codes"
	"yula/internal/models"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/session"
	"yula/internal/pkg/user"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUsecase    user.UserUsecase
	sessionUsecase session.SessionUsecase
}

func NewUserHandler(userUsecase user.UserUsecase, sessionUsecase session.SessionUsecase) *UserHandler {
	return &UserHandler{
		userUsecase:    userUsecase,
		sessionUsecase: sessionUsecase,
	}
}

func (uh *UserHandler) Routing(r *mux.Router, sm *middleware.SessionMiddleware) {
	r.HandleFunc("/signup", uh.SignUpHandler).Methods(http.MethodPost)
	r.Handle("/profile", sm.CheckAuthorized(http.HandlerFunc(uh.GetProfileHandler))).Methods(http.MethodGet)
}

func (uh *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signUpUser models.UserSignUp

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&signUpUser)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := models.HttpError{Code: http.StatusBadRequest, Message: err.Error()}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	user, servErr := uh.userUsecase.Create(&signUpUser)
	if servErr != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		httpStat := codes.ServerErrorToHttpStatus(servErr)
		response := models.HttpError{Code: httpStat.Code, Message: httpStat.Message}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	userSession, err := uh.sessionUsecase.Create(user.Id)
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
	w.Header().Add("Location", r.Host+"/signin") // указываем в качестве перенаправления страницу входа
	w.WriteHeader(http.StatusCreated)

	response := models.HttpUser{Code: http.StatusCreated, Message: "user created successfully",
		Body: models.HttpBodyUser{User: user.RemovePassword()}}
	js, _ := json.Marshal(response)

	w.Write(js)
}

func (uh *UserHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.ContextUserId).(int64)

	log.Printf("User %d opened profile", userId)

	profile, serverErr := uh.userUsecase.GetById(userId)
	if serverErr != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		httpStat := codes.ServerErrorToHttpStatus(serverErr)
		response := models.HttpError{Code: httpStat.Code, Message: httpStat.Message}
		js, _ := json.Marshal(response)

		w.Write(js)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.HttpBodyInterface{Code: http.StatusOK, Message: "profile opened",
		Body: models.HttpBodyProfile{Profile: *profile}}
	js, _ := json.Marshal(response)

	w.Write(js)
}
