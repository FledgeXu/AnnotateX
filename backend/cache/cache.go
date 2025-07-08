package cache

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client   *redis.Client
	initOnce sync.Once
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func InitRedis(cfg RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}
