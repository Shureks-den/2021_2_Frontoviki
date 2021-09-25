package models

import (
	"time"
)

type UserData struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
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
