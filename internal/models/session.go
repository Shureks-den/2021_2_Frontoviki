package models

import "time"

type Session struct {
	Value     string
	UserId    int64
	ExpiresAt time.Time
}
