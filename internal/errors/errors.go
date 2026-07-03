package errors

type Code string

const (
	ErrValidation Code = "VALIDATION_ERROR"
	ErrNotFound   Code = "NOT_FOUND"
	ErrConflict   Code = "CONFLICT"
	ErrInternal   Code = "INTERNAL_ERROR"
)

type AppError struct {
	Code    Code
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func Validation(msg string, err error) *AppError {
	return &AppError{
		Code:    ErrValidation,
		Message: msg,
		Err:     err,
	}
}

func Internal(err error) *AppError {
	return &AppError{
		Code:    ErrInternal,
		Message: "internal server error",
		Err:     err,
	}
}

func NotFound(msg string) *AppError {
	return &AppError{
		Code:    ErrNotFound,
		Message: msg,
	}
}

func Conflict(msg string) *AppError {
	return &AppError{
		Code:    ErrConflict,
		Message: msg,
	}
}
