package codes

import "net/http"

type HttpStatus struct {
	Code    int
	Message string
}

var httpStatusMap = map[ServerErrorType]*HttpStatus{
	UserAlreadyExist: {Code: http.StatusConflict, Message: "user with this email already exist"},
	UserNotExist:     {Code: http.StatusNotFound, Message: "user with this email not exist"},
	InternalError:    {Code: http.StatusInternalServerError, Message: "something went wrong"},
	Unauthorized:     {Code: http.StatusUnauthorized, Message: "no rights to access this resource"},
	NotFound:         {Code: http.StatusNotFound, Message: "resource can not be found"},
}

func ServerErrorToHttpStatus(error *ServerError) *HttpStatus {
	if error == nil {
		return &HttpStatus{Code: http.StatusOK}
	}

	return httpStatusMap[error.ErrorCode]
}

func GenCustomStatus(httpCode int, message string) *HttpStatus {
	return &HttpStatus{Code: httpCode, Message: message}
}
