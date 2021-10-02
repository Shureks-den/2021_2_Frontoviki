package delivery

import (
	"encoding/json"
	"net/http"
	"time"
	internalError "yula/internal/error"
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
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	user, err := sh.userUsecase.GetByEmail(signInUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	err = sh.userUsecase.CheckPassword(user, signInUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	userSession, err := sh.sessionUsecase.Create(user.Id)
	if err != nil {
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    userSession.Value,
		Expires:  userSession.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "signin successfully", nil))
}

func (sh *SessionHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.Unauthorized)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	err = sh.sessionUsecase.Delete(session.Value)
	if err != nil {
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	session.Expires = time.Now().Add(-time.Minute)
	http.SetCookie(w, session)

	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "logout successfully", nil))
}
