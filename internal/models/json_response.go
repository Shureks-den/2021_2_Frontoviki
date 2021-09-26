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

type HttpBodyProfile struct {
	Profile Profile `json:"profile"`
}

type HttpBodyInterface struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}
