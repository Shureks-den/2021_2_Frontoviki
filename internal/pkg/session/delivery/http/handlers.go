package delivery

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
	internalError "yula/internal/error"
	auth "yula/proto/generated/auth"

	"yula/internal/models"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/user"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
)

type SessionHandler struct {
	sessionUsecase auth.AuthClient
	userUsecase    user.UserUsecase
}

func NewSessionHandler(sessionUsecase auth.AuthClient, userUsecase user.UserUsecase) *SessionHandler {
	return &SessionHandler{
		sessionUsecase: sessionUsecase, userUsecase: userUsecase,
	}
}

func (sh *SessionHandler) Routing(r *mux.Router) {
	r.HandleFunc("/signin", sh.SignInHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/logout", sh.LogOutHandler).Methods(http.MethodPost, http.MethodOptions)
}

var (
	logger logging.Logger = logging.GetLogger()
)

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
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var signInUser models.UserSignIn

	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warnf("cannot convert body to bytes: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	err = easyjson.Unmarshal(buf, &signInUser)
	if err != nil {
		logger.Warnf("cannot unmarshal: %s", err.Error())
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
		logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		return
	}

	user, err := sh.userUsecase.GetByEmail(signInUser.Email)
	if err != nil {
		logger.Warnf("can not get by email: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	err = sh.userUsecase.CheckPassword(user, signInUser.Password)
	if err != nil {
		logger.Warnf("wrong password check: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	userSession, err := sh.sessionUsecase.Create(context.Background(), &auth.UserID{ID: user.Id})
	if err != nil {
		logger.Warnf("can not create user: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    userSession.SessionID,
		Expires:  userSession.ExpireAt.AsTime(),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Secure:   true,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "signin successfully", nil))
	logger.Debug("signin successfully")
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
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	session, err := r.Cookie("session_id")
	if err != nil {
		logger.Warnf("unauthorized: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.Unauthorized)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	_, err = sh.sessionUsecase.Delete(context.Background(), &auth.SessionID{ID: session.Value})
	if err != nil {
		logger.Warnf("can not delete session: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	session.Expires = time.Now().Add(-time.Minute)
	http.SetCookie(w, session)

	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "logout successfully", nil))
	logger.Debug("logout successfully")
}
