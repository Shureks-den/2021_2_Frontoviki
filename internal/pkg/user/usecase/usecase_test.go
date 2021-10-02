package usecase

/*
import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
	"yula/internal/codes"
	"yula/internal/config"
	"yula/internal/database"
	"yula/internal/models"

	userRep "yula/internal/pkg/user/repository"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
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
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	uu := NewUserUsecase(ur)

	reqUser := models.UserSignUp{
		Username: "username",
		Password: "password",
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + fmt.Sprint(rand.Int()) + "@TEST.ru",
	}

	createdUser, error := uu.Create(&reqUser)
	if error != nil {
		t.Fatalf(error.Message)
	}

	assert.Equal(t, reqUser.Email, createdUser.Email)
}

func TestTwiceCreate(t *testing.T) {
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
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	uu := NewUserUsecase(ur)

	reqUser := models.UserSignUp{
		Username: "username",
		Password: "password",
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TEST.ru",
	}

	_, error := uu.Create(&reqUser)
	if error != nil {
		t.Fatalf(error.Message)
	}
	_, error = uu.Create(&reqUser)

	assert.Equal(t, error, codes.StatusMap[codes.UserAlreadyExist])
}

func TestGetByEmail(t *testing.T) {
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
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	uu := NewUserUsecase(ur)

	reqUser := models.UserSignUp{
		Username: "username",
		Password: "password",
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TEST.ru",
	}

	created_user, error := uu.Create(&reqUser)
	if error != nil {
		t.Fatalf(error.Message)
	}
	user, error := uu.GetByEmail(created_user.Email)
	if error != nil {
		t.Fatalf(error.Message)
	}

	assert.Equal(t, user.Email, created_user.Email)
}

func TestGetByEmailUserNotExist(t *testing.T) {
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
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	uu := NewUserUsecase(ur)

	reqUser := models.UserSignUp{
		Username: "username",
		Password: "password",
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TEST.ru",
	}

	_, error := uu.GetByEmail(reqUser.Email)
	assert.Equal(t, error, codes.StatusMap[codes.UserNotExist])
}

func TestCheckPassword(t *testing.T) {
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
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	uu := NewUserUsecase(ur)

	reqUser := models.UserSignUp{
		Username: "username",
		Password: "password",
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TEST.ru",
	}

	created_user, error := uu.Create(&reqUser)
	if error != nil {
		t.Fatalf(error.Message)
	}

	var srverr *codes.ServerError = nil
	error = uu.CheckPassword(created_user, reqUser.Password)
	assert.Equal(t, error, srverr)
}

func TestCheckPasswordInvalid(t *testing.T) {
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
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	uu := NewUserUsecase(ur)

	reqUser := models.UserSignUp{
		Username: "username",
		Password: "password",
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TEST.ru",
	}

	created_user, error := uu.Create(&reqUser)
	if error != nil {
		t.Fatalf(error.Message)
	}

	error = uu.CheckPassword(created_user, reqUser.Password+"0104")
	assert.Equal(t, error, codes.StatusMap[codes.Unauthorized])
}

func TestGetById(t *testing.T) {
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
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	uu := NewUserUsecase(ur)

	reqUser := models.UserSignUp{
		Username: "username",
		Password: "password",
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TEST.ru",
	}

	created_user, error := uu.Create(&reqUser)
	if error != nil {
		t.Fatalf(error.Message)
	}
	user, error := uu.GetById(created_user.Id)
	if error != nil {
		t.Fatalf(error.Message)
	}

	assert.Equal(t, user.Email, created_user.Email)
}

func TestGetByIdUserNotExist(t *testing.T) {
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
		t.Fatal("Connection not opened")
	}
	defer postgres.Close()

	ur := userRep.NewUserRepository(postgres.GetDbPool())
	uu := NewUserUsecase(ur)

	reqUser := models.UserSignUp{
		Username: "username",
		Password: "password",
		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TEST.ru",
	}

	_, error := uu.GetById(reqUser.Id)
	assert.Equal(t, error, codes.StatusMap[codes.UserNotExist])
}*/

// func TestUpdateUserProfile(t *testing.T) {
// 	// loads values from .env into the system
// 	pwd, err := os.Getwd()
// 	folders := strings.Split(pwd, "/")
// 	pwd = strings.Join(folders[:len(folders)-4], "/")
// 	fmt.Println(pwd, err)

// 	if err := godotenv.Load(pwd + "/.env"); err != nil {
// 		t.Fatal("No .env file found")
// 	}

// 	cnfg := config.NewConfig()
// 	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
// 	if err != nil {
// 		t.Fatal("Connection not opened")
// 	}
// 	defer postgres.Close()

// 	ur := userRep.NewUserRepository(postgres.GetDbPool())
// 	uu := NewUserUsecase(ur)

// 	reqUser := models.UserSignUp{
// 		Username: "username",
// 		Password: "password",
// 		Email:    fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int()) + "@TEST.ru",
// 	}

// 	createdUser, error := uu.Create(&reqUser)
// 	if error != nil {
// 		t.Fatalf(error.Message)
// 	}

// 	createdUser.Password = "wrmkgwprg"
// 	error = uu.UpdateProfile(createdUser.Id, createdUser)
// 	if error != nil {
// 		t.Fatalf(error.Message)
// 	}

// 	var srverr *codes.ServerError = nil
// 	error = uu.CheckPassword(createdUser, "wrmkgwprg")
// 	assert.Equal(t, error, srverr)
// }
