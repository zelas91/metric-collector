package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/zelas91/metric-collector/internal/logger"
)

var (
	addr           *string
	pollInterval   *int
	reportInterval *int
	key            *string
	rateLimit      *int
	log            = logger.New()
)

func init() {
	addr = flag.String("a", "localhost:8080", "endpoint start server")
	pollInterval = flag.Int("p", 2, " poll interval ")
	reportInterval = flag.Int("r", 10, " poll interval ")
	key = flag.String("k", "", "key hash")
	rateLimit = flag.Int("l", 0, "rate_limit")
}

type Config struct {
	BaseURL        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Errorf("read env error=%v", err)
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
	if cfg.Key == "" {
		cfg.Key = *key
	}
	cfg.BaseURL = fmt.Sprintf("http://%s/updates", cfg.BaseURL)
	return &cfg
}
