package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

var addr *string

func init() {
	addr = flag.String("a", "localhost:8080", "endpoint start server")
}

type Config struct {
	Addr string `env:"ADDRESS"`
}

func NewConfig() *Config {
	var cfg Config
	initLogger()
	err := env.Parse(&cfg)
	if err != nil {
		logrus.Debugf("read env error=%v", err)
	}
	if cfg.Addr != "" {
		return &cfg
	}
	flag.Parse()
	return &Config{
		Addr: *addr,
	}
}
func initLogger() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetLevel(logrus.DebugLevel)
	//logrus.SetReportCaller(true)
}
