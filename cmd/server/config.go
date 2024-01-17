package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/config"
	"os"
)

const (
	defaultAddr          = "localhost:8080"
	defaultCryptoKey     = ""
	defaultStoreInterval = 1
	defaultStoreFile     = "/tmp/metrics-db.json"
	defaultDB            = ""
	defaultRestore       = true
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
	jsonCfg       *string
	trustedSubnet *string
	addrGRPC      *string
	cryptoCert    *string
)

func init() {
	addr = flag.String("a", defaultAddr, "endpoint start server")
	storeInterval = flag.Int("i", defaultStoreInterval, "store interval")
	restore = flag.Bool("r", defaultRestore, "load file metrics")
	filePath = flag.String("f", defaultStoreFile, "file path ")
	database = flag.String("d", defaultDB, "Database URL")
	key = flag.String("k", "", "key hash")
	cryptoKey = flag.String("crypto-key", defaultCryptoKey, "private key")
	jsonCfg = flag.String("c", "", "config json")
	trustedSubnet = flag.String("t", "", "cidr ip")
	addrGRPC = flag.String("grpc-addr", "localhost:3200", "default localhost:3200")
	cryptoCert = flag.String("cert-key", "", "pub key")
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

	if cfg.AddrGRPC == "" {
		cfg.AddrGRPC = *addrGRPC
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
	if cfg.CertPath == "" {
		cfg.CertPath = *cryptoCert
	}
	if cfg.JSONConfig == "" {
		cfg.JSONConfig = *jsonCfg
	}
	if cfg.TrustedSubnet == "" {
		cfg.TrustedSubnet = *trustedSubnet
	}
	if cfg.JSONConfig != "" {
		if data, err := os.ReadFile(cfg.JSONConfig); err == nil {
			configJSON := &config.Config{}
			if err = json.Unmarshal(data, configJSON); err != nil {
				log.Errorf("read json config server err:%v", err)
				return &cfg
			}
			if *cfg.Addr == defaultAddr {
				cfg.Addr = configJSON.Addr
			}
			if cfg.CryptoCertPath == defaultCryptoKey {
				cfg.CryptoCertPath = configJSON.CryptoCertPath
			}
			if *cfg.Restore {
				cfg.Restore = configJSON.Restore
			}
			if *cfg.FilePath == defaultStoreFile {
				cfg.FilePath = configJSON.FilePath
			}

			if *cfg.FilePath == defaultStoreFile {
				cfg.FilePath = configJSON.FilePath
			}
			if *cfg.Database == defaultDB {
				cfg.Database = configJSON.Database
			}
			if *cfg.StoreInterval == defaultStoreInterval {
				cfg.StoreInterval = configJSON.StoreInterval
			}
			if cfg.TrustedSubnet == "" {
				cfg.TrustedSubnet = configJSON.TrustedSubnet
			}
		}
	}
	return &cfg
}

func printVersion() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
