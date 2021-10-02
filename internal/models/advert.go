package models

import (
	"time"
)

type Advert struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	City        string    `json:"location"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	PublishedAt time.Time `json:"published_at"`
	DateClose   time.Time `json:"date_close"`
	IsActive    bool      `json:"is_active"`
	PublisherId int64     `json:"publisher_id"`
	Category    string    `json:"category"`
	Images      []string  `json:"images"`
}
