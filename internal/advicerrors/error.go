package advicerrors

import "net/http"

var (
	ErrNotFound         = NewAppError(nil, "not found", http.StatusNotFound)
	ErrBadRequest       = NewAppError(nil, "bad request ", http.StatusBadRequest)
	ErrMethodNotAllowed = NewAppError(nil, "status method not allowed", http.StatusMethodNotAllowed)
)

type AppError struct {
	err     error
	message string
	code    int
}

func (e *AppError) Error() string {
	return e.message
}
func NewAppError(err error, message string, code int) *AppError {
	return &AppError{err: err, message: message, code: code}
}
