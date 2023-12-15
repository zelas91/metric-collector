// Package logger to initialize the logger

package logger

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger
var once sync.Once

// Shutdown stop logger.
func Shutdown() {
	if err := logger.Sync(); err != nil {
		log.Printf("logger sync %v", err)
	}
}

// New creates a logger once.
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

		if err = json.Unmarshal(file, &cfg); err != nil {
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
