package advicerrors

import (
	"errors"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					http.Error(w, ErrNotFound.message, ErrNotFound.code)
					return
				}
				if errors.Is(err, ErrBadRequest) {
					http.Error(w, ErrBadRequest.message, ErrBadRequest.code)
					return
				}
				if errors.Is(err, ErrMethodNotAllowed) {
					http.Error(w, ErrMethodNotAllowed.message, ErrMethodNotAllowed.code)
					return
				}
			}
		}
	}
}
