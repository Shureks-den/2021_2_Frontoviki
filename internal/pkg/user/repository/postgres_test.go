package repository

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"yula/internal/config"
	"yula/internal/database"
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

	assert.Nil(t, ur)
}
