package advicerrors

import "net/http"

type AppError struct {
	message string
	code    int
}

func (e *AppError) Error() string {
	return e.message
}
func NewAppError(message string, code int) *AppError {
	return &AppError{message: message, code: code}
}

func NewErrNotFound(message string) *AppError {
	return NewAppError(message, http.StatusNotFound)
}

func NewErrBadRequest(message string) *AppError {
	return NewAppError(message, http.StatusBadRequest)
}

func NewErrMethodNotAllowed(message string) *AppError {
	return NewAppError(message, http.StatusMethodNotAllowed)
}
