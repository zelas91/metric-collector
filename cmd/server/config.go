package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
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
	env.Parse(&cfg)
	if cfg.Addr != "" {
		return &cfg
	}
	flag.Parse()
	return &Config{
		Addr: *addr,
	}
}
