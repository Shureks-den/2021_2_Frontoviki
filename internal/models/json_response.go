package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	Profile Profile `json:"profile"`
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
	Advert   Advert  `json:"advert"`
	Salesman Profile `json:"salesman"`
}

type HttpBodySalesmanPage struct {
	Salesman Profile        `json:"salesman"`
	Adverts  []*AdvertShort `json:"adverts"`
}
