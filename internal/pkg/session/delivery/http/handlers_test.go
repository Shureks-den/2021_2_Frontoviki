package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"yula/internal/config"
	"yula/internal/database"
	"yula/internal/models"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"
	userRep "yula/internal/pkg/user/repository"
	userUse "yula/internal/pkg/user/usecase"

	sessRep "yula/internal/pkg/session/repository"
	sessUse "yula/internal/pkg/session/usecase"

	imageloaderRepo "yula/internal/pkg/image_loader/repository"
	imageloaderUse "yula/internal/pkg/image_loader/usecase"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var password = "c0mplex"
var testUser = &models.UserSignUp{
	Email:    "test15@email.com",
	Password: password,
	Name:     "test",
	Surname:  "surtest",
}

func TestSession_SignInHandler_Success(t *testing.T) {
	logger := logging.GetLogger()
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer postgres.Close()

	bytes := bytes.NewReader([]byte(fmt.Sprintf(`
	{
		"email": "test@email.com",
		"password": "%s",
		""
	}
	`, password)))

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)

	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)

	uu := userUse.NewUserUsecase(ur, ilu)
	su := sessUse.NewSessionUsecase(sr)

	_, serverErr := uu.Create(testUser)
	if serverErr != nil {
		t.Fatal()
	}

	// сам тест
	r := httptest.NewRequest("POST", "/signin", bytes)
	w := httptest.NewRecorder()

	sh := NewSessionHandler(su, uu, logger)

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh.Routing(router)

	sh.SignInHandler(w, r)

	var Answer models.HttpError
	err = json.NewDecoder(w.Body).Decode(&Answer)
	if err != nil {
		t.Fatal("invalid serialization")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "signin successfully")
}

func TestSession_SignInHandler_InvalidEmail(t *testing.T) {
	logger := logging.GetLogger()
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer postgres.Close()

	bytes := bytes.NewReader([]byte(fmt.Sprintf(`
	{
		"email": "invalid",
		"password": "%s"
	}
	`, password)))

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)
	uu := userUse.NewUserUsecase(ur, ilu)
	su := sessUse.NewSessionUsecase(sr)

	// сам тест
	r := httptest.NewRequest("POST", "/signin", bytes)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh := NewSessionHandler(su, uu, logger)
	sh.Routing(router)

	sh.SignInHandler(w, r)

	var Answer models.HttpError
	err = json.NewDecoder(w.Body).Decode(&Answer)
	if err != nil {
		t.Fatal("invalid serialization")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, Answer.Code, 404)
	assert.Equal(t, Answer.Message, "user with this email not exist")
}

func TestSession_SignInHandler_InvalidPassword(t *testing.T) {
	logger := logging.GetLogger()
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer postgres.Close()

	bytes := bytes.NewReader([]byte(`
	{
		"email": "test@email.com",
		"password": "baobab"
	}
	`))

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)

	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)

	uu := userUse.NewUserUsecase(ur, ilu)
	su := sessUse.NewSessionUsecase(sr)

	// сам тест
	r := httptest.NewRequest("POST", "/signin", bytes)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh := NewSessionHandler(su, uu, logger)
	sh.Routing(router)

	sh.SignInHandler(w, r)

	var Answer models.HttpError
	err = json.NewDecoder(w.Body).Decode(&Answer)
	if err != nil {
		t.Fatal("invalid serialization")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, Answer.Code, 401)
	assert.Equal(t, Answer.Message, "no rights to access this resource")
}

func TestSession_SignInHandler_InvalidBody(t *testing.T) {
	logger := logging.GetLogger()
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer postgres.Close()

	bytes := bytes.NewReader(nil)

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)

	uu := userUse.NewUserUsecase(ur, ilu)
	su := sessUse.NewSessionUsecase(sr)

	// сам тест
	r := httptest.NewRequest("POST", "/signin", bytes)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh := NewSessionHandler(su, uu, logger)
	sh.Routing(router)

	sh.SignInHandler(w, r)

	var Answer models.HttpError
	err = json.NewDecoder(w.Body).Decode(&Answer)
	if err != nil {
		t.Fatal("invalid serialization")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, Answer.Code, 400)
	assert.Equal(t, Answer.Message, "EOF")
}

func TestSession_LogOutHandler_Success(t *testing.T) {
	logger := logging.GetLogger()
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer postgres.Close()

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)
	uu := userUse.NewUserUsecase(ur, ilu)
	su := sessUse.NewSessionUsecase(sr)

	session, err := su.Create(0)
	assert.Equal(t, err, nil)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    session.Value,
		Expires:  session.ExpiresAt,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}

	r := httptest.NewRequest("POST", "/logout", nil)
	r.AddCookie(&cookie)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh := NewSessionHandler(su, uu, logger)
	sh.Routing(router)

	sh.LogOutHandler(w, r)

	var Answer models.HttpError
	err = json.NewDecoder(w.Body).Decode(&Answer)
	if err != nil {
		t.Fatal("invalid serialization")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, Answer.Code, 200)
	assert.Equal(t, Answer.Message, "logout successfully")
}

func TestSession_LogOutHandler_InvalidName(t *testing.T) {
	logger := logging.GetLogger()
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer postgres.Close()

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)

	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)

	uu := userUse.NewUserUsecase(ur, ilu)
	su := sessUse.NewSessionUsecase(sr)

	session, err := su.Create(0)
	assert.Equal(t, err, nil)
	cookie := http.Cookie{
		Name:     "no_session_id",
		Value:    session.Value,
		Expires:  session.ExpiresAt,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}

	r := httptest.NewRequest("POST", "/logout", nil)
	r.AddCookie(&cookie)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh := NewSessionHandler(su, uu, logger)
	sh.Routing(router)

	sh.LogOutHandler(w, r)

	var Answer models.HttpError
	err = json.NewDecoder(w.Body).Decode(&Answer)
	if err != nil {
		t.Fatal("invalid serialization")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, Answer.Code, 401)
	assert.Equal(t, Answer.Message, "no rights to access this resource")
}

func TestSession_LogOutHandler_InvalidValue(t *testing.T) {
	logger := logging.GetLogger()
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer postgres.Close()

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)

	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)

	uu := userUse.NewUserUsecase(ur, ilu)
	su := sessUse.NewSessionUsecase(sr)

	session, err := su.Create(0)
	assert.Equal(t, err, nil)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    "invalid_value",
		Expires:  session.ExpiresAt,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}

	r := httptest.NewRequest("POST", "/logout", nil)
	r.AddCookie(&cookie)
	w := httptest.NewRecorder()

	router := mux.NewRouter().PathPrefix("/").Subrouter()
	router.Use(middleware.LoggerMiddleware)
	sh := NewSessionHandler(su, uu, logger)
	sh.Routing(router)

	sh.LogOutHandler(w, r)

	var Answer models.HttpError
	err = json.NewDecoder(w.Body).Decode(&Answer)
	if err != nil {
		t.Fatal("invalid serialization")
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, Answer.Code, 500)
	assert.Equal(t, Answer.Message, "something went wrong")
}
