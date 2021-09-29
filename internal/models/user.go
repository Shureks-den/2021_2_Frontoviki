package models

import (
	"time"
)

type UserData struct {
	// main info
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`

	// additional info
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Image   string `json:"image"`
}

type UserSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignUp struct {
	Id        int64     `json:"-"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
}

type UserWithoutPassword struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (user *UserData) RemovePassword() UserWithoutPassword {
	return UserWithoutPassword{
		Id: user.Id, Username: user.Username, Email: user.Email, CreatedAt: user.CreatedAt,
	}
}

type Profile struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`

	Name    string `json:"name"`
	Surname string `json:"surname"`
	Image   string `json:"image"`
}

func (user *UserData) ToProfile() *Profile {
	return &Profile{
		Id: user.Id, Username: user.Username, Email: user.Email, CreatedAt: user.CreatedAt,
		Name: user.Name, Surname: user.Surname, Image: user.Image,
	}
}
