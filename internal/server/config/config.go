package config

type Config struct {
	Addr           *string `env:"ADDRESS" json:"address"`
	StoreInterval  *int    `env:"STORE_INTERVAL" json:"store_interval"`
	FilePath       *string `env:"FILE_STORAGE_PATH" json:"store_file"`
	Restore        *bool   `env:"RESTORE" json:"restore"`
	Database       *string `env:"DATABASE_DSN" json:"database_dsn"`
	Key            *string `env:"KEY"`
	CryptoCertPath string  `env:"CRYPTO_KEY" json:"crypto_key"`
	JSONConfig     string  `env:"CONFIG"`
	TrustedSubnet  string  `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}
