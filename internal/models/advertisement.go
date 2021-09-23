package models

import (
	"time"

	"github.com/google/uuid"
)

type Advertisement struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Location    string    `json:"location"`
	Published   time.Time `json:"published"`
}
