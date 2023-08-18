package logger

import (
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"os"
)

var Log *zap.SugaredLogger

func init() {
	file, err := os.ReadFile("config/config.json")
	if err != nil {
		log.Fatal("ERROR")
	}
	var cfg zap.Config

	if err := json.Unmarshal(file, &cfg); err != nil {
		log.Fatal(err)
	}
	l, err := cfg.Build()
	Log = l.Sugar()
	if err != nil {
		log.Fatal(err)
	}
}
