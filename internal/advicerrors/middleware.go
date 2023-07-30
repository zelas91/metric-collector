package advicerrors

import (
	"net/http"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) *AppError

func Middleware(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			http.Error(w, err.message, err.code)
			return
		}
	}
}
