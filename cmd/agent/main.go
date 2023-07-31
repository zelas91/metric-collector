package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type AppErr struct {
	message string
}
type AppHandler func(w http.ResponseWriter, r *http.Request) *AppErr
type middleware func(h http.Handler) http.Handler

func MiddlewareAdvice(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			http.Error(w, err.message, http.StatusBadRequest)
			return
		}
	}
}

// Middleware function to log incoming requests
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Middleware function to add a timestamp header to the response
func timestampMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("X-Timestamp", time.Now().Format(time.RFC3339))
		next.ServeHTTP(w, r)
	})
}

// Final request handler
func handler(w http.ResponseWriter, r *http.Request) *AppErr {
	fmt.Fprintln(w, "Hello, World!")
	return &AppErr{message: "ASD"}
}

// Chain the middlewares and handler
func applyMiddleware(handler http.Handler, middlewares ...middleware) http.Handler {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}

func main() {
	// Create the final handler by chaining the middlewares and the final handler itself
	finalHandler := applyMiddleware(http.HandlerFunc(MiddlewareAdvice(handler)), logMiddleware, timestampMiddleware)

	// Serve the final handler
	log.Println("Server listen on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", finalHandler))
}
