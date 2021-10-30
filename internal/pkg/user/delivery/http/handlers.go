package delivery

import (
	"bytes"
	"encoding/json"
	"net/http"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/session"
	"yula/internal/pkg/user"

	"github.com/sirupsen/logrus"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUsecase    user.UserUsecase
	sessionUsecase session.SessionUsecase
	logger         logging.Logger
}

func NewUserHandler(userUsecase user.UserUsecase, sessionUsecase session.SessionUsecase, logger logging.Logger) *UserHandler {
	return &UserHandler{
		userUsecase:    userUsecase,
		sessionUsecase: sessionUsecase,
		logger:         logger,
	}
}

func (uh *UserHandler) Routing(r *mux.Router, sm *middleware.SessionMiddleware) {
	r.HandleFunc("/signup", uh.SignUpHandler).Methods(http.MethodPost, http.MethodOptions)

	s := r.PathPrefix("/users").Subrouter()
	s.Handle("/profile", sm.CheckAuthorized(http.HandlerFunc(uh.GetProfileHandler))).Methods(http.MethodGet, http.MethodOptions)
	s.Handle("/profile", sm.CheckAuthorized(http.HandlerFunc(uh.UpdateProfileHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/profile/upload", sm.CheckAuthorized(http.HandlerFunc(uh.UploadProfileImageHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/profile/password", sm.CheckAuthorized(http.HandlerFunc(uh.ChangePasswordHandler))).Methods(http.MethodPost, http.MethodOptions)
}

// SignUpHandler godoc
// @Summary Sign up
// @Description Sign up
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param user body models.UserSignUp true "User sign up data"
// @Success 200 {object} models.HttpBodyInterface{body=models.Profile}
// @failure default {object} models.HttpError
// @Router /signup [post]
func (uh *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signUpUser models.UserSignUp
	uh.logger = uh.logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&signUpUser)
	if err != nil {
		uh.logger.Warnf("bad request: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	_, err = govalidator.ValidateStruct(signUpUser)
	if err != nil {
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(signUpUser)
		uh.logger.Warnf("invalid data: %s", buf.String())

		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		return
	}

	user, servErr := uh.userUsecase.Create(&signUpUser)
	if servErr != nil {
		uh.logger.Warnf("can not create user: %s", servErr.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(servErr)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	userSession, err := uh.sessionUsecase.Create(user.Id)
	if err != nil {
		uh.logger.Warnf("can not create session based on user %d: %s", user.Id, err.Error())
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

	w.Header().Add("Location", r.Host+"/signin") // указываем в качестве перенаправления страницу входа
	w.WriteHeader(http.StatusOK)

	body := models.HttpBodyProfile{Profile: *user.ToProfile()}
	w.Write(models.ToBytes(http.StatusCreated, "user created successfully", body))
	uh.logger.Debugf("user %d created successfully", user.Id)
}

// GetProfileHandler godoc
// @Summary Get user's profile
// @Description Get user's profile
// @Tags user
// @Accept application/json
// @Produce application/json
// @Success 200 {object} models.HttpBodyInterface{body=models.Profile}
// @failure default {object} models.HttpError
// @Router /users/profile [get]
func (uh *UserHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	var userId int64
	uh.logger = uh.logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	profile, err := uh.userUsecase.GetById(userId)
	if err != nil {
		uh.logger.Warnf("can not get user with id %d: %s", userId, err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)

	body := models.HttpBodyProfile{Profile: *profile}
	w.Write(models.ToBytes(http.StatusOK, "profile provided", body))
	uh.logger.Debugf("user %d created successfully", userId)
}

// GetProfileHandler godoc
// @Summary Get user's profile
// @Description Get user's profile
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param profile body models.Profile true "New profile"
// @Success 200 {object} models.HttpBodyInterface{body=models.Profile}
// @failure default {object} models.HttpError
// @Router /users/profile [post]
func (uh *UserHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var userId int64
	uh.logger.WithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	userNew := models.UserData{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&userNew)
	if err != nil {
		uh.logger.Warnf("bad request: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	_, err = govalidator.ValidateStruct(userNew)
	if err != nil {
		uh.logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		return
	}

	profile, err := uh.userUsecase.UpdateProfile(userId, &userNew)
	if err != nil {
		uh.logger.Warnf("can not update profile: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)

	body := models.HttpBodyProfile{Profile: *profile}
	w.Write(models.ToBytes(http.StatusOK, "profile updated", body))
	uh.logger.Debugf("user %d profile updated", userId)
}

// GetProfileHandler godoc
// @Summary Get user's profile
// @Description Get user's profile
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param avatar formData file true "Uploaded avatar"
// @Success 200 {object} models.HttpBodyInterface{body=models.Profile}
// @failure default {object} models.HttpError
// @Router /users/profile/upload [post]
func (uh *UserHandler) UploadProfileImageHandler(w http.ResponseWriter, r *http.Request) {
	uh.logger = uh.logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	defer r.Body.Close()
	err := r.ParseMultipartForm(2 << 20) // 2Мб
	if err != nil {
		uh.logger.Warnf("can not parsemultipart: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.InternalError)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	if len(r.MultipartForm.File["avatar"]) == 0 {
		uh.logger.Warnf("avatar len is 0: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.EmptyImageForm)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	file := r.MultipartForm.File["avatar"][0]
	user, err := uh.userUsecase.UploadAvatar(file, userId)
	if err != nil {
		uh.logger.Warnf("can not upload user %d avatar: %s", userId, err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.EmptyImageForm)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyProfile{Profile: *user.ToProfile()}
	w.Write(models.ToBytes(http.StatusOK, "avatar uploaded successfully", body))
	uh.logger.Debugf("user %d avatar uploaded successfully", userId)
}

// ChangePasswordHandler godoc
// @Summary Change password
// @Description Change password
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param profile body models.ChangePassword true "Change password model"
// @Success 200 {object} models.HttpBodyInterface
// @failure default {object} models.HttpError
// @Router /users/profile/password [post]
func (uh *UserHandler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	uh.logger = uh.logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	changePassword := models.ChangePassword{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&changePassword)
	if err != nil {
		uh.logger.Warnf("bad request: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	_, err = govalidator.ValidateStruct(changePassword)
	if err != nil {
		uh.logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		return
	}

	err = uh.userUsecase.UpdatePassword(userId, &changePassword)
	if err != nil {
		uh.logger.Warnf("password not updated: %s", err.Error())

		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(models.ToBytes(http.StatusOK, "password changed", nil))
	uh.logger.Debugf("user %d changed password successfully", userId)
}
