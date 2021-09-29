package models

import (
	"time"
)

type AdvtData struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Location    string    `json:"location"`
	PublishedAt time.Time `json:"published_at"`
	Image       string    `json:"image"`
	PublisherId int64     `json:"publisher_id"`
	IsActive    bool      `json:"is_active"`
}
