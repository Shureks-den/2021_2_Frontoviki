package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"yula/internal/models"
	"yula/internal/pkg/middleware"

	myerr "yula/internal/error"

	userMock "yula/internal/pkg/user/mocks"

	sessMock "yula/internal/pkg/session/mocks"

	imageloader "yula/internal/pkg/image_loader"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var password string = "c0mplex"
var testUser = &models.UserSignUp{
	Email:    "tes5trtrgrwt4145@email.com",
	Password: password,
	Name:     "test",
	Surname:  "surtest",
}

func TestSession_SignInHandler_Success(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	sh := NewSessionHandler(&su, &uu)

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh.Routing(router)

	srv := httptest.NewServer(router)
	defer srv.Close()

	reqUser := models.UserSignIn{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	user := models.UserData{
		Id:        258,
		Email:     reqUser.Email,
		Password:  "aboba",
		CreatedAt: time.Now(),
		Image:     imageloader.DefaultAdvertImage,
		Rating:    0,
	}

	sessionCreated := models.Session{
		Value:     uuid.NewString(),
		UserId:    user.Id,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	uu.On("GetByEmail", reqUser.Email).Return(&user, nil)
	uu.On("CheckPassword", &user, reqUser.Password).Return(nil)

	su.On("Create", user.Id).Return(&sessionCreated, nil).Once()

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(reqUser)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/signin", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "signin successfully")
}

func TestSession_SignInHandler_InvalidEmail(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	sh := NewSessionHandler(&su, &uu)

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh.Routing(router)

	srv := httptest.NewServer(router)
	defer srv.Close()

	reqUser := models.UserSignIn{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	uu.On("GetByEmail", reqUser.Email).Return(nil, myerr.NotExist)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(reqUser)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/signin", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 404)
	assert.Equal(t, Answer.Message, "not exist")
}

func TestSession_SignInHandler_InvalidPassword(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	sh := NewSessionHandler(&su, &uu)

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh.Routing(router)

	srv := httptest.NewServer(router)
	defer srv.Close()

	reqUser := models.UserSignIn{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	user := models.UserData{
		Id:        258,
		Email:     reqUser.Email,
		Password:  "aboba",
		CreatedAt: time.Now(),
		Image:     imageloader.DefaultAdvertImage,
		Rating:    0,
	}

	uu.On("GetByEmail", reqUser.Email).Return(&user, nil)
	uu.On("CheckPassword", &user, reqUser.Password).Return(myerr.PasswordMismatch)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(reqUser)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/signin", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 401)
	assert.Equal(t, Answer.Message, "password mismatch")
}

func TestSession_SignInHandler_InvalidBody(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	sh := NewSessionHandler(&su, &uu)

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh.Routing(router)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/signin", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, http.StatusInternalServerError)
}

func TestSession_LogOutHandler_Success(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	sh := NewSessionHandler(&su, &uu)

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh.Routing(router)

	srv := httptest.NewServer(router)
	defer srv.Close()

	session := models.Session{Value: uuid.NewString(), UserId: 255159, ExpiresAt: time.Now().Add(time.Hour)}
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    session.Value,
		Expires:  session.ExpiresAt,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}

	su.On("Delete", cookie.Value).Return(nil)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/logout", srv.URL), nil)
	assert.Nil(t, err)

	req.AddCookie(&cookie)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpError
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "logout successfully")
}

func TestSession_LogOutHandler_InvalidName(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	sh := NewSessionHandler(&su, &uu)

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh.Routing(router)

	srv := httptest.NewServer(router)
	defer srv.Close()

	session := models.Session{Value: uuid.NewString(), UserId: 255159, ExpiresAt: time.Now().Add(time.Hour)}
	cookie := http.Cookie{
		Name:     "not_session_id",
		Value:    session.Value,
		Expires:  session.ExpiresAt,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/logout", srv.URL), nil)
	assert.Nil(t, err)

	req.AddCookie(&cookie)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpError
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 401)
	assert.Equal(t, Answer.Message, "unauthorized")
}

func TestSession_LogOutHandler_InvalidValue(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	sh := NewSessionHandler(&su, &uu)

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh.Routing(router)

	srv := httptest.NewServer(router)
	defer srv.Close()

	session := models.Session{Value: uuid.NewString(), UserId: 255159, ExpiresAt: time.Now().Add(time.Hour)}
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    "aboba",
		Expires:  session.ExpiresAt,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}

	su.On("Delete", cookie.Value).Return(myerr.NotExist)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/logout", srv.URL), nil)
	assert.Nil(t, err)

	req.AddCookie(&cookie)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpError
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 404)
	assert.Equal(t, Answer.Message, "not exist")
}
