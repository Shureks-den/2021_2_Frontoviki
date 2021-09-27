package repository

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"yula/internal/config"
	"yula/internal/database"
	"yula/internal/models"
)

func TestInit(t *testing.T) {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer postgres.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	ur := NewUserRepository(postgres.GetDbPool())

	assert.NotNil(t, ur)
}

func TestInsertSelect(t *testing.T) {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer postgres.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	ur := NewUserRepository(postgres.GetDbPool())
	assert.NotNil(t, ur)

	ud := models.UserData{
		Email:    "stringwfwrwf@mail.ru",
		Password: "3191031",
	}

	ur.Insert(&ud)
	user, _ := ur.SelectByEmail(ud.Email)

	assert.Equal(t, user.Email, ud.Email)
	assert.Equal(t, user.Id, ud.Id)
}

func TestUpdate(t *testing.T) {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer postgres.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	ur := NewUserRepository(postgres.GetDbPool())
	assert.NotNil(t, ur)

	ud := models.UserData{
		Email:    "stringwfwrwf@mail.ru",
		Password: "3191031",
	}

	ur.Insert(&ud)
	ud.Email = "eoqoepqpqeoqoqepqe@mail.ru"
	ur.Update(&ud)

	user, _ := ur.SelectById(ud.Id)

	assert.Equal(t, user.Email, ud.Email)
	assert.Equal(t, user.Id, ud.Id)
}
