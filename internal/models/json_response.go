package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var defaultJsonAnswer string = fmt.Sprintf(
	`{ "code": %d, "message": "%s" }`,
	http.StatusInternalServerError,
	"internal error",
)

func ToBytes(code int, message string, body interface{}) []byte {
	var response interface{}
	if body == nil {
		response = HttpError{Code: code, Message: message}
	} else {
		response = HttpBodyInterface{Code: code, Message: message, Body: body}
	}

	js := new(bytes.Buffer)
	err := json.NewEncoder(js).Encode(response)
	if err != nil {
		return []byte(defaultJsonAnswer)
	}

	return []byte(js.Bytes())
}

type HttpError struct {
	Code    int    `json:"code" enums:"400,401,403,404,409,500"`
	Message string `json:"message" example:"bad request"`
}

type HttpBodyInterface struct {
	Code    int         `json:"code" enums:"200,201"`
	Message string      `json:"message" enums:"success,created"`
	Body    interface{} `json:"body"`
}

type HttpBodyProfile struct {
	Profile Profile    `json:"profile"`
	Rating  RatingStat `json:"rating"`
}

type HttpBodyAdverts struct {
	Advert []*Advert `json:"adverts"`
}

type HttpBodyAdvertShort struct {
	AdvertShort AdvertShort `json:"advert"`
}

type HttpBodyAdvert struct {
	Advert Advert `json:"advert"`
}

type HttpBodyAdvertDetail struct {
	Advert        Advert         `json:"advert"`
	Salesman      Profile        `json:"salesman"`
	Rating        RatingStat     `json:"rating"`
	PriceHistory  []*AdvertPrice `json:"price_history"`
	FavoriteCount int64          `json:"favorite_count"`
}

type HttpBodySalesmanPage struct {
	Salesman Profile        `json:"salesman"`
	Adverts  []*AdvertShort `json:"adverts"`
	Rating   RatingStat     `json:"rating"`
}

type HttpBodyCartAll struct {
	Cart    []*Cart   `json:"cart"`
	Adverts []*Advert `json:"adverts"`
	Hints   []string  `json:"hints"`
}

type HttpBodyCart struct {
	Cart    []*Cart   `json:"cart"`
	Adverts []*Advert `json:"adverts"`
}

type HttpBodyOrder struct {
	Salesman Profile `json:"salesman"`
	Order    Cart    `json:"order"`
}

type HttpBodyCategories struct {
	Categories []*Category `json:"categories"`
}

type HttpBodyChatHistory struct {
	Messages []*Message `json:"messages"`
}

type HttpDialog struct {
	Id      int64  `json:"id" valid:"int"`
	Name    string `json:"name" valid:"type(string)"`
	Surname string `json:"surname" valid:"type(string)"`

	Adv AdvertShort `json:"adv_info"`

	CreatedAt time.Time `json:"created_at" valid:"-" swaggerignore:"true"`
}

type HttpBodyDialogs struct {
	Dialogs []*HttpDialog `json:"dialogs"`
}

type HttpBodyPriceHistory struct {
	History []*AdvertPrice `json:"history"`
}
