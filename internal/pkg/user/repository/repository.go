package repository

/*
import (
	"errors"
	"fmt"
	"yula/internal/models"
	"yula/internal/pkg/user"

	"github.com/google/uuid"
)

type UserRepository struct {
	users []models.UserData
}

func NewUserRepository(users []models.UserData) user.UserRepository {
	return &UserRepository{
		users: users,
	}
}

func (ur *UserRepository) Insert(user *models.UserData) error {
	ur.users = append(ur.users, *user)

	// проверка что вставка произошла - вывод всех юзеров
	ur.PrintAllUsers()
	return nil
}

func (ur *UserRepository) SelectByEmail(email string) (*models.UserData, error) {
	for _, usr := range ur.users {
		if usr.Email == email {
			return &usr, nil
		}
	}
	return nil, errors.New("user not exist")
}

func (ur *UserRepository) SelectById(userId uuid.UUID) (*models.UserData, error) {
	for _, usr := range ur.users {
		if usr.Id == userId {
			return &usr, nil
		}
	}
	return nil, errors.New("user not exist")
}

func (ur *UserRepository) Update(user *models.UserData) error {
	_, err := ur.SelectById(user.Id)

	if err != nil {
		return errors.New("user not exist")
	}

	for ind, usr := range ur.users {
		if usr.Id == user.Id {
			ur.users[ind] = *user
		}
	}

	return nil
}

func (ur *UserRepository) PrintAllUsers() {
	for i, usr := range ur.users {
		fmt.Printf("User %d:\n\tId: %s\n\tUsername: %s\n\tEmail: %s\n\tPassword: %s\n\tCreated at: %s\n\n",
			i, usr.Id, usr.Username, usr.Email, usr.Password, usr.CreatedAt)
	}
}
*/
