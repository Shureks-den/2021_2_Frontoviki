package usecase

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"yula/internal/config"
	sessRep "yula/internal/pkg/session/repository"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSession_SignInHandler_CreateSuccess(t *testing.T) {
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-4], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()

	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := NewSessionUsecase(sr)

	session, err := su.Create(0)
	assert.Equal(t, nil, err)
	assert.Equal(t, session.UserId, int64(0))
}

func TestSession_SignInHandler_DeleteSuccess(t *testing.T) {
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-4], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()

	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := NewSessionUsecase(sr)

	session, err := su.Create(0)
	assert.Equal(t, nil, err)
	assert.Equal(t, session.UserId, int64(0))

	err = su.Delete(session.Value)
	assert.Equal(t, nil, err)
}

func TestSession_SignInHandler_DeleteError(t *testing.T) {
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-4], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()

	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := NewSessionUsecase(sr)

	err = su.Delete("session_value")
	assert.Equal(t, "empty_row", err.Error())
}

func TestSession_SignInHandler_CheckSuccess(t *testing.T) {
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-4], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()

	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := NewSessionUsecase(sr)

	session, err := su.Create(0)
	assert.Equal(t, nil, err)
	assert.Equal(t, session.UserId, int64(0))

	sessionNew, err := su.Check(session.Value)
	assert.Equal(t, nil, err)
	assert.Equal(t, session.UserId, int64(0))
	assert.Equal(t, session.Value, sessionNew.Value)

	err = su.Delete(session.Value)
	assert.Equal(t, nil, err)
}

func TestSession_SignInHandler_CheckError(t *testing.T) {
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-4], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal(err.Error())
	}

	cnfg := config.NewConfig()

	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	su := NewSessionUsecase(sr)

	sessionNew, err := su.Check("session_value")
	assert.Equal(t, true, sessionNew == nil)
	assert.Equal(t, "empty_row", err.Error())
}
