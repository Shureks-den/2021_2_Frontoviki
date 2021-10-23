package models

import (
	"time"
)

type UserData struct {
	Id        int64     `json:"id" valid:"-"`
	Email     string    `json:"email" valid:"email"`
	Password  string    `json:"password" valid:"type(string),minstringlength(4),optional"`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	Name      string    `json:"name" valid:"type(string),minstringlength(2)"`
	Surname   string    `json:"surname" valid:"type(string),minstringlength(2)"`
	Image     string    `json:"image" valid:"-"`
	Rating    float32   `json:"rating" valid:"range(0|10),optional"`
}

type UserSignIn struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"type(string),minstringlength(4)"`
}

type UserSignUp struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"type(string),minstringlength(4)"`
	Name     string `json:"name" valid:"type(string),minstringlength(2)"`
	Surname  string `json:"surname" valid:"type(string),minstringlength(2)"`
}

type Profile struct {
	Id        int64     `json:"id" valid:"-"`
	Email     string    `json:"email" valid:"email"`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	Name      string    `json:"name" valid:"type(string),minstringlength(2)"`
	Surname   string    `json:"surname" valid:"type(string),minstringlength(2)"`
	Image     string    `json:"image" valid:"-"`
	Rating    float32   `json:"rating" valid:"range(0|10),optional"`
}

func (user *UserData) ToProfile() *Profile {
	return &Profile{
		Id: user.Id, Email: user.Email, CreatedAt: user.CreatedAt, Name: user.Name,
		Surname: user.Surname, Image: user.Image, Rating: user.Rating,
	}
}

type ChangePassword struct {
	Email       string `json:"email" valid:"email"`
	Password    string `json:"password" valid:"type(string),minstringlength(4)"`
	NewPassword string `json:"new_password" valid:"type(string),minstringlength(4)"`
}
