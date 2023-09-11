package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/config"
)

var log = logger.New()

func NewConfig() *config.Config {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Errorf("read env error=%v", err)
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
	if cfg.Database == nil {
		cfg.Database = new(string)
		flag.StringVar(cfg.Database, "d", "", "Database URL")
	}
	flag.Parse()
	return &cfg
}
