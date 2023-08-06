package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

var addr *string
var pollInterval *int
var reportInterval *int

func init() {
	addr = flag.String("a", "localhost:8080", "endpoint start server")
	pollInterval = flag.Int("p", 2, " poll interval ")
	reportInterval = flag.Int("r", 10, " poll interval ")

}

type Config struct {
	BaseURL        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		logrus.Debugf("read env error=%v", err)
	}
	flag.Parse()
	if cfg.BaseURL == "" {
		cfg.BaseURL = *addr
	}
	if cfg.ReportInterval <= 0 {
		cfg.ReportInterval = *reportInterval
	}

	if cfg.PollInterval <= 0 {
		cfg.PollInterval = *pollInterval
	}
	cfg.BaseURL = fmt.Sprintf("http://%s/update", cfg.BaseURL)
	initLogger()
	return &cfg
}
func initLogger() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}
