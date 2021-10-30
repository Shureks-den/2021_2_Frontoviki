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

	myerr "yula/internal/error"
	"yula/internal/pkg/middleware"
	userMock "yula/internal/pkg/user/mocks"

	sessMock "yula/internal/pkg/session/mocks"

	imageloader "yula/internal/pkg/image_loader"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSignUpHandlerValid(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	r := mux.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	sm := middleware.NewSessionMiddleware(&su)
	uh.Routing(r, sm)

	srv := httptest.NewServer(r)
	defer srv.Close()

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(reqUser)
	assert.Nil(t, err)

	userCreated := models.UserData{
		Id:        258,
		Email:     reqUser.Email,
		Password:  "aboba",
		CreatedAt: time.Now(),
		Image:     imageloader.DefaultAdvertImage,
		Rating:    0,
	}
	uu.On("Create", &reqUser).Return(&userCreated, nil).Once()

	sessionCreated := models.Session{
		Value:     uuid.NewString(),
		UserId:    userCreated.Id,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	su.On("Create", userCreated.Id).Return(&sessionCreated, nil).Once()

	reader := bytes.NewReader(reqBodyBuffer.Bytes())
	res, err := http.Post(fmt.Sprintf("%s/signup", srv.URL), http.DetectContentType(reqBodyBuffer.Bytes()), reader)
	assert.Nil(t, err)

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(res.Body)
	// newStr := buf.String()

	// t.Fatal(res.Body)

	profile := models.Profile{}
	resp := models.HttpBodyInterface{
		Body: profile,
	}
	err = json.NewDecoder(res.Body).Decode(&resp)
	assert.Nil(t, err)

	assert.Equal(t, (((resp.Body.(map[string]interface{}))["profile"]).(map[string]interface{}))["email"], userCreated.ToProfile().Email)
}

func TestSignUpHandlerUserNotValid(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	r := mux.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	sm := middleware.NewSessionMiddleware(&su)
	uh.Routing(r, sm)

	srv := httptest.NewServer(r)
	defer srv.Close()

	res, err := http.Post(fmt.Sprintf("%s/signup", srv.URL), http.DetectContentType(nil), nil)
	assert.Nil(t, err)

	decodedRes := models.HttpError{}
	err = json.NewDecoder(res.Body).Decode(&decodedRes)
	assert.Nil(t, err)

	assert.Equal(t, decodedRes.Code, http.StatusBadRequest)
}

func TestSignUpHandlerSameEmail(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	r := mux.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	sm := middleware.NewSessionMiddleware(&su)
	uh.Routing(r, sm)

	srv := httptest.NewServer(r)
	defer srv.Close()

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    "superchel@shibanov.jp",
	}

	userCreated := models.UserData{
		Id:        258,
		Email:     reqUser.Email,
		Password:  "aboba",
		CreatedAt: time.Now(),
		Image:     imageloader.DefaultAdvertImage,
		Rating:    0,
	}
	uu.On("Create", &reqUser).Return(&userCreated, nil).Once()
	uu.On("Create", &reqUser).Return(nil, myerr.AlreadyExist)

	sessionCreated := models.Session{
		Value:     uuid.NewString(),
		UserId:    userCreated.Id,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	su.On("Create", userCreated.Id).Return(&sessionCreated, nil).Once()

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(reqUser)
	assert.Nil(t, err)

	reader := bytes.NewReader(reqBodyBuffer.Bytes())
	_, err = http.Post(fmt.Sprintf("%s/signup", srv.URL), http.DetectContentType(reqBodyBuffer.Bytes()), reader)
	assert.Nil(t, err)

	reqBodyBuffer = new(bytes.Buffer)
	err = json.NewEncoder(reqBodyBuffer).Encode(reqUser)
	assert.Nil(t, err)

	reader = bytes.NewReader(reqBodyBuffer.Bytes())
	res, err := http.Post(fmt.Sprintf("%s/signup", srv.URL), http.DetectContentType(reqBodyBuffer.Bytes()), reader)
	assert.Nil(t, err)

	resError := models.HttpError{}
	err = json.NewDecoder(res.Body).Decode(&resError)
	assert.Nil(t, err)

	assert.Equal(t, resError.Code, http.StatusForbidden)
}

func TestUpdateProfileHandlerUserNotValid(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	r := mux.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	s := r.PathPrefix("/users").Subrouter()
	s.Handle("/profile", http.HandlerFunc(uh.UpdateProfileHandler))

	srv := httptest.NewServer(r)
	defer srv.Close()

	res, err := http.Post(fmt.Sprintf("%s/users/profile", srv.URL), http.DetectContentType(nil), nil)
	assert.Nil(t, err)

	resError := models.HttpError{}
	err = json.NewDecoder(res.Body).Decode(&resError)
	assert.Nil(t, err)

	assert.Equal(t, resError.Code, http.StatusBadRequest)
}

func TestGetProfileHandlerFailToAccessPage(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	r := mux.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	s := r.PathPrefix("/users").Subrouter()
	s.Handle("/profile", http.HandlerFunc(uh.GetProfileHandler))

	srv := httptest.NewServer(r)
	defer srv.Close()

	uu.On("GetById", int64(-1)).Return(nil, myerr.NotExist)

	res, err := http.Get(fmt.Sprintf("%s/users/profile", srv.URL))
	assert.Nil(t, err)

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(res.Body)
	// newStr := buf.String()

	resError := models.HttpError{}
	err = json.NewDecoder(res.Body).Decode(&resError)
	assert.Nil(t, err)

	assert.Equal(t, resError.Code, http.StatusNotFound)
}
