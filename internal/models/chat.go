package models

import "time"

type Message struct {
	UserFrom  int64     `json:"from" valid:"int"`
	UserTo    int64     `json:"to" valid:"int"`
	Message   string    `json:"message" valid:"type(string)"`
	CreatedAt time.Time `json:"created_at" valid:"-" swaggerignore:"true"`
	UpdatedAt time.Time `json:"created_at" valid:"-" swaggerignore:"true"`
}
