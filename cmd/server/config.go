package main

import (
	"flag"
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
)

func init() {
	addr = flag.String("a", "localhost:8080", "endpoint start server")
	storeInterval = flag.Int("i", 300, "store interval")
	restore = flag.Bool("r", true, "load file metrics")
	filePath = flag.String("f", "/tmp/metrics-db.json", "file path ")
	database = flag.String("d", "", "Database URL")
	key = flag.String("k", "", "key hash")
}

func NewConfig() *config.Config {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Errorf("read env error=%v", err)
	}

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
	flag.Parse()
	return &cfg
}
