package delivery

import (
	"encoding/json"
	"net/http"
	"time"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/session"
	"yula/internal/pkg/user"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
)

type SessionHandler struct {
	sessionUsecase session.SessionUsecase
	userUsecase    user.UserUsecase
	logger         logging.Logger
}

func NewSessionHandler(sessionUsecase session.SessionUsecase, userUsecase user.UserUsecase, logger logging.Logger) *SessionHandler {
	return &SessionHandler{
		sessionUsecase: sessionUsecase, userUsecase: userUsecase, logger: logger,
	}
}

func (sh *SessionHandler) Routing(r *mux.Router) {
	r.HandleFunc("/signin", sh.SignInHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/logout", sh.LogOutHandler).Methods(http.MethodPost, http.MethodOptions)
}

// SignInHandler godoc
// @Summary Sign in
// @Description Sign in
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param user body models.UserSignIn true "User sign in data"
// @Success 200 {object} models.HttpBodyInterface
// @failure default {object} models.HttpError
// @Router /signin [post]
func (sh *SessionHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	sh.logger = sh.logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var signInUser models.UserSignIn

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&signInUser)
	if err != nil {
		sh.logger.Warnf("bad request: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	sanitizer := bluemonday.UGCPolicy()
	signInUser.Email = sanitizer.Sanitize(signInUser.Email)
	signInUser.Password = sanitizer.Sanitize(signInUser.Password)

	_, err = govalidator.ValidateStruct(signInUser)
	if err != nil {
		sh.logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		return
	}

	user, err := sh.userUsecase.GetByEmail(signInUser.Email)
	if err != nil {
		sh.logger.Warnf("can not get by email: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	err = sh.userUsecase.CheckPassword(user, signInUser.Password)
	if err != nil {
		sh.logger.Warnf("wrong password check: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	userSession, err := sh.sessionUsecase.Create(user.Id)
	if err != nil {
		sh.logger.Warnf("can not create user: %s", err.Error())
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
	sh.logger.Debug("signin successfully")
}

// SignInHandler godoc
// @Summary Log out
// @Description Log out
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} models.HttpBodyInterface
// @failure default {object} models.HttpError
// @Router /logout [post]
func (sh *SessionHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	sh.logger = sh.logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	session, err := r.Cookie("session_id")
	if err != nil {
		sh.logger.Warnf("unauthorized: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.Unauthorized)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	err = sh.sessionUsecase.Delete(session.Value)
	if err != nil {
		sh.logger.Warnf("can not delete session: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	session.Expires = time.Now().Add(-time.Minute)
	http.SetCookie(w, session)

	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "logout successfully", nil))
	sh.logger.Debug("logout successfully")
}
