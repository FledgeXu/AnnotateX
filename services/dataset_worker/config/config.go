package config

import (
	"sync"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

var (
	AppConfig *Config
	once      sync.Once
)

type Config struct {
	S3BucketName string `env:"S3_BUCKET_NAME" envDefault:"dev"`
	S3Endpoint   string `env:"S3_ENDPOINT" envDefault:"localhost:9000"`
	S3AccessKey  string `env:"S3_ACCESS_KEY" envDefault:"minioadmin"`
	S3SecretKey  string `env:"S3_SECRET_KEY" envDefault:"minioadmin"`
	S3UseSSL     bool   `env:"S3_USE_SSL" envDefault:"FALSE"`
}

func GetConfig() *Config {
	once.Do(func() {
		AppConfig = load()
	})
	return AppConfig
}

func load() *Config {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		panic(err)
	}
	return &cfg
}
