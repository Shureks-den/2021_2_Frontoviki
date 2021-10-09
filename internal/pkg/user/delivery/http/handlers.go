package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	internalError "yula/internal/error"
	"yula/internal/models"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/session"
	"yula/internal/pkg/user"

	"github.com/asaskevich/govalidator"
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
	r.HandleFunc("/signup", uh.SignUpHandler).Methods(http.MethodPost, http.MethodOptions)

	s := r.PathPrefix("/users").Subrouter()
	s.Handle("/profile", sm.CheckAuthorized(http.HandlerFunc(uh.GetProfileHandler))).Methods(http.MethodGet, http.MethodOptions)
	s.Handle("/profile", sm.CheckAuthorized(http.HandlerFunc(uh.UpdateProfileHandler))).Methods(http.MethodPost, http.MethodOptions)
	s.Handle("/profile/upload", sm.CheckAuthorized(http.HandlerFunc(uh.UploadProfileImageHandler))).Methods(http.MethodPost, http.MethodOptions)
	// r.Handle("profile/upload", sm.CheckAuthorized(http.HandlerFunc(uh.UploadProfileImageHandler))).Methods(http.MethodPost)
	// - пока не работает
}

func (uh *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signUpUser models.UserSignUp

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&signUpUser)
	if err != nil {
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	_, err = govalidator.ValidateStruct(signUpUser)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		return
	}

	user, servErr := uh.userUsecase.Create(&signUpUser)
	if servErr != nil {
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(servErr)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	userSession, err := uh.sessionUsecase.Create(user.Id)
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

	w.Header().Add("Location", r.Host+"/signin") // указываем в качестве перенаправления страницу входа
	w.WriteHeader(http.StatusOK)

	body := models.HttpBodyProfile{Profile: *user.ToProfile()}
	w.Write(models.ToBytes(http.StatusCreated, "user created successfully", body))
}

func (uh *UserHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	profile, err := uh.userUsecase.GetById(userId)
	if err != nil {
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)

	body := models.HttpBodyProfile{Profile: *profile}
	w.Write(models.ToBytes(http.StatusOK, "profile provided", body))
}

func (uh *UserHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	userNew := models.UserData{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&userNew)
	if err != nil {
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	_, err = govalidator.ValidateStruct(userNew)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(http.StatusBadRequest, "invalid data", nil))
		return
	}

	profile, err := uh.userUsecase.UpdateProfile(userId, &userNew)
	if err != nil {
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)

	body := models.HttpBodyProfile{Profile: *profile}
	w.Write(models.ToBytes(http.StatusOK, "profile updated", body))
}

func (uh *UserHandler) UploadProfileImageHandler(w http.ResponseWriter, r *http.Request) {
	var userId int64
	if r.Context().Value(middleware.ContextUserId) != nil {
		userId = r.Context().Value(middleware.ContextUserId).(int64)
	}

	defer r.Body.Close()
	err := r.ParseMultipartForm(2 << 20) // 2Мб
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.InternalError)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	if len(r.MultipartForm.File["avatar"]) == 0 {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.EmptyImageForm)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	file := r.MultipartForm.File["avatar"][0]
	user, err := uh.userUsecase.UploadAvatar(file, userId)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.EmptyImageForm)
		w.Write(models.ToBytes(metaCode, metaMessage, nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyProfile{Profile: *user.ToProfile()}
	w.Write(models.ToBytes(http.StatusOK, "avatar uploaded successfully", body))
}
