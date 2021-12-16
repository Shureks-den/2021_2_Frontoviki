package models

import (
	"time"
)

type IMessage struct {
	IdFrom int64 `json:"from" valid:"int"`
	IdTo   int64 `json:"to" valid:"int"`
	IdAdv  int64 `json:"adv" valid:"int"`
}

type Message struct {
	MI IMessage `json:"info"`

	Msg string `json:"message" valid:"type(string)"`

	CreatedAt time.Time `json:"created_at" valid:"-" swaggerignore:"true"`
}

func (iMsg *IMessage) ToMessage(Msg string, CreatedAt time.Time) *Message {
	return &Message{
		MI:        *iMsg,
		Msg:       Msg,
		CreatedAt: CreatedAt,
	}
}

func (msg *Message) ToIMessage() *IMessage {
	return &msg.MI
}

type IDialog struct {
	Id1   int64 `json:"user1" valid:"int"`
	Id2   int64 `json:"user2" valid:"int"`
	IdAdv int64 `json:"adv" valid:"int"`
}

type Dialog struct {
	DI IDialog `json:"info"`

	CreatedAt time.Time `json:"created_at" valid:"-" swaggerignore:"true"`
}

func (dialog *Dialog) ToIDialog() *IDialog {
	return &dialog.DI
}

func (iDialog *IDialog) ToDialog(CreatedAt time.Time) *Dialog {
	return &Dialog{
		DI:        *iDialog,
		CreatedAt: CreatedAt,
	}
}
