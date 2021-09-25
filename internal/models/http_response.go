package models

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HttpBodyUser struct {
	User UserWithoutPassword `json:"user"`
}

type HttpUser struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Body    HttpBodyUser `json:"body"`
}
