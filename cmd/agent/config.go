package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/zelas91/metric-collector/internal/logger"
)

const (
	defaultAddr           = "localhost:8080"
	defaultCryptoKey      = ""
	defaultReportInterval = 10
	defaultPoolInterval   = 2
)

var (
	addr           *string
	pollInterval   *int
	reportInterval *int
	key            *string
	rateLimit      *int
	log            = logger.New()
	buildVersion   = "N/A"
	buildDate      = "N/A"
	buildCommit    = "N/A"
	cryptoKey      *string
	jsonCfg        *string
)

func init() {
	addr = flag.String("a", defaultAddr, "endpoint start server")
	pollInterval = flag.Int("p", defaultPoolInterval, " poll interval ")
	reportInterval = flag.Int("r", defaultReportInterval, " poll interval ")
	key = flag.String("k", "", "key hash")
	rateLimit = flag.Int("l", 1, "rate_limit")
	cryptoKey = flag.String("crypto-key", defaultCryptoKey, "public key")
	jsonCfg = flag.String("c", "", "config json")
	printVersion()
}

type Config struct {
	BaseURL        string `env:"ADDRESS" json:"address"`
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval"`
	Key            string `env:"KEY"`
	RateLimit      *int   `env:"RATE_LIMIT"`
	CryptoCertPath string `env:"CRYPTO_KEY" json:"crypto_key"`
	JSONConfig     string `env:"CONFIG"`
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
	if cfg.CryptoCertPath == "" {
		cfg.CryptoCertPath = *cryptoKey
	}
	if cfg.JSONConfig != "" {
		cfg.JSONConfig = *jsonCfg
	}
	if cfg.JSONConfig != "" {
		if data, err := os.ReadFile(cfg.JSONConfig); err == nil {
			configJSON := &Config{}
			if err = json.Unmarshal(data, configJSON); err != nil {
				log.Errorf("read json config agent err:%v", err)
				return &cfg
			}
			if cfg.BaseURL == defaultAddr {
				cfg.BaseURL = configJSON.BaseURL
			}
			if cfg.CryptoCertPath == defaultCryptoKey {
				cfg.CryptoCertPath = configJSON.CryptoCertPath
			}
			if cfg.PollInterval == defaultPoolInterval {
				cfg.PollInterval = configJSON.PollInterval
			}
			if cfg.ReportInterval == defaultReportInterval {
				cfg.ReportInterval = configJSON.ReportInterval
			}
		}
	}
	cfg.BaseURL = fmt.Sprintf("http://%s/updates", cfg.BaseURL)
	return &cfg
}
func printVersion() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
