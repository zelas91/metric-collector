package main

import (
	"context"
	"log"
	"net/http"
)

type keyLogger struct{}

func withLogger(ctx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(ctx, keyLogger{}, logger)
}

func getLogger(ctx context.Context) *log.Logger {
	logger, ok := ctx.Value(keyLogger{}).(*log.Logger)
	if !ok {
		// Установка логгера по умолчанию, если он не был передан через контекст
		return log.Default()
	}
	return logger
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	logger := getLogger(req.Context())
	logger.Println("Обработка запроса...")

	// Добавьте здесь свою обработку запроса

	logger.Println("Запрос обработан.")
	handleRequest2(req.Context())
}
func handleRequest2(ctx context.Context) {
	l := getLogger(ctx)
	l.Println("entry  logger 2")
}

func main() {
	//Создание логгера
	logger := log.New(log.Writer(), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Создание HTTP-обработчика с передачей логгера через контекст
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ctx := withLogger(req.Context(), logger)
		handleRequest(w, req.WithContext(ctx))
	})

	// Запуск HTTP-сервера
	http.ListenAndServe(":8000", nil)
}
