package models

import (
	"strconv"
	"time"
	internalError "yula/internal/error"
)

type Advert struct {
	Id          int64     `json:"id" valid:"-" swaggerignore:"true"`
	Name        string    `json:"name" valid:"type(string),stringlength(1|100)" example:"anime's t-shirt"`
	Description string    `json:"description" valid:"optional,stringlength(1|2000)" example:"advert's description"`
	Price       int       `json:"price" valid:"optional,type(int)" example:"100"`
	Location    string    `json:"location" valid:"type(string)" example:"Moscow"`
	Latitude    float64   `json:"latitude" valid:"latitude" example:"55.751244"`
	Longitude   float64   `json:"longitude" valid:"longitude" example:"37.618423"`
	PublishedAt time.Time `json:"published_at" valid:"-" swaggerignore:"true"`
	DateClose   time.Time `json:"date_close" valid:"-" swaggerignore:"true"`
	IsActive    bool      `json:"is_active" valid:"-" swaggerignore:"true"`
	PublisherId int64     `json:"publisher_id" valid:"-" swaggerignore:"true"`
	Category    string    `json:"category" valid:"type(string)" example:"clothes"`
	Images      []string  `json:"images" valid:"-" swaggerignore:"true"`
	Views       int64     `json:"views" valid:"-" swaggerignore:"true"`
}

type AdvertShort struct {
	Id       int64  `json:"id" example:"1"`
	Name     string `json:"name" example:"anime's t-shirt"`
	Price    int    `json:"price" example:"100"`
	Location string `json:"location" example:"Moscow"`
	Image    string `json:"image" example:"/static/advert_images/default_image.png"`
}

func (a *Advert) ToShort() *AdvertShort {
	var imageStr string
	if len(a.Images) == 0 {
		imageStr = ""
	} else {
		imageStr = a.Images[0]
	}
	return &AdvertShort{
		Id: a.Id, Name: a.Name, Price: a.Price, Location: a.Location, Image: imageStr,
	}
}

type Page struct {
	PageNum int64
	Count   int64
}

const (
	DefaultPageNum     int64 = 1
	DefaultCountAdvert int64 = 50
)

func NewPage(pageNumS string, countS string) (*Page, error) {
	pageNum := DefaultPageNum
	count := DefaultCountAdvert
	var err error
	if pageNumS != "" {
		pageNum, err = strconv.ParseInt(pageNumS, 10, 64)
		if err != nil {
			return nil, internalError.BadRequest
		}
	}

	if countS != "" {
		count, err = strconv.ParseInt(countS, 10, 64)
		if err != nil {
			return nil, internalError.BadRequest
		}
	}

	return &Page{PageNum: pageNum - 1, Count: count}, nil
}
