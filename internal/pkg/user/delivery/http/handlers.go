package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/user"
	proto "yula/proto/generated/auth"

	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUsecase    user.UserUsecase
	sessionUsecase proto.AuthClient
}

func NewUserHandler(userUsecase user.UserUsecase, sessionUsecase proto.AuthClient) *UserHandler {
	return &UserHandler{
		userUsecase:    userUsecase,
		sessionUsecase: sessionUsecase,
	}
}

func (uh *UserHandler) Routing(r *mux.Router, sm *middleware.SessionMiddleware) {
	r.HandleFunc("/signup", uh.SignUpHandler).Methods(http.MethodPost, http.MethodOptions)

	s := r.PathPrefix("/users").Subrouter()
	s.Handle("/profile", middleware.SetSCRFToken(sm.CheckAuthorized(http.HandlerFunc(uh.GetProfileHandler)))).Methods(http.MethodGet, http.MethodOptions)
	s.Handle("/profile", sm.CheckAuthorized(http.HandlerFunc(uh.UpdateProfileHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/profile/upload", sm.CheckAuthorized(http.HandlerFunc(uh.UploadProfileImageHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/profile/password", sm.CheckAuthorized(http.HandlerFunc(uh.ChangePasswordHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/profile/rating", sm.CheckAuthorized(http.HandlerFunc(uh.RatingHandler))).Methods(http.MethodPost, http.MethodOptions)
}

var (
	logger logging.Logger = logging.GetLogger()
)

// SignUpHandler godoc
// @Summary Sign up
// @Description Sign up
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param user body models.UserSignUp true "User sign up data"
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyProfile}
// @failure default {object} models.HttpError
// @Router /signup [post]
func (uh *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signUpUser models.UserSignUp
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warnf("cannot convert body to bytes: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = easyjson.Unmarshal(buf, &signUpUser)
	if err != nil {
		logger.Warnf("cannot unmarshal: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	sanitizer := bluemonday.UGCPolicy()
	signUpUser.Email = sanitizer.Sanitize(signUpUser.Email)
	signUpUser.Password = sanitizer.Sanitize(signUpUser.Password)
	signUpUser.Name = sanitizer.Sanitize(signUpUser.Name)
	signUpUser.Surname = sanitizer.Sanitize(signUpUser.Surname)

	_, err = govalidator.ValidateStruct(signUpUser)
	if err != nil {
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(signUpUser)
		if err != nil {
			logger.Warnf("trouble with encoder %s", err.Error())
		}
		logger.Warnf("invalid data: %s", buf.String())

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}
	user, servErr := uh.userUsecase.Create(&signUpUser)
	if servErr != nil {
		logger.Warnf("can not create user: %s", servErr.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(servErr)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	protoUserSession, err := uh.sessionUsecase.Create(context.Background(), &proto.UserID{ID: user.Id})
	if err != nil {
		logger.Warnf("can not create session based on user %d: %s", user.Id, err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}
	userSession := models.Session{
		Value:     protoUserSession.SessionID,
		ExpiresAt: protoUserSession.ExpireAt.AsTime(),
		UserId:    protoUserSession.UserID,
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    userSession.Value,
		Expires:  userSession.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Secure:   true,
	})

	w.Header().Add("Location", r.Host+"/signin") // указываем в качестве перенаправления страницу входа
	w.WriteHeader(http.StatusOK)

	body := models.HttpBodyProfile{Profile: *user.ToProfile()}
	_, err = w.Write(models.ToBytes(http.StatusCreated, "user created successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
	logger.Debugf("user %d created successfully", user.Id)
}

// GetProfileHandler godoc
// @Summary Get user's profile
// @Description Get user's profile
// @Tags user
// @Accept application/json
// @Produce application/json
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyProfile}
// @failure default {object} models.HttpError
// @Router /users/profile [get]
func (uh *UserHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	var userId int64 = -1
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	profile, err := uh.userUsecase.GetById(userId)
	if err != nil {
		if userId != int64(-1) {
			logger.Warnf("can not get user with id %d: %s", userId, err.Error())
		}
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	rateStat, err := uh.userUsecase.GetRating(userId, userId)
	if err != nil {
		logger.Debugf("can not get user's statistic with id %d: %s", userId, err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyProfile{Profile: *profile, Rating: *rateStat}
	_, err = w.Write(models.ToBytes(http.StatusOK, "profile provided", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
	logger.Debugf("user %d got successfully", userId)
}

// GetProfileHandler godoc
// @Summary Get user's profile
// @Description Get user's profile
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param profile body models.Profile true "New profile"
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyProfile}
// @failure default {object} models.HttpError
// @Router /users/profile [post]
func (uh *UserHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var userId int64
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	userNew := models.UserData{}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warnf("cannot convert body to bytes: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = easyjson.Unmarshal(buf, &userNew)
	if err != nil {
		logger.Warnf("cannot unmarshal: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	sanitizer := bluemonday.UGCPolicy()
	userNew.Email = sanitizer.Sanitize(userNew.Email)
	userNew.Password = sanitizer.Sanitize(userNew.Password)
	userNew.Name = sanitizer.Sanitize(userNew.Name)
	userNew.Surname = sanitizer.Sanitize(userNew.Surname)
	userNew.Image = sanitizer.Sanitize(userNew.Image)

	_, err = govalidator.ValidateStruct(userNew)
	if err != nil {
		logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	profile, err := uh.userUsecase.UpdateProfile(userId, &userNew)
	if err != nil {
		logger.Warnf("can not update profile: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	rateStat, err := uh.userUsecase.GetRating(userId, userId)
	if err != nil {
		logger.Debugf("can not get user's statistic with id %d: %s", userId, err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyProfile{Profile: *profile, Rating: *rateStat}
	_, err = w.Write(models.ToBytes(http.StatusOK, "profile updated", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
	logger.Debugf("user %d profile updated", userId)
}

// GetProfileHandler godoc
// @Summary Get user's profile
// @Description Get user's profile
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param avatar formData file true "Uploaded avatar"
// @Success 200 {object} models.HttpBodyInterface{body=models.HttpBodyProfile}
// @failure default {object} models.HttpError
// @Router /users/profile/upload [post]
func (uh *UserHandler) UploadProfileImageHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	defer r.Body.Close()
	err := r.ParseMultipartForm(2 << 20) // 2Мб
	if err != nil {
		logger.Warnf("can not parsemultipart: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.InternalError)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	if len(r.MultipartForm.File["avatar"]) == 0 {
		logger.Warnf("avatar len is 0: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.EmptyImageForm)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	file := r.MultipartForm.File["avatar"][0]
	user, err := uh.userUsecase.UploadAvatar(file, userId)
	if err != nil {
		logger.Warnf("can not upload user %d avatar: %s", userId, err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.EmptyImageForm)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	rateStat, err := uh.userUsecase.GetRating(userId, userId)
	if err != nil {
		logger.Debugf("can not get user's statistic with id %d: %s", userId, err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyProfile{Profile: *user.ToProfile(), Rating: *rateStat}
	_, err = w.Write(models.ToBytes(http.StatusOK, "avatar uploaded successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
	logger.Debugf("user %d avatar uploaded successfully", userId)
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
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	changePassword := models.ChangePassword{}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warnf("cannot convert body to bytes: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = easyjson.Unmarshal(buf, &changePassword)
	if err != nil {
		logger.Warnf("cannot unmarshal: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	sanitizer := bluemonday.UGCPolicy()
	changePassword.Email = sanitizer.Sanitize(changePassword.Email)
	changePassword.Password = sanitizer.Sanitize(changePassword.Password)
	changePassword.NewPassword = sanitizer.Sanitize(changePassword.NewPassword)

	_, err = govalidator.ValidateStruct(changePassword)
	if err != nil {
		logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = uh.userUsecase.UpdatePassword(userId, &changePassword)
	if err != nil {
		logger.Warnf("password not updated: %s", err.Error())

		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, "password changed", nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}

	logger.Debugf("user %d changed password successfully", userId)
}

// RatingHandler godoc
// @Summary Rate users
// @Description Rate users
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param body body models.Rating true "Rate user model"
// @Success 200 {object} models.HttpBodyInterface
// @failure default {object} models.HttpError
// @Router /users/profile/rating [post]
func (uh *UserHandler) RatingHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	defer r.Body.Close()
	inputRating := &models.Rating{}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warnf("cannot convert body to bytes: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	err = easyjson.Unmarshal(buf, inputRating)
	if err != nil {
		logger.Warnf("cannot unmarshal: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	_, err = govalidator.ValidateStruct(inputRating)
	if err != nil {
		logger.Warnf("invalid data: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	inputRating.UserFrom = userId
	err = uh.userUsecase.SetRating(inputRating)
	if err != nil {
		logger.Warnf("cannot set rating: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, "user appreciated", nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
}
