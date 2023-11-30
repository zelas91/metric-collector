// Package config - config web service.
package config

type Config struct {
	Addr          *string `env:"ADDRESS"`
	StoreInterval *int    `env:"STORE_INTERVAL"`
	FilePath      *string `env:"FILE_STORAGE_PATH"`
	Restore       *bool   `env:"RESTORE"`
	Database      *string `env:"DATABASE_DSN"`
	Key           *string `env:"KEY"`
}
