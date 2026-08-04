package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config is the validated application config - no silent defaults for secrets (SAD).
type Config struct {
	Env              string `mapstructure:"env"` // local, staging, production, pilot
	Port             int    `mapstructure:"port"`
	DatabaseURL      string `mapstructure:"database_url"`
	RedisURL         string `mapstructure:"redis_url"`
	MinIOEndpoint    string `mapstructure:"minio_endpoint"`
	MinIOAccessKey   string `mapstructure:"minio_access_key"`
	MinIOSecretKey   string `mapstructure:"minio_secret_key"`
	MinIOBucket      string `mapstructure:"minio_bucket"`
	MinIOUseSSL      bool   `mapstructure:"minio_use_ssl"`
	JWTSecret        string `mapstructure:"jwt_secret"`
	ConnectorEncKey  string `mapstructure:"connector_encryption_key"`
	FaydaMode        string `mapstructure:"fayda_mode"` // mock | live
	FaydaPartnerCode string `mapstructure:"fayda_partner_code"`
	FaydaPartnerKey  string `mapstructure:"fayda_partner_key"`
	FaydaBaseURL     string `mapstructure:"fayda_base_url"`
	RedisCacheTTL    int    `mapstructure:"redis_cache_ttl_seconds"`
}

func Load() (*Config, error) {
	_ = godotenv.Load() // optional for local

	viper.AutomaticEnv()
	viper.SetDefault("env", "local")
	viper.SetDefault("port", 8080)
	viper.SetDefault("minio_bucket", "apexpay-vault")
	viper.SetDefault("minio_use_ssl", false)
	viper.SetDefault("fayda_mode", "mock")
	viper.SetDefault("fayda_base_url", "https://id.gov.et/api")
	viper.SetDefault("redis_cache_ttl_seconds", 300)

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Required secrets - fail fast, no silent defaults
	if cfg.DatabaseURL == "" {
		if env := os.Getenv("DATABASE_URL"); env != "" {
			cfg.DatabaseURL = env
		} else {
			return nil, fmt.Errorf("DATABASE_URL required")
		}
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET required")
	}
	if cfg.ConnectorEncKey == "" && cfg.Env != "local" {
		return nil, fmt.Errorf("CONNECTOR_ENCRYPTION_KEY required in non-local env")
	}
	if cfg.FaydaMode == "live" && cfg.FaydaPartnerKey == "" {
		return nil, fmt.Errorf("FAYDA_PARTNER_KEY required in live mode")
	}

	return &cfg, nil
}
