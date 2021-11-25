package models

import "time"

type Message struct {
	IdFrom    int64     `json:"from" valid:"int"`
	IdTo      int64     `json:"to" valid:"int"`
	IdAdv     int64     `json:"adv" valid:"int"`
	Msg       string    `json:"message" valid:"type(string)"`
	CreatedAt time.Time `json:"created_at" valid:"-" swaggerignore:"true"`
}

type Dialog struct {
	Id1       int64     `json:"user1" valid:"int"`
	Id2       int64     `json:"user2" valid:"int"`
	IdAdv     int64     `json:"adv" valid:"int"`
	CreatedAt time.Time `json:"created_at" valid:"-" swaggerignore:"true"`
}
