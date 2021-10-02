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
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HttpBodyInterface struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}

type HttpBodyProfile struct {
	Profile Profile `json:"profile"`
}

type HttpBodyAdvts struct {
	Advert []*Advert `json:"advert"`
}
