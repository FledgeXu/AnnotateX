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
