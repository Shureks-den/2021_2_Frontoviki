package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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
	"github.com/stretchr/testify/mock"

	faker "github.com/jaswdr/faker"
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

	// t.Fatal(newStr)

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

func TestSignUpHandlerFailCreateSession(t *testing.T) {
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
	su.On("Create", userCreated.Id).Return(nil, myerr.InternalError).Once()

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(reqUser)
	assert.Nil(t, err)

	reader := bytes.NewReader(reqBodyBuffer.Bytes())
	res, err := http.Post(fmt.Sprintf("%s/signup", srv.URL), http.DetectContentType(reqBodyBuffer.Bytes()), reader)
	assert.Nil(t, err)

	resError := models.HttpError{}
	err = json.NewDecoder(res.Body).Decode(&resError)
	assert.Nil(t, err)

	assert.Equal(t, resError.Code, http.StatusInternalServerError)
}

func TestUpdateProfileSuccess(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	r := mux.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	r.Handle("/users/profile", http.HandlerFunc(uh.UpdateProfileHandler)).Methods(http.MethodPost, http.MethodOptions)

	userNew := models.UserData{
		Id:    0,
		Email: "aboba@baobab.com",
		Name:  "baobab",
	}

	srv := httptest.NewServer(r)
	defer srv.Close()

	uu.On("UpdateProfile", userNew.Id, &userNew).Return(userNew.ToProfile(), nil)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(userNew)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	res, err := http.Post(fmt.Sprintf("%s/users/profile", srv.URL), http.DetectContentType(nil), reader)
	assert.Nil(t, err)

	Answer := models.HttpBodyInterface{}
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "profile updated")
}

func TestUpdateProfileFailUpdate(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	r := mux.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	r.Handle("/users/profile", http.HandlerFunc(uh.UpdateProfileHandler)).Methods(http.MethodPost, http.MethodOptions)

	userNew := models.UserData{
		Id:    0,
		Email: "aboba@baobab.com",
		Name:  "baobab",
	}

	srv := httptest.NewServer(r)
	defer srv.Close()

	uu.On("UpdateProfile", userNew.Id, &userNew).Return(nil, myerr.InternalError)

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(userNew)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	res, err := http.Post(fmt.Sprintf("%s/users/profile", srv.URL), http.DetectContentType(nil), reader)
	assert.Nil(t, err)

	Answer := models.HttpBodyInterface{}
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
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

func TestUploadImageSuccess(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	router := mux.NewRouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/profile/upload", http.HandlerFunc(uh.UploadProfileImageHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	faker := faker.New()
	p := faker.Person()
	fakeimg := p.Image()

	user := models.UserData{
		Id:    0,
		Email: "aboba@baobab.com",
	}

	uu.On("UploadAvatar", mock.AnythingOfType("*multipart.FileHeader"), int64(0)).Return(&user, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("avatar", fakeimg.Name())
	assert.Nil(t, err)
	file, err := os.Open(fakeimg.Name())
	assert.Nil(t, err)
	_, err = io.Copy(fw, file)
	assert.Nil(t, err)
	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/profile/upload", srv.URL), bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "avatar uploaded successfully")
}

func TestUploadImageFailParse(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	router := mux.NewRouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/profile/upload", http.HandlerFunc(uh.UploadProfileImageHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/profile/upload", srv.URL), nil)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}

func TestUploadImageFailUploadAvatar(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	router := mux.NewRouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/profile/upload", http.HandlerFunc(uh.UploadProfileImageHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	faker := faker.New()
	p := faker.Person()
	fakeimg := p.Image()

	uu.On("UploadAvatar", mock.AnythingOfType("*multipart.FileHeader"), int64(0)).Return(nil, myerr.InternalError)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("avatar", fakeimg.Name())
	assert.Nil(t, err)
	file, err := os.Open(fakeimg.Name())
	assert.Nil(t, err)
	_, err = io.Copy(fw, file)
	assert.Nil(t, err)
	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/profile/upload", srv.URL), bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 400)
	assert.Equal(t, Answer.Message, "require image")
}

func TestChangePasswordSuccess(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	router := mux.NewRouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/profile/password", http.HandlerFunc(uh.ChangePasswordHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	changePw := models.ChangePassword{
		Email:       "aboba@baobab.com",
		Password:    "aboba",
		NewPassword: "baobab",
	}

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(changePw)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	uu.On("UpdatePassword", int64(0), &changePw).Return(nil)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/profile/password", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "password changed")
}

func TestChangePasswordFailParse(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	router := mux.NewRouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/profile/password", http.HandlerFunc(uh.ChangePasswordHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/profile/password", srv.URL), nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 400)
	assert.Equal(t, Answer.Message, "bad request")
}

func TestChangePasswordFailUpdate(t *testing.T) {
	su := sessMock.SessionUsecase{}
	uu := userMock.UserUsecase{}
	uh := NewUserHandler(&uu, &su)

	router := mux.NewRouter()
	router.Use(middleware.LoggerMiddleware)
	router.Handle("/profile/password", http.HandlerFunc(uh.ChangePasswordHandler)).Methods(http.MethodPost, http.MethodOptions)

	srv := httptest.NewServer(router)
	defer srv.Close()

	changePw := models.ChangePassword{
		Email:       "aboba@baobab.com",
		Password:    "aboba",
		NewPassword: "baobab",
	}

	reqBodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBuffer).Encode(changePw)
	assert.Nil(t, err)
	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	uu.On("UpdatePassword", int64(0), &changePw).Return(myerr.InternalError)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/profile/password", srv.URL), reader)
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)

	var Answer models.HttpBodyInterface
	err = json.NewDecoder(res.Body).Decode(&Answer)
	assert.Nil(t, err)

	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "internal error")
}
