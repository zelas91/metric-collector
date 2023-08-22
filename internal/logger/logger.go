package logger

import (
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"os"
	"sync"
)

var logger *zap.SugaredLogger
var once sync.Once

func Shutdown() {
	if err := logger.Sync(); err != nil {
		log.Fatal(err)
	}
}
func New() *zap.SugaredLogger {

	once.Do(func() {
		file, err := os.ReadFile("config/config.json")
		if err != nil {
			log.Println(err)
			l, err := zap.NewDevelopment()
			if err != nil {
				log.Fatal(err)
			}
			logger = l.Sugar()
			return
		}

		var cfg zap.Config

		if err := json.Unmarshal(file, &cfg); err != nil {
			log.Fatal(err)
		}
		l, err := cfg.Build()
		if err != nil {
			log.Fatal(err)
		}
		logger = l.Sugar()
	})
	return logger
}
