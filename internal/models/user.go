package models

import (
	"time"
)

type UserData struct {
	Id        int64     `json:"id" valid:"-"`
	Email     string    `json:"email" valid:"email"`
	Phone     string    `json:"phone" valid:"stringlength(11|11),optional"`
	Password  string    `json:"password" valid:"type(string),minstringlength(4),optional"`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	Name      string    `json:"name" valid:"type(string),minstringlength(2)"`
	Surname   string    `json:"surname" valid:"type(string),minstringlength(2)"`
	Image     string    `json:"image" valid:"-"`
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
	Phone     string    `json:"phone" valid:"stringlength(11|11)"`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	Name      string    `json:"name" valid:"type(string),minstringlength(2)"`
	Surname   string    `json:"surname" valid:"type(string),minstringlength(2)"`
	Image     string    `json:"image" valid:"-"`
}

func (user *UserData) ToProfile() *Profile {
	return &Profile{
		Id: user.Id, Email: user.Email, Phone: user.Phone,
		CreatedAt: user.CreatedAt, Name: user.Name,
		Surname: user.Surname, Image: user.Image,
	}
}

type ChangePassword struct {
	Email       string `json:"email" valid:"email"`
	Password    string `json:"password" valid:"type(string),minstringlength(4)"`
	NewPassword string `json:"new_password" valid:"type(string),minstringlength(4)"`
}

type Rating struct {
	UserFrom int64 `json:"from" valid:"optional"`
	UserTo   int64 `json:"to" valid:"int"`
	Rating   int   `json:"rating" valid:"range(0|5),optional"`
}

type RatingStat struct {
	RatingSum    int64   `json:"sum" valid:"type(int)"`
	RatingCount  int64   `json:"count" valid:"type(int)"`
	RatingAvg    float32 `json:"avg" valid:"type(float)"`
	PersonalRate int     `json:"rate" valid:"range(0|5)"`
	IsRated      bool    `json:"is_rated" valid:"required"`
}
