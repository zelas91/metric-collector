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

func New() *zap.SugaredLogger {

	once.Do(func() {
		file, err := os.ReadFile("config/config.json")
		if err != nil {
			log.Println(err)
			l, err := zap.NewDevelopment()
			if err != nil {
				log.Println(err)
				return
			}
			logger = l.Sugar()
			return
		}

		var cfg zap.Config

		if err := json.Unmarshal(file, &cfg); err != nil {
			log.Fatal(err)
		}
		l, err := cfg.Build()
		logger = l.Sugar()

		if err != nil {
			log.Fatal(err)
		}
	})
	return logger
}
