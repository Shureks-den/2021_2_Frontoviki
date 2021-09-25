package codes

type ServerErrorType uint8

const (
	UserNotExist ServerErrorType = iota + 1
	UserAlreadyExist
	UnexpectedError
	InternalError
	Unauthorized
)

type ServerError struct {
	ErrorCode ServerErrorType
	Message   string
}

var StatusMap = map[ServerErrorType]*ServerError{
	UserNotExist:     {ErrorCode: UserNotExist, Message: "user not exist"},
	UserAlreadyExist: {ErrorCode: UserAlreadyExist, Message: "user already exist"},
	UnexpectedError:  {ErrorCode: UnexpectedError, Message: "unexpected error"},
	InternalError:    {ErrorCode: InternalError, Message: "internal error"},
	Unauthorized:     {ErrorCode: Unauthorized, Message: "no rights to access this resource"},
}

func NewServerError(error ServerErrorType) *ServerError {
	return StatusMap[error]
}
