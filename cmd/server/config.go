package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/config"
)

var (
	log           = logger.New()
	addr          *string
	storeInterval *int
	restore       *bool
	filePath      *string
	database      *string
	key           *string
	buildVersion  = "N/A"
	buildDate     = "N/A"
	buildCommit   = "N/A"
	cryptoKey     *string
)

func init() {
	addr = flag.String("a", "localhost:8080", "endpoint start server")
	storeInterval = flag.Int("i", 300, "store interval")
	restore = flag.Bool("r", true, "load file metrics")
	filePath = flag.String("f", "/tmp/metrics-db.json", "file path ")
	database = flag.String("d", "", "Database URL")
	key = flag.String("k", "", "key hash")
	cryptoKey = flag.String("crypto-key", "", "private key")
	printVersion()
}

// NewConfig initialize struct config by environment variables and flags.
func NewConfig() *config.Config {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Errorf("read env error=%v", err)
	}
	flag.Parse()
	if cfg.Addr == nil {
		cfg.Addr = addr
	}
	if cfg.StoreInterval == nil {
		cfg.StoreInterval = storeInterval

	}

	if cfg.Restore == nil {
		cfg.Restore = restore
	}
	if cfg.FilePath == nil {
		cfg.FilePath = filePath

	}
	if cfg.Database == nil {
		cfg.Database = database
	}
	if cfg.Key == nil {
		cfg.Key = key
	}

	if cfg.CryptoCertPath == "" {
		cfg.CryptoCertPath = *cryptoKey
	}

	return &cfg
}

func printVersion() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
