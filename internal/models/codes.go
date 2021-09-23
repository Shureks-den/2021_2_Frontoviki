package models

import (
	"encoding/json"
	"net/http"
)

type StatusCodeType uint8

const (
	OK StatusCodeType = iota + 1
	Created
	UserNotExist
	InternalError
	BadRequest
)

type Status struct {
	StatusCode StatusCodeType `json:"code"`
	HttpCode   int            `json:"-"`
	Message    string         `json:"message"`
}

var StatusMap = map[StatusCodeType]*Status{
	OK:            {StatusCode: OK, HttpCode: http.StatusOK, Message: "succes"},
	Created:       {StatusCode: Created, HttpCode: http.StatusCreated, Message: "object created"},
	UserNotExist:  {StatusCode: UserNotExist, HttpCode: http.StatusNotFound, Message: "user not exist"},
	InternalError: {StatusCode: InternalError, HttpCode: http.StatusInternalServerError, Message: "unidentified error"},
	BadRequest:    {StatusCode: BadRequest, HttpCode: http.StatusBadRequest, Message: "bad request"},
}

func StatusByCode(code StatusCodeType) *Status {
	status, isExist := StatusMap[code]
	if !isExist {
		return StatusMap[InternalError]
	}
	return status
}

func ToJson(status *Status) []byte {
	jsonStatus, err := json.Marshal(status)
	if err != nil {
		jsonStatus = []byte("")
	}

	return jsonStatus
}
