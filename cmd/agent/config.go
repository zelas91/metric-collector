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
	buildVersion   string
	buildDate      string
	buildCommit    string
	cryptoKey      *string
)

func init() {
	addr = flag.String("a", "localhost:8080", "endpoint start server")
	pollInterval = flag.Int("p", 2, " poll interval ")
	reportInterval = flag.Int("r", 10, " poll interval ")
	key = flag.String("k", "", "key hash")
	rateLimit = flag.Int("l", 1, "rate_limit")
	cryptoKey = flag.String("crypto-key", "", "public key")
	printVersion()
}

type Config struct {
	BaseURL        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      *int   `env:"RATE_LIMIT"`
	CryptoKey      string `env:"CRYPTO_KEY"`
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
	if cfg.RateLimit == nil {
		cfg.RateLimit = rateLimit
	}
	if cfg.CryptoKey == "" {
		cfg.CryptoKey = *cryptoKey
	}
	cfg.BaseURL = fmt.Sprintf("http://%s/updates", cfg.BaseURL)
	return &cfg
}
func printVersion() {
	fmt.Printf("Build version: %s\n", getBuildValue(buildVersion))
	fmt.Printf("Build date: %s\n", getBuildValue(buildDate))
	fmt.Printf("Build commit: %s\n", getBuildValue(buildCommit))
}
func getBuildValue(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}
