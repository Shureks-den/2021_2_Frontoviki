package repository

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"yula/internal/config"
	"yula/internal/database"
	"yula/internal/models"
)

func TestInit(t *testing.T) {
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-4], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal("No .env file found")
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer postgres.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	ur := NewUserRepository(postgres.GetDbPool())

	assert.NotNil(t, ur)
}

func TestInsertSelect(t *testing.T) {
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-4], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal("No .env file found")
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer postgres.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	ur := NewUserRepository(postgres.GetDbPool())
	assert.NotNil(t, ur)

	ud := models.UserData{
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TESTINSERTSELECT.ru",
		Password: "3191022331",
	}

	ur.Insert(&ud)
	user, _ := ur.SelectByEmail(ud.Email)

	assert.Equal(t, user.Email, ud.Email)
	assert.Equal(t, user.Id, ud.Id)
}

func TestUpdate(t *testing.T) {
	// loads values from .env into the system
	pwd, err := os.Getwd()
	folders := strings.Split(pwd, "/")
	pwd = strings.Join(folders[:len(folders)-4], "/")
	fmt.Println(pwd, err)

	if err := godotenv.Load(pwd + "/.env"); err != nil {
		t.Fatal("No .env file found")
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer postgres.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	ur := NewUserRepository(postgres.GetDbPool())
	assert.NotNil(t, ur)

	ud := models.UserData{
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TESTUPDATE.ru",
		Password: "3191031",
	}

	ur.Insert(&ud)
	ud.Email = fmt.Sprint(time.Now().Unix()) + "TESTUPDATE_aboba@mail.ru"
	ur.Update(&ud)

	user, _ := ur.SelectById(ud.Id)

	assert.Equal(t, user.Email, ud.Email)
	assert.Equal(t, user.Id, ud.Id)
}
