package codes

type DatabaseErrorType uint8

const (
	EmptyRow DatabaseErrorType = iota + 1
	UnexpectedDbError
)

type DatabaseError struct {
	Error DatabaseErrorType
}

func NewDatabaseError(error DatabaseErrorType) *DatabaseError {
	return &DatabaseError{Error: error}
}
