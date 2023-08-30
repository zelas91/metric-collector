package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/zelas91/metric-collector/internal/logger"
)

var log = logger.New()

type Config struct {
	Addr          *string `env:"ADDRESS"`
	StoreInterval *int    `env:"STORE_INTERVAL"`
	FilePath      *string `env:"FILE_STORAGE_PATH"`
	Restore       *bool   `env:"RESTORE"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Debugf("read env error=%v", err)
	}

	if cfg.Addr == nil {
		cfg.Addr = new(string)
		flag.StringVar(cfg.Addr, "a", "localhost:8080", "endpoint start server")
	}
	if cfg.StoreInterval == nil {
		cfg.StoreInterval = new(int)
		flag.IntVar(cfg.StoreInterval, "i", 300, "store interval")
	}

	if cfg.Restore == nil {
		cfg.Restore = new(bool)
		flag.BoolVar(cfg.Restore, "r", true, "load file metrics")
	}
	if cfg.FilePath == nil {
		cfg.FilePath = new(string)
		flag.StringVar(cfg.FilePath, "f", "/tmp/metrics-db.json", "file path ")

	}
	flag.Parse()
	return &cfg
}
