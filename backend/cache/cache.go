package cache

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func InitRedis(cfg RedisConfig) *redis.Client {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		})
	})
	return client
}
