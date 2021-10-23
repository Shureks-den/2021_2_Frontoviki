package error

import (
	"fmt"
	"net/http"
)

type ServerAnswer struct {
	Code    int
	Message string
}

func (se ServerAnswer) Error() string {
	return fmt.Sprintf("error with code %d happened: %s", se.Code, se.Message)
}

var (
	// определяем ошибки баз данных
	DatabaseError error = ServerAnswer{
		Code:    http.StatusInternalServerError,
		Message: "database error",
	}

	InvalidQuery error = ServerAnswer{
		Code:    http.StatusInternalServerError,
		Message: "invalid query",
	}

	EmptyQuery error = ServerAnswer{
		Code:    http.StatusNotFound,
		Message: "empty rows",
	}

	NotUpdated error = ServerAnswer{
		Code:    http.StatusInternalServerError,
		Message: "not apdated",
	}

	NotCreated error = ServerAnswer{
		Code:    http.StatusConflict,
		Message: "not created",
	}

	RollbackError error = ServerAnswer{
		Code:    http.StatusInternalServerError,
		Message: "rollback error",
	}

	NotCommited error = ServerAnswer{
		Code:    http.StatusInternalServerError,
		Message: "don't commited",
	}

	// определяем ошибки уровня usecase
	NotExist error = ServerAnswer{
		Code:    http.StatusNotFound,
		Message: "not exist",
	}

	AlreadyExist error = ServerAnswer{
		Code:    http.StatusForbidden,
		Message: "already exist",
	}

	InternalError error = ServerAnswer{
		Code:    http.StatusInternalServerError,
		Message: "internal error",
	}

	PasswordMismatch error = ServerAnswer{
		Code:    http.StatusUnauthorized,
		Message: "password mismatch",
	}

	Conflict error = ServerAnswer{
		Code:    http.StatusConflict,
		Message: "unable to access this resource",
	}

	NotEnoughCopies error = ServerAnswer{
		Code:    http.StatusConflict,
		Message: "not enough copies",
	}

	// определяем ошибки уровня http
	BadRequest error = ServerAnswer{
		Code:    http.StatusBadRequest,
		Message: "bad request",
	}

	Unauthorized error = ServerAnswer{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
	}

	// ошибки обработки картинок
	EmptyImageForm error = ServerAnswer{
		Code:    http.StatusBadRequest,
		Message: "require image",
	}

	UnknownExtension error = ServerAnswer{
		Code:    http.StatusBadRequest,
		Message: "file format is not allowed (only PNG, JPEG)",
	}

	UnableToReadFile error = ServerAnswer{
		Code:    http.StatusBadRequest,
		Message: "unable to read file",
	}

	UnableToRemove error = ServerAnswer{
		Code:    http.StatusInternalServerError,
		Message: "unable to remove",
	}
)

func ToMetaStatus(err error) (int, string) {
	answer, ok := err.(ServerAnswer)
	if ok {
		return answer.Code, answer.Message
	}
	return http.StatusInternalServerError, answer.Error()
}

func SetMaxCopies(amount int64) error {
	return ServerAnswer{
		Code:    http.StatusConflict,
		Message: fmt.Sprintf("not enough copies. max amount: %d", amount),
	}
}
