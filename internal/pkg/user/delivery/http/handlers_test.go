package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"yula/internal/config"
	"yula/internal/database"
	"yula/internal/models"

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

func TestSignUpHandlerValid(t *testing.T) {
	// loads values from .env into the system
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal("No .env file found")
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := sessUse.NewSessionUsecase(sr)
	uu := userUse.NewUserUsecase(ur, ilu)
	uh := NewUserHandler(uu, su)

	r := mux.NewRouter()
	sm := middleware.NewSessionMiddleware(su)
	uh.Routing(r, sm)

	srv := httptest.NewServer(r)
	defer srv.Close()

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TEST.ru",
	}

	reqBodyBuffer := new(bytes.Buffer)
	err = json.NewEncoder(reqBodyBuffer).Encode(reqUser)
	if err != nil {
		t.Fatalf("Bad json encode")
	}

	reader := bytes.NewReader(reqBodyBuffer.Bytes())
	_, err = http.Post(fmt.Sprintf("%s/signup", srv.URL), http.DetectContentType(reqBodyBuffer.Bytes()), reader)
	if err != nil {
		t.Fatalf("Could not post request on signup")
	}

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(res.Body)
	// newStr := buf.String()

	// respUser := models.HttpUser{}
	// err = json.NewDecoder(res.Body).Decode(&respUser)
	// if err != nil {
	// 	t.Fatalf("Could not serialize user from response")
	// }

	// assert.Equal(t, respUser.Body.User.Email, reqUser.Email)
}

func TestSignUpHandlerUserNotValid(t *testing.T) {
	// loads values from .env into the system
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal("No .env file found")
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := sessUse.NewSessionUsecase(sr)
	uu := userUse.NewUserUsecase(ur, ilu)
	uh := NewUserHandler(uu, su)

	r := mux.NewRouter()
	sm := middleware.NewSessionMiddleware(su)
	uh.Routing(r, sm)

	srv := httptest.NewServer(r)
	defer srv.Close()

	res, err := http.Post(fmt.Sprintf("%s/signup", srv.URL), http.DetectContentType(nil), nil)
	if err != nil {
		t.Fatalf("Could not post request on signup")
	}

	resError := models.HttpError{}
	err = json.NewDecoder(res.Body).Decode(&resError)
	if err != nil {
		t.Fatalf("Could not serialize error from response")
	}
	assert.Equal(t, resError.Code, 400)
	assert.Equal(t, resError.Message, "EOF")
}

func TestSignUpHandlerDBConnectionNotOpened(t *testing.T) {
	// loads values from .env into the system
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal("No .env file found")
	}

	cnfg := config.Config{
		DbConfig: config.DatabaseConfig{DatabaseUrl: ""},
		TarantoolCfg: config.TarantoolConfig{
			TarantoolServerAddress: config.GetEnv("TARANTOOL_ADDRESS", "localhost:3302"),
			TarantoolOpts: config.TarantoolOptions{
				User: config.GetEnv("TARANTOOL_USER", "admin"),
				Pass: config.GetEnv("TARANTOOL_PASS", "pass"),
			},
		},
	}
	_, err = database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	assert.NotNil(t, err)
}

func TestSignUpHandlerSameEmail(t *testing.T) {
	// loads values from .env into the system
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal("No .env file found")
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	uu := userUse.NewUserUsecase(ur, ilu)
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := sessUse.NewSessionUsecase(sr)
	uh := NewUserHandler(uu, su)

	r := mux.NewRouter()
	sm := middleware.NewSessionMiddleware(su)
	uh.Routing(r, sm)

	srv := httptest.NewServer(r)
	defer srv.Close()

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TEST.ru",
	}

	reqBodyBuffer := new(bytes.Buffer)
	err = json.NewEncoder(reqBodyBuffer).Encode(reqUser)
	if err != nil {
		t.Fatalf("Bad json encode")
	}

	reader := bytes.NewReader(reqBodyBuffer.Bytes())
	_, err = http.Post(fmt.Sprintf("%s/signup", srv.URL), http.DetectContentType(reqBodyBuffer.Bytes()), reader)
	if err != nil {
		t.Fatalf("Could not post request on signup")
	}

	json.NewEncoder(reqBodyBuffer).Encode(reqUser)
	reader = bytes.NewReader(reqBodyBuffer.Bytes())

	res, err := http.Post(fmt.Sprintf("%s/signup", srv.URL), http.DetectContentType(reqBodyBuffer.Bytes()), reader)
	if err != nil {
		t.Fatalf("Could not post request on signup")
	}

	resError := models.HttpError{}
	err = json.NewDecoder(res.Body).Decode(&resError)
	if err != nil {
		t.Fatalf("Could not serialize error from response")
	}
	assert.Equal(t, resError.Code, 409)
	assert.Equal(t, resError.Message, "user with this email already exist")
}

func TestUpdateProfileHandlerFailToAccessPage(t *testing.T) {
	// loads values from .env into the system
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal("No .env file found")
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := sessUse.NewSessionUsecase(sr)
	uu := userUse.NewUserUsecase(ur, ilu)
	uh := NewUserHandler(uu, su)

	r := mux.NewRouter()
	s := r.PathPrefix("/users").Subrouter()
	s.Handle("/profile", http.HandlerFunc(uh.UpdateProfileHandler))

	srv := httptest.NewServer(r)
	defer srv.Close()

	reqUser := models.UserSignUp{
		Password: "password",
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TEST.ru",
	}

	reqBodyBuffer := new(bytes.Buffer)
	err = json.NewEncoder(reqBodyBuffer).Encode(reqUser)
	if err != nil {
		t.Fatalf("Bad json encode")
	}

	reader := bytes.NewReader(reqBodyBuffer.Bytes())

	res, err := http.Post(fmt.Sprintf("%s/users/profile", srv.URL), http.DetectContentType(reqBodyBuffer.Bytes()), reader)
	if err != nil {
		t.Fatalf("Could not post request on profile")
	}

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(res.Body)
	// newStr := buf.String()

	resError := models.HttpError{}
	err = json.NewDecoder(res.Body).Decode(&resError)
	if err != nil {
		t.Fatalf("Could not serialize error from response")
	}
	assert.Equal(t, resError.Code, 404)
	assert.Equal(t, resError.Message, "user with this email not exist")
}

func TestUpdateProfileHandlerUserNotValid(t *testing.T) {
	// loads values from .env into the system
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal("No .env file found")
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := sessUse.NewSessionUsecase(sr)
	uu := userUse.NewUserUsecase(ur, ilu)
	uh := NewUserHandler(uu, su)

	r := mux.NewRouter()
	s := r.PathPrefix("/users").Subrouter()
	s.Handle("/profile", http.HandlerFunc(uh.UpdateProfileHandler))

	srv := httptest.NewServer(r)
	defer srv.Close()

	res, err := http.Post(fmt.Sprintf("%s/users/profile", srv.URL), http.DetectContentType(nil), nil)
	if err != nil {
		t.Fatalf("Could not post request on profile")
	}
	resError := models.HttpError{}
	err = json.NewDecoder(res.Body).Decode(&resError)
	if err != nil {
		t.Fatalf("Could not serialize error from response")
	}
	assert.Equal(t, resError.Code, 400)
	assert.Equal(t, resError.Message, "EOF")
}

func TestGetProfileHandlerFailToAccessPage(t *testing.T) {
	// loads values from .env into the system
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-5], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal("No .env file found")
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ilr := imageloaderRepo.NewImageLoaderRepository()
	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := sessUse.NewSessionUsecase(sr)
	uu := userUse.NewUserUsecase(ur, ilu)
	uh := NewUserHandler(uu, su)

	r := mux.NewRouter()

	s := r.PathPrefix("/users").Subrouter()
	s.Handle("/profile", http.HandlerFunc(uh.GetProfileHandler))

	srv := httptest.NewServer(r)
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s/users/profile", srv.URL))
	if err != nil {
		t.Fatalf("Could not post request on profile")
	}

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(res.Body)
	// newStr := buf.String()

	resError := models.HttpError{}
	err = json.NewDecoder(res.Body).Decode(&resError)
	if err != nil {
		t.Fatalf("Could not serialize error from response")
	}
	assert.Equal(t, resError.Code, 404)
	assert.Equal(t, resError.Message, "user with this email not exist")

}
